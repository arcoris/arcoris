/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package retry

import (
	"context"
	"time"
)

// run executes op with retry orchestration.
//
// run is the private retry engine behind Do and DoValue. It owns the runtime
// control flow that combines operation execution, retryability classification,
// retry-owned limits, backoff sequence consumption, clock-backed delays, context
// interruption, and observer events.
//
// run does not define operation semantics. The caller remains responsible for
// idempotency, replay safety, transactional safety, and external side effects.
// run also does not define protocol-specific retry rules, storage retry rules,
// retry budgets, circuit breakers, hedging, metrics exporters, tracing exporters,
// or logging backends.
//
// The returned value is meaningful only when err is nil. On failure, run returns
// the zero value of T and an error describing the terminal retry decision.
func run[T any](
	ctx context.Context,
	op ValueOperation[T],
	config config,
) (T, error) {
	requireContext(ctx)
	requireValueOperation(op)
	requireClock(config.clock)
	requireBackoff(config.backoff)
	requireClassifier(config.classifier)
	requireMaxAttempts(config.maxAttempts)
	requireMaxElapsed(config.maxElapsed)

	var zero T

	sequence := config.backoff.NewSequence()
	requireBackoffSequence(sequence)

	startedAt := config.clock.Now()

	var attempts uint
	var lastAttempt Attempt
	var lastErr error

	for {
		if err := contextStopError(ctx); err != nil {
			outcome := newOutcome(
				config,
				startedAt,
				attempts,
				lastErr,
				StopReasonInterrupted,
			)
			emitRetryEvent(ctx, config, stopEvent(outcome, lastAttempt))
			return zero, err
		}

		attempts++
		attempt := Attempt{
			Number:    attempts,
			StartedAt: config.clock.Now(),
		}
		lastAttempt = attempt

		emitRetryEvent(ctx, config, Event{
			Kind:    EventAttemptStart,
			Attempt: attempt,
		})

		value, err := op(ctx)
		if err == nil {
			outcome := newOutcome(
				config,
				startedAt,
				attempts,
				nil,
				StopReasonSucceeded,
			)
			emitRetryEvent(ctx, config, stopEvent(outcome, attempt))
			return value, nil
		}

		lastErr = err

		emitRetryEvent(ctx, config, Event{
			Kind:    EventAttemptFailure,
			Attempt: attempt,
			Err:     err,
		})

		if !config.classifier.Retryable(err) {
			outcome := newOutcome(
				config,
				startedAt,
				attempts,
				err,
				StopReasonNonRetryable,
			)
			emitRetryEvent(ctx, config, stopEvent(outcome, attempt))
			return zero, err
		}

		if attempts >= config.maxAttempts {
			outcome := newOutcome(
				config,
				startedAt,
				attempts,
				err,
				StopReasonMaxAttempts,
			)
			emitRetryEvent(ctx, config, stopEvent(outcome, attempt))
			return zero, NewExhaustedError(outcome)
		}

		if interruptErr := contextStopError(ctx); interruptErr != nil {
			outcome := newOutcome(
				config,
				startedAt,
				attempts,
				err,
				StopReasonInterrupted,
			)
			emitRetryEvent(ctx, config, stopEvent(outcome, attempt))
			return zero, interruptErr
		}

		delay, ok := sequence.Next()
		requireBackoffDelay(delay, ok)

		if !ok {
			outcome := newOutcome(
				config,
				startedAt,
				attempts,
				err,
				StopReasonBackoffExhausted,
			)
			emitRetryEvent(ctx, config, stopEvent(outcome, attempt))
			return zero, NewExhaustedError(outcome)
		}

		if maxElapsedWouldBeExceeded(config, startedAt, delay) {
			outcome := newOutcome(
				config,
				startedAt,
				attempts,
				err,
				StopReasonMaxElapsed,
			)
			emitRetryEvent(ctx, config, stopEvent(outcome, attempt))
			return zero, NewExhaustedError(outcome)
		}

		emitRetryEvent(ctx, config, Event{
			Kind:    EventRetryDelay,
			Attempt: attempt,
			Delay:   delay,
			Err:     err,
		})

		if waitErr := waitDelay(ctx, config.clock, delay); waitErr != nil {
			outcome := newOutcome(
				config,
				startedAt,
				attempts,
				err,
				StopReasonInterrupted,
			)
			emitRetryEvent(ctx, config, stopEvent(outcome, attempt))
			return zero, waitErr
		}
	}
}

// newOutcome constructs retry completion metadata using the configured clock.
//
// The helper centralizes FinishedAt assignment so every terminal path records
// completion time consistently.
func newOutcome(
	config config,
	startedAt time.Time,
	attempts uint,
	lastErr error,
	reason StopReason,
) Outcome {
	return Outcome{
		Attempts:   attempts,
		StartedAt:  startedAt,
		FinishedAt: config.clock.Now(),
		LastErr:    lastErr,
		Reason:     reason,
	}
}

// stopEvent constructs the terminal observer event for outcome.
//
// When retry stops before any operation attempt, the event carries only Outcome.
// When retry stops after one or more operation calls, the event carries the last
// Attempt and mirrors Outcome.LastErr through Event.Err.
func stopEvent(outcome Outcome, attempt Attempt) Event {
	if outcome.Attempts == 0 {
		return Event{
			Kind:    EventRetryStop,
			Outcome: outcome,
		}
	}

	return Event{
		Kind:    EventRetryStop,
		Attempt: attempt,
		Err:     outcome.LastErr,
		Outcome: outcome,
	}
}

// emitRetryEvent notifies configured observers about event.
//
// Observers are called synchronously in registration order. Observer failures are
// not represented in retry's error model because Observer does not return error.
// The retry package does not recover observer panics.
func emitRetryEvent(ctx context.Context, config config, event Event) {
	for _, observer := range config.observers {
		observer.ObserveRetry(ctx, event)
	}
}

// maxElapsedWouldBeExceeded reports whether waiting for delay would violate the
// configured max-elapsed boundary.
//
// A zero maxElapsed disables elapsed-time limiting. The check is deliberately
// conservative: if the selected delay would consume all remaining time, retry
// stops with StopReasonMaxElapsed instead of sleeping until no execution budget
// remains for the next attempt.
func maxElapsedWouldBeExceeded(
	config config,
	startedAt time.Time,
	delay time.Duration,
) bool {
	if config.maxElapsed == 0 {
		return false
	}

	elapsed := config.clock.Since(startedAt)
	if elapsed >= config.maxElapsed {
		return true
	}

	remaining := config.maxElapsed - elapsed
	return delay >= remaining
}

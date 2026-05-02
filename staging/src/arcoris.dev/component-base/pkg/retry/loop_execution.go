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

	"arcoris.dev/component-base/pkg/backoff"
)

// retryExecution owns the mutable state for one Do or DoValue run.
//
// The state is intentionally single-owner, not concurrency-safe, and never
// shared across retry executions. It does not store context.Context; context is
// passed explicitly to every helper that observes cancellation or emits observer
// events, following Go's context ownership rule.
type retryExecution struct {
	config   config
	sequence backoff.Sequence

	startedAt time.Time
	attempts  uint

	lastAttempt Attempt
	lastErr     error
}

// newRetryExecution validates retry runtime dependencies and starts one run.
//
// Public configuration stores a reusable backoff.Schedule. Runtime execution
// immediately creates a fresh Sequence so delay state is owned by this run only.
func newRetryExecution(config config) *retryExecution {
	requireClock(config.clock)
	requireBackoff(config.backoff)
	requireClassifier(config.classifier)
	requireMaxAttempts(config.maxAttempts)
	requireMaxElapsed(config.maxElapsed)

	sequence := config.backoff.NewSequence()
	requireBackoffSequence(sequence)

	return &retryExecution{
		config:    config,
		sequence:  sequence,
		startedAt: config.clock.Now(),
	}
}

func (e *retryExecution) contextStop(ctx context.Context) error {
	return contextStopError(ctx)
}

// nextAttempt records and emits the next operation attempt.
//
// Attempt numbers are one-based and include the initial call. The emitted start
// event is observer-visible before the operation is invoked.
func (e *retryExecution) nextAttempt(ctx context.Context) Attempt {
	e.attempts++
	e.lastAttempt = Attempt{
		Number:    e.attempts,
		StartedAt: e.config.clock.Now(),
	}
	e.emit(ctx, Event{
		Kind:    EventAttemptStart,
		Attempt: e.lastAttempt,
	})

	return e.lastAttempt
}

// recordFailure stores the operation-owned failure for later retry decisions.
//
// lastErr is reserved for operation errors. Retry-owned context interruption is
// returned separately and must not replace the last operation error in Outcome.
func (e *retryExecution) recordFailure(ctx context.Context, attempt Attempt, err error) {
	e.lastAttempt = attempt
	e.lastErr = err
	e.emit(ctx, Event{
		Kind:    EventAttemptFailure,
		Attempt: attempt,
		Err:     err,
	})
}

func (e *retryExecution) retryable(err error) bool {
	return e.config.classifier.Retryable(err)
}

func (e *retryExecution) maxAttemptsReached() bool {
	return e.attempts >= e.config.maxAttempts
}

// nextDelay consumes the execution-owned backoff sequence.
//
// A negative delay with ok=true violates the Sequence contract and panics via
// the stable retry diagnostic. ok=false is not a backoff error; run maps it to
// retry-owned backoff exhaustion.
func (e *retryExecution) nextDelay() (time.Duration, bool) {
	delay, ok := e.sequence.Next()
	requireBackoffDelay(delay, ok)
	return delay, ok
}

// retryDelay emits the delay selected after the current failed retryable attempt.
//
// It assumes recordFailure has already stored lastAttempt and lastErr for the
// retry boundary that selected this delay.
func (e *retryExecution) retryDelay(ctx context.Context, delay time.Duration) {
	e.emit(ctx, Event{
		Kind:    EventRetryDelay,
		Attempt: e.lastAttempt,
		Delay:   delay,
		Err:     e.lastErr,
	})
}

// waitDelay waits at the retry-owned boundary between two attempts.
//
// The method exists only to keep run focused on control flow while preserving
// the package-level waitDelay semantics: context interruption observed here is
// retry-owned and must be returned as ErrInterrupted.
func (e *retryExecution) waitDelay(ctx context.Context, delay time.Duration) error {
	return waitDelay(ctx, e.config.clock, delay)
}

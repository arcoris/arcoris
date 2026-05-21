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

	"arcoris.dev/chrono/delay"
)

// retryExecution owns the mutable state for one Do or DoValue run.
//
// The state is intentionally single-owner, not concurrency-safe, and never
// shared across retry executions. It does not store context.Context; context is
// passed explicitly to every helper that observes cancellation or emits observer
// events, following Go's context ownership rule.
type retryExecution struct {
	// config is the normalized immutable configuration used by this execution.
	//
	// The value is copied from configOf before the loop starts. Options are not
	// retained and cannot mutate it after execution begins.
	config config

	// sequence is the delay stream owned by this execution.
	//
	// It is created from config.delay exactly once. The sequence is mutable,
	// single-owner state and must not be shared between retry executions.
	sequence delay.Sequence

	// startedAt is the clock timestamp captured before the first attempt.
	//
	// Elapsed-time limits and terminal outcomes use this same timestamp so event
	// metadata and limit checks agree on the execution boundary.
	startedAt time.Time

	// attempts is the number of operation attempts that have started.
	//
	// It includes the initial operation call and is incremented before each
	// operation invocation.
	attempts uint

	// lastAttempt is the most recent operation attempt metadata.
	//
	// Terminal stop events reuse it when at least one attempt was started.
	lastAttempt Attempt

	// lastErr is the most recent operation-owned error.
	//
	// Retry-owned context interruption is returned separately and must not
	// replace this field.
	lastErr error
}

// newRetryExecution validates retry runtime dependencies and starts one run.
//
// Public configuration stores a reusable delay.Schedule. Runtime execution
// immediately creates a fresh Sequence so delay state is owned by this run only.
func newRetryExecution(cfg config) *retryExecution {
	requireClock(cfg.clock)
	requireDelaySchedule(cfg.delay)
	requireClassifier(cfg.classifier)
	requireMaxAttempts(cfg.maxAttempts)
	requireMaxElapsed(cfg.maxElapsed)

	seq := cfg.delay.NewSequence()
	requireDelaySequence(seq)

	return &retryExecution{
		config:    cfg,
		sequence:  seq,
		startedAt: cfg.clock.Now(),
	}
}

// contextStop returns retry-owned context interruption if ctx has stopped.
//
// The wrapper keeps context observation readable at call sites and centralizes
// the rule that retry-owned context errors are returned as ErrInterrupted.
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
func (e *retryExecution) recordFailure(ctx context.Context, a Attempt, err error) {
	e.lastAttempt = a
	e.lastErr = err
	e.emit(ctx, Event{
		Kind:    EventAttemptFailure,
		Attempt: a,
		Err:     err,
	})
}

// retryable reports whether err is retryable according to the configured
// classifier.
//
// The error remains operation-owned. This method only delegates the policy
// decision and does not wrap, classify, or store the error itself.
func (e *retryExecution) retryable(err error) bool {
	return e.config.classifier.Retryable(err)
}

// maxAttemptsReached reports whether the attempt budget has been consumed.
//
// The comparison uses attempts after the current attempt has already started, so
// reaching the limit means retry must stop before scheduling another attempt.
func (e *retryExecution) maxAttemptsReached() bool {
	return e.attempts >= e.config.maxAttempts
}

// nextDelay consumes the execution-owned delay sequence.
//
// A negative delay with ok=true violates the Sequence contract and panics via
// the stable retry diagnostic. ok=false is not a delay error; run maps it to
// retry-owned delay exhaustion.
func (e *retryExecution) nextDelay() (time.Duration, bool) {
	d, ok := e.sequence.Next()
	requireDelay(d, ok)
	return d, ok
}

// retryDelay emits the delay selected after the current failed retryable attempt.
//
// It assumes recordFailure has already stored lastAttempt and lastErr for the
// retry boundary that selected this delay.
func (e *retryExecution) retryDelay(ctx context.Context, d time.Duration) {
	e.emit(ctx, Event{
		Kind:    EventRetryDelay,
		Attempt: e.lastAttempt,
		Delay:   d,
		Err:     e.lastErr,
	})
}

// waitDelay waits at the retry-owned boundary between two attempts.
//
// The method exists only to keep run focused on control flow while preserving
// the package-level waitDelay semantics: context interruption observed here is
// retry-owned and must be returned as ErrInterrupted.
func (e *retryExecution) waitDelay(ctx context.Context, d time.Duration) error {
	return waitDelay(ctx, e.config.clock, d)
}

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
// The state is intentionally package-private and single-owner. Public options
// still store a reusable backoff.Schedule; each execution creates and consumes
// its own backoff.Sequence.
type retryExecution struct {
	ctx      context.Context
	config   config
	sequence backoff.Sequence

	startedAt time.Time
	attempts  uint

	lastAttempt Attempt
	lastErr     error
}

func newRetryExecution(ctx context.Context, config config) *retryExecution {
	requireClock(config.clock)
	requireBackoff(config.backoff)
	requireClassifier(config.classifier)
	requireMaxAttempts(config.maxAttempts)
	requireMaxElapsed(config.maxElapsed)

	sequence := config.backoff.NewSequence()
	requireBackoffSequence(sequence)

	return &retryExecution{
		ctx:       ctx,
		config:    config,
		sequence:  sequence,
		startedAt: config.clock.Now(),
	}
}

func (e *retryExecution) contextStop() error {
	return contextStopError(e.ctx)
}

func (e *retryExecution) nextAttempt() Attempt {
	e.attempts++
	e.lastAttempt = Attempt{
		Number:    e.attempts,
		StartedAt: e.config.clock.Now(),
	}
	e.emit(Event{
		Kind:    EventAttemptStart,
		Attempt: e.lastAttempt,
	})

	return e.lastAttempt
}

func (e *retryExecution) recordFailure(attempt Attempt, err error) {
	e.lastAttempt = attempt
	e.lastErr = err
	e.emit(Event{
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

func (e *retryExecution) nextDelay() (time.Duration, bool) {
	delay, ok := e.sequence.Next()
	requireBackoffDelay(delay, ok)
	return delay, ok
}

func (e *retryExecution) retryDelay(delay time.Duration) {
	e.emit(Event{
		Kind:    EventRetryDelay,
		Attempt: e.lastAttempt,
		Delay:   delay,
		Err:     e.lastErr,
	})
}

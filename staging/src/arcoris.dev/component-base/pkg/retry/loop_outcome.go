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

import "context"

// succeeded emits the terminal success event.
//
// Success must not preserve earlier operation errors after a later retry attempt
// succeeds, so LastErr is always nil for StopReasonSucceeded.
func (e *retryExecution) succeeded(ctx context.Context) {
	e.stop(ctx, StopReasonSucceeded, nil)
}

// nonRetryable emits the terminal event and returns the original operation error.
//
// Non-retryable errors are operation-owned results, so retry must not wrap them
// as ErrExhausted or ErrInterrupted.
func (e *retryExecution) nonRetryable(ctx context.Context, err error) error {
	e.stop(ctx, StopReasonNonRetryable, err)
	return err
}

// exhausted emits a retry-owned exhaustion outcome.
//
// Callers must pass only exhausted StopReason values so NewExhaustedError
// preserves the package invariant that ErrExhausted belongs to retry-owned
// attempt, elapsed, or backoff exhaustion.
func (e *retryExecution) exhausted(ctx context.Context, reason StopReason) error {
	outcome := e.stop(ctx, reason, e.lastErr)
	return NewExhaustedError(outcome)
}

// interrupted emits a retry-owned interruption outcome and returns err unchanged.
//
// The interruption cause is not stored in Outcome.LastErr. LastErr is reserved
// for the last operation-owned error, if an operation attempt happened before
// retry observed context cancellation.
func (e *retryExecution) interrupted(ctx context.Context, err error) error {
	e.stop(ctx, StopReasonInterrupted, e.lastErr)
	return err
}

// stop centralizes terminal Outcome construction and EventRetryStop emission.
func (e *retryExecution) stop(ctx context.Context, reason StopReason, lastErr error) Outcome {
	outcome := e.newOutcome(reason, lastErr)
	e.emit(ctx, stopEvent(outcome, e.lastAttempt))
	return outcome
}

// newOutcome assigns terminal metadata consistently for every stop path.
func (e *retryExecution) newOutcome(reason StopReason, lastErr error) Outcome {
	return Outcome{
		Attempts:   e.attempts,
		StartedAt:  e.startedAt,
		FinishedAt: e.config.clock.Now(),
		LastErr:    lastErr,
		Reason:     reason,
	}
}

// stopEvent constructs the observer-visible terminal event for outcome.
//
// Stops before the first operation omit Attempt and Err. Stops after an
// operation mirror Outcome.LastErr through Event.Err so observer metadata stays
// structurally valid for every terminal reason.
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

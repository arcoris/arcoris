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

import "time"

// Event describes one observer-visible retry event.
//
// Event is immutable retry metadata produced by retry execution and delivered to
// configured observers. It does not execute operations, decide retryability,
// compute retry delays, wait on timers, wrap errors, mutate retry state, or
// change the retry decision.
//
// Event values are copyable. They must not contain locks, channels, mutable
// delay sequences, operations, timers, or references back to the retry
// execution that produced them.
//
// Field interpretation depends on Kind:
//
//   - EventAttemptStart carries Attempt;
//   - EventAttemptFailure carries Attempt and Err;
//   - EventRetryDelay carries Attempt, Err, and Delay;
//   - EventRetryStop carries Outcome and may carry the last Attempt and Err.
//
// The zero Event is invalid and represents an omitted event value.
type Event struct {
	// Kind identifies the event shape and determines which payload fields are
	// meaningful.
	Kind EventKind

	// Attempt describes the operation attempt associated with this event.
	//
	// Attempt is required for attempt-scoped events. For EventRetryStop, Attempt
	// is required only when Outcome.Attempts is non-zero; in that case it should
	// describe the last operation attempt observed by the retry execution.
	Attempt Attempt

	// Delay is the retry delay selected after a failed retryable attempt.
	//
	// Delay is meaningful only for EventRetryDelay. A zero delay is valid and
	// represents an immediate retry boundary. Negative delays are invalid because
	// delay sequences must not produce negative wait durations for retry.
	Delay time.Duration

	// Err is the operation-owned error associated with this event.
	//
	// Err is required for EventAttemptFailure and EventRetryDelay. For
	// EventRetryStop, Err should mirror the nilness of Outcome.LastErr. Successful
	// outcomes and interruptions before the first attempt do not carry Err.
	Err error

	// Outcome is the terminal retry execution metadata.
	//
	// Outcome is meaningful only for EventRetryStop. Non-terminal events must keep
	// Outcome as the zero value.
	Outcome Outcome
}

// IsZero reports whether e is the zero Event value.
//
// The zero Event is not valid observer metadata. It is useful as an omitted
// value in tests or future structs that may optionally include an event.
func (e Event) IsZero() bool {
	return e.Kind == 0 &&
		e.Attempt.IsZero() &&
		e.Delay == 0 &&
		e.Err == nil &&
		e.Outcome.IsZero()
}

// IsValid reports whether e is structurally consistent for its event kind.
//
// IsValid checks event-shape invariants only. It does not verify that the event
// was emitted by a real retry loop, that Err is retryable, that Delay came from a
// specific delay sequence, or that timestamps came from a particular clock.
func (e Event) IsValid() bool {
	switch e.Kind {
	case EventAttemptStart:
		return e.isValidAttemptStartEvent()
	case EventAttemptFailure:
		return e.isValidAttemptFailureEvent()
	case EventRetryDelay:
		return e.isValidRetryDelayEvent()
	case EventRetryStop:
		return e.isValidRetryStopEvent()
	default:
		return false
	}
}

// isValidAttemptStartEvent reports whether e is a valid attempt-start event.
//
// Attempt-start events are emitted before the operation is called. They carry
// only the attempt metadata because no error, delay, or terminal outcome exists
// yet.
func (e Event) isValidAttemptStartEvent() bool {
	return e.Attempt.IsValid() &&
		e.Delay == 0 &&
		e.Err == nil &&
		e.Outcome.IsZero()
}

// isValidAttemptFailureEvent reports whether e is a valid attempt-failure event.
//
// Attempt-failure events are emitted after an operation attempt returns a non-nil
// operation-owned error. They do not carry retry delay or terminal outcome
// metadata because retry has not necessarily decided what happens next.
func (e Event) isValidAttemptFailureEvent() bool {
	return e.Attempt.IsValid() &&
		e.Delay == 0 &&
		e.Err != nil &&
		e.Outcome.IsZero()
}

// isValidRetryDelayEvent reports whether e is a valid retry-delay event.
//
// Retry-delay events are emitted after retry has classified a failed attempt as
// retryable, selected a delay from the configured delay sequence, and decided
// to wait before another attempt. A zero delay is valid; a negative delay is not.
func (e Event) isValidRetryDelayEvent() bool {
	return e.Attempt.IsValid() &&
		e.Delay >= 0 &&
		e.Err != nil &&
		e.Outcome.IsZero()
}

// isValidRetryStopEvent reports whether e is a valid retry-stop event.
//
// Stop events carry a valid Outcome. When retry stops before the first operation
// call, Attempt and Err must be omitted. When retry stops after one or more
// operation calls, Attempt must describe the last attempt and Err must match the
// nilness of Outcome.LastErr.
//
// The method intentionally compares only error nilness and not concrete error
// identity: direct comparison between two non-nil error interface values can
// panic when the dynamic error value is not comparable.
func (e Event) isValidRetryStopEvent() bool {
	if !e.Outcome.IsValid() || e.Delay != 0 {
		return false
	}

	if e.Outcome.Attempts == 0 {
		return e.Attempt.IsZero() && e.Err == nil
	}

	if !e.Attempt.IsValid() || e.Attempt.Number != e.Outcome.Attempts {
		return false
	}

	if e.Outcome.LastErr == nil {
		return e.Err == nil
	}

	return e.Err != nil
}

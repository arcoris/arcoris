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

// EventKind identifies the shape of a retry observer event.
//
// EventKind is part of the observer-facing retry metadata model. It tells
// observers which payload fields in Event are meaningful. It does not represent
// retryability, delay policy, stop reason, operation status, context state, or
// protocol-specific behavior.
//
// The zero value is invalid. This prevents an empty Event from accidentally
// looking like a valid observer notification.
//
// EventKind values are intended for diagnostics, tests, observer adapters,
// metrics adapters, tracing adapters, and logs. They are not stable external
// wire-format values. If a future public API needs serialized retry events, that
// API should define its own compatibility contract instead of serializing this
// runtime enum directly.
type EventKind uint8

const (
	// EventAttemptStart is emitted immediately before an operation attempt is
	// called.
	//
	// The corresponding Event carries valid Attempt metadata and no error,
	// delay, or outcome payload. At this point the operation has not completed
	// and retry does not yet know whether the attempt will succeed, fail, be
	// retried, or stop execution.
	EventAttemptStart EventKind = iota + 1

	// EventAttemptFailure is emitted after an operation attempt returns a non-nil
	// operation-owned error.
	//
	// The corresponding Event carries the failed Attempt and the operation-owned
	// error. It does not decide whether another attempt will be scheduled. Retry
	// still has to apply the classifier, limits, context state, and delay
	// sequence after the failed attempt is observed.
	EventAttemptFailure

	// EventRetryDelay is emitted after a failed retryable attempt when retry has
	// selected a delay before the next attempt.
	//
	// The corresponding Event carries the failed Attempt, the operation-owned
	// error that led to the delay, and the selected delay. A zero delay is valid
	// and means the next attempt may be scheduled immediately after retry-owned
	// context checks.
	EventRetryDelay

	// EventRetryStop is emitted when retry execution reaches a terminal decision.
	//
	// The corresponding Event carries a valid Outcome. It may also carry the last
	// Attempt and last operation-owned error when retry stopped after at least one
	// operation call.
	EventRetryStop
)

// String returns the canonical lower-case diagnostic name of k.
//
// The returned value is intended for diagnostics, tests, logs, observer events,
// and error messages. It is not a versioned serialization format. Unknown values
// return "invalid" so callers never accidentally render an unknown numeric value
// as a valid event kind.
func (k EventKind) String() string {
	switch k {
	case EventAttemptStart:
		return "attempt_start"
	case EventAttemptFailure:
		return "attempt_failure"
	case EventRetryDelay:
		return "retry_delay"
	case EventRetryStop:
		return "retry_stop"
	default:
		return "invalid"
	}
}

// IsValid reports whether k is one of the event kinds defined by this package.
//
// IsValid is useful at package boundaries, in observer tests, and in defensive
// code that receives an EventKind from caller-controlled input. The zero value
// is invalid and must not be treated as a no-op or as EventRetryStop.
func (k EventKind) IsValid() bool {
	switch k {
	case EventAttemptStart,
		EventAttemptFailure,
		EventRetryDelay,
		EventRetryStop:
		return true
	default:
		return false
	}
}

// IsAttemptScoped reports whether k describes a specific operation attempt.
//
// Attempt-scoped event kinds require valid Attempt metadata and do not carry a
// terminal Outcome. EventRetryStop is not attempt-scoped because retry execution
// can stop before the first operation call.
func (k EventKind) IsAttemptScoped() bool {
	switch k {
	case EventAttemptStart,
		EventAttemptFailure,
		EventRetryDelay:
		return true
	default:
		return false
	}
}

// IsTerminal reports whether k records the terminal retry observer event.
//
// Only EventRetryStop is terminal. Terminal here refers to retry observer event
// flow, not to lifecycle state, operation state, or external protocol state.
func (k EventKind) IsTerminal() bool {
	return k == EventRetryStop
}

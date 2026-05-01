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

// StopReason identifies why one retry execution stopped.
//
// StopReason is completion metadata. It explains the terminal reason recorded in
// an Outcome or stop event, but it does not carry the operation error, context
// cause, backoff sequence, attempt metadata, observer state, or retry
// configuration that produced that reason.
//
// The zero value is invalid. This keeps an empty Outcome from accidentally
// looking like a successful retry execution. Valid reasons are assigned
// explicitly by the retry loop when execution reaches a terminal decision.
//
// StopReason values are intended for diagnostics, tests, observer events,
// structured retry errors, and future metrics or tracing adapters. They are not
// stable external wire-format values. If a future API needs serialized retry
// outcomes, that API should define its own compatibility contract instead of
// serializing this runtime enum directly.
type StopReason uint8

const (
	// StopReasonSucceeded means an operation attempt returned nil error.
	//
	// A successful retry execution may complete on the initial attempt or on a
	// later retry attempt. The reason does not encode which attempt succeeded;
	// that belongs to Outcome.Attempts and observer attempt metadata.
	StopReasonSucceeded StopReason = iota + 1

	// StopReasonNonRetryable means an operation attempt failed with an error that
	// the configured classifier did not allow to be retried.
	//
	// This reason does not mean retry created a wrapper error. Non-retryable
	// operation errors should normally be returned unchanged because they remain
	// operation-owned results rather than retry-owned failures.
	StopReasonNonRetryable

	// StopReasonMaxAttempts means retry execution stopped because the configured
	// maximum number of operation attempts had been reached.
	//
	// Max-attempt exhaustion is retry-owned exhaustion. It is expected to be
	// represented by ErrExhausted with the last operation error preserved as the
	// underlying cause when available.
	StopReasonMaxAttempts

	// StopReasonMaxElapsed means retry execution stopped because scheduling
	// another attempt would exceed the configured maximum elapsed runtime.
	//
	// Max-elapsed exhaustion is retry-owned exhaustion. The retry loop should
	// detect this before sleeping for a delay that cannot fit inside the elapsed
	// limit.
	StopReasonMaxElapsed

	// StopReasonBackoffExhausted means the configured backoff sequence had no
	// next delay for another retry attempt.
	//
	// The backoff package represents finite sequence exhaustion with ok=false.
	// That is not a backoff error. At the retry layer, however, it means a
	// retryable operation failure could not be followed by another scheduled
	// attempt, so the retry execution is exhausted.
	StopReasonBackoffExhausted

	// StopReasonInterrupted means retry execution stopped because the retry-owned
	// context was cancelled or its deadline expired.
	//
	// Interruption is owned by the retry loop only when retry observes the context
	// stop at its own boundary, such as before an attempt or while waiting between
	// attempts. Raw context errors returned by an operation remain operation-owned
	// unless the retry loop explicitly classifies its own context stop.
	StopReasonInterrupted
)

// String returns the canonical lower-case diagnostic name of r.
//
// The returned value is intended for diagnostics, tests, logs, observer events,
// and error messages. It is not a versioned serialization format. Unknown values
// return "invalid" so callers never accidentally render an unknown numeric value
// as a valid stop reason.
func (r StopReason) String() string {
	switch r {
	case StopReasonSucceeded:
		return "succeeded"
	case StopReasonNonRetryable:
		return "non_retryable"
	case StopReasonMaxAttempts:
		return "max_attempts"
	case StopReasonMaxElapsed:
		return "max_elapsed"
	case StopReasonBackoffExhausted:
		return "backoff_exhausted"
	case StopReasonInterrupted:
		return "interrupted"
	default:
		return "invalid"
	}
}

// IsValid reports whether r is one of the stop reasons defined by this package.
//
// IsValid is useful at package boundaries, in tests, and in defensive code that
// receives a StopReason from caller-controlled input. The zero value is invalid
// and must not be treated as success, failure, exhaustion, or interruption.
func (r StopReason) IsValid() bool {
	switch r {
	case StopReasonSucceeded,
		StopReasonNonRetryable,
		StopReasonMaxAttempts,
		StopReasonMaxElapsed,
		StopReasonBackoffExhausted,
		StopReasonInterrupted:
		return true
	default:
		return false
	}
}

// Succeeded reports whether r records a successful retry execution.
//
// Only StopReasonSucceeded is successful. Invalid reasons are not successful.
func (r StopReason) Succeeded() bool {
	return r == StopReasonSucceeded
}

// Failed reports whether r records an unsuccessful retry execution.
//
// Failed returns true for valid non-success reasons. Invalid reasons return
// false so malformed metadata is not accidentally classified as an ordinary
// retry failure.
func (r StopReason) Failed() bool {
	switch r {
	case StopReasonNonRetryable,
		StopReasonMaxAttempts,
		StopReasonMaxElapsed,
		StopReasonBackoffExhausted,
		StopReasonInterrupted:
		return true
	default:
		return false
	}
}

// Exhausted reports whether r records retry-owned exhaustion.
//
// Exhaustion means retry wanted or was allowed to continue only until a
// retry-owned boundary was reached: attempt limit, elapsed-time limit, or finite
// backoff sequence exhaustion. Non-retryable operation errors and retry-owned
// interruptions are not exhaustion.
func (r StopReason) Exhausted() bool {
	switch r {
	case StopReasonMaxAttempts,
		StopReasonMaxElapsed,
		StopReasonBackoffExhausted:
		return true
	default:
		return false
	}
}

// Interrupted reports whether r records retry-owned interruption.
//
// Interruption is distinct from exhaustion. It means retry execution stopped
// because the retry-owned context was cancelled or expired at a retry boundary.
func (r StopReason) Interrupted() bool {
	return r == StopReasonInterrupted
}

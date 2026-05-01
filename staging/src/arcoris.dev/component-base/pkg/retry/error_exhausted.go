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

import "errors"

const (
	panicInvalidExhaustedOutcome   = "retry: invalid exhausted outcome"
	panicNonExhaustedOutcomeReason = "retry: non-exhausted outcome reason"
)

// ErrExhausted identifies retry execution that stopped because a retry-owned
// exhaustion boundary was reached.
//
// Exhaustion is retry-owned. It is returned when retry execution cannot schedule
// another attempt because one of its configured retry boundaries has been
// reached:
//
//   - maximum operation attempts;
//   - maximum elapsed retry runtime;
//   - configured backoff sequence exhaustion.
//
// ErrExhausted is a classification sentinel. Retry implementations SHOULD return
// NewExhaustedError with a concrete Outcome instead of returning this sentinel
// directly when the exhaustion reason is known.
//
// Non-retryable operation errors are not ErrExhausted. They remain operation-
// owned errors and should normally be returned unchanged.
var ErrExhausted error = retryErrorKindExhausted

// Exhausted reports whether err identifies retry-owned exhaustion.
//
// Exhausted is a classification helper over errors.Is(err, ErrExhausted). It
// does not classify non-retryable operation errors, retry-owned context
// interruptions, or arbitrary caller errors as exhaustion.
func Exhausted(err error) bool {
	return errors.Is(err, ErrExhausted)
}

// NewExhaustedError returns an error classified as ErrExhausted.
//
// outcome must be valid retry completion metadata whose reason records
// retry-owned exhaustion: StopReasonMaxAttempts, StopReasonMaxElapsed, or
// StopReasonBackoffExhausted. Invalid outcomes and non-exhausted reasons are
// programming errors and cause panic.
//
// The returned error preserves outcome. Callers can recover it with
// ExhaustedOutcome. The returned error also unwraps to outcome.LastErr when it is
// non-nil, so errors.Is and errors.As can still inspect the last operation-owned
// error that caused retry execution to continue until exhaustion.
func NewExhaustedError(outcome Outcome) error {
	if !outcome.IsValid() {
		panic(panicInvalidExhaustedOutcome)
	}
	if !outcome.Exhausted() {
		panic(panicNonExhaustedOutcomeReason)
	}

	return exhaustedError{outcome: outcome}
}

// ExhaustedOutcome extracts retry exhaustion metadata from err.
//
// The helper searches the error chain for an error produced by
// NewExhaustedError. It returns false when err is nil, when err is not classified
// as retry-owned exhaustion, or when only the bare ErrExhausted sentinel is
// present without concrete Outcome metadata.
func ExhaustedOutcome(err error) (Outcome, bool) {
	var exhausted exhaustedError
	if errors.As(err, &exhausted) {
		return exhausted.outcome, true
	}

	return Outcome{}, false
}

// exhaustedError marks retry execution as exhausted while preserving completion
// metadata and the last operation-owned error.
//
// exhaustedError is intentionally private. Public callers should classify
// exhaustion with Exhausted or errors.Is(err, ErrExhausted), and should extract
// completion metadata with ExhaustedOutcome.
type exhaustedError struct {
	// outcome is the retry completion metadata that explains which retry-owned
	// exhaustion boundary stopped execution.
	outcome Outcome
}

// Error returns the exhaustion message.
//
// When Outcome.LastErr is available, the message includes both the retry-level
// classification and the underlying operation error. Error classification must
// still use errors.Is, not string matching.
func (e exhaustedError) Error() string {
	return retryErrorMessage(ErrExhausted, e.outcome.LastErr)
}

// Unwrap returns the last operation-owned error.
//
// A nil result means the exhausted outcome did not carry a lower-level operation
// error. Valid exhausted outcomes normally carry LastErr, but the method remains
// total and mirrors standard Go error wrapper behavior.
func (e exhaustedError) Unwrap() error {
	return e.outcome.LastErr
}

// Is reports whether target matches the retry exhaustion classification.
func (e exhaustedError) Is(target error) bool {
	return target == ErrExhausted
}

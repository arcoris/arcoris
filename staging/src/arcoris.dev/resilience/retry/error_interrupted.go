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

// ErrInterrupted identifies retry execution that stopped because retry observed
// its owning context stop at a retry boundary.
//
// Interruption is retry-owned only when retry itself observes the context stop,
// such as before an operation attempt or while waiting between attempts. Raw
// context.Canceled and context.DeadlineExceeded values returned by an operation
// remain operation-owned errors unless retry explicitly wraps its own context
// observation with NewInterruptedError.
//
// ErrInterrupted is a classification sentinel. Retry implementations SHOULD
// return NewInterruptedError with a concrete cause instead of returning this
// sentinel directly when the interruption reason is known.
var ErrInterrupted error = retryErrorKindInterrupted

// Interrupted reports whether err identifies retry-owned interruption.
//
// Interrupted is a classification helper over errors.Is(err, ErrInterrupted). It
// intentionally does not treat raw context.Canceled or context.DeadlineExceeded
// as retry interruptions. Operations can return those errors for their own work,
// and retry must not infer ownership unless it observed the context stop itself.
func Interrupted(err error) bool {
	return errors.Is(err, ErrInterrupted)
}

// NewInterruptedError returns an error classified as ErrInterrupted.
//
// cause describes why retry execution was interrupted. The returned error unwraps
// to cause, so callers can still use errors.Is and errors.As for the underlying
// reason, including context.Canceled, context.DeadlineExceeded, or a custom
// cancellation cause.
//
// If cause is nil, NewInterruptedError still returns a non-nil error classified
// as ErrInterrupted.
//
// If cause is already classified as ErrInterrupted, it is returned unchanged to
// avoid adding redundant retry-level wrappers.
func NewInterruptedError(cause error) error {
	if cause == nil {
		return interruptedError{}
	}
	if errors.Is(cause, ErrInterrupted) {
		return cause
	}

	return interruptedError{cause: cause}
}

// interruptedError marks retry execution as interrupted while preserving the
// lower-level cause.
//
// interruptedError is intentionally private. Public callers should classify retry
// interruption with Interrupted or errors.Is(err, ErrInterrupted), not by
// depending on concrete error types.
type interruptedError struct {
	// cause is the lower-level reason retry execution was interrupted.
	//
	// The value MAY be nil when retry knows execution was interrupted but has no
	// more specific cause to expose.
	cause error
}

// Error returns the interruption message.
//
// When a concrete cause is available, the message includes both the retry-level
// classification and the underlying cause. Error classification must still use
// errors.Is, not string matching.
func (e interruptedError) Error() string {
	return retryErrorMessage(ErrInterrupted, e.cause)
}

// Unwrap returns the lower-level interruption cause.
//
// A nil result means the interruption has no lower-level cause. This still
// remains a valid ErrInterrupted classification through interruptedError.Is.
func (e interruptedError) Unwrap() error {
	return e.cause
}

// Is reports whether target matches the retry interruption classification.
func (e interruptedError) Is(target error) bool {
	return target == ErrInterrupted
}

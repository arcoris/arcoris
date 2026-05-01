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

package wait

import "errors"

// ErrInterrupted identifies a wait operation that stopped before its condition
// completed successfully and without the condition returning its own domain
// error.
//
// ErrInterrupted is a classification sentinel. It describes wait-loop ownership
// of the stop reason; it does not describe arbitrary condition failures. Raw
// context.Canceled and context.DeadlineExceeded values are not classified as
// wait interruptions unless a wait implementation wraps them with
// NewInterruptedError or NewTimeoutError.
//
// Wait implementations SHOULD return NewInterruptedError with a concrete cause
// instead of returning this sentinel directly when the interruption reason is
// known.
//
// Timeout is a specialized form of interruption. Errors classified as ErrTimeout
// also match ErrInterrupted through errors.Is.
var ErrInterrupted error = waitErrorKindInterrupted

// Interrupted reports whether err identifies a wait-owned interruption.
//
// Interrupted is a classification helper over errors.Is(err, ErrInterrupted). It
// reports true for timeout errors because timeout is a specialized form of
// interruption.
//
// Interrupted intentionally does not treat raw context.Canceled or
// context.DeadlineExceeded as wait interruptions. A condition can return those
// errors for condition-owned work, and this package cannot infer that such
// errors came from the wait loop itself.
//
// Wait implementations that stop because their own context was cancelled SHOULD
// return NewInterruptedError(ctx.Err()). Wait implementations that stop because
// their own deadline, timeout, or finite retry budget expired SHOULD return
// NewTimeoutError(cause).
func Interrupted(err error) bool {
	return errors.Is(err, ErrInterrupted)
}

// NewInterruptedError returns an error classified as ErrInterrupted.
//
// cause describes why the wait operation was interrupted. The returned error
// unwraps to cause, so callers can still use errors.Is and errors.As for the
// underlying reason.
//
// If cause is nil, NewInterruptedError still returns a non-nil error classified
// as ErrInterrupted.
//
// If cause is already classified as ErrInterrupted, it is returned unchanged to
// avoid adding redundant wait-level wrappers. This includes timeout errors,
// because timeout is a specialized form of interruption.
func NewInterruptedError(cause error) error {
	if cause == nil {
		return interruptedError{}
	}
	if errors.Is(cause, ErrInterrupted) {
		return cause
	}

	return interruptedError{cause: cause}
}

// interruptedError marks a wait operation as interrupted while preserving the
// lower-level cause.
//
// interruptedError is intentionally private. Public callers should classify wait
// interruption with Interrupted or errors.Is(err, ErrInterrupted), not by
// depending on concrete error types.
type interruptedError struct {
	// cause is the lower-level reason the wait operation was interrupted.
	//
	// The value MAY be nil when the wait layer knows that the operation was
	// interrupted but has no more specific cause to expose.
	cause error
}

// Error returns the interruption message.
//
// When a concrete cause is available, the message includes both the wait-level
// classification and the underlying cause. Error classification must still use
// errors.Is, not string matching.
func (e interruptedError) Error() string {
	return waitErrorMessage(ErrInterrupted, e.cause)
}

// Unwrap returns the lower-level interruption cause.
//
// A nil result means the interruption has no lower-level cause. This still
// remains a valid ErrInterrupted classification through interruptedError.Is.
func (e interruptedError) Unwrap() error {
	return e.cause
}

// Is reports whether target matches the wait interruption classification.
//
// interruptedError matches only ErrInterrupted. It does not match ErrTimeout
// because not every wait interruption is a timeout.
func (e interruptedError) Is(target error) bool {
	return target == ErrInterrupted
}

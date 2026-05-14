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

// ErrTimeout identifies a wait operation whose condition was not satisfied
// before the wait-owned timeout, deadline, or finite retry budget expired.
//
// ErrTimeout is a classification sentinel. It describes timeout ownership by the
// wait operation itself. A condition can return context.DeadlineExceeded for its
// own work; that raw error is not a wait timeout unless the wait loop explicitly
// wraps it with NewTimeoutError.
//
// Wait implementations SHOULD return NewTimeoutError with a concrete cause
// instead of returning this sentinel directly when the timeout reason is known.
//
// ErrTimeout is more specific than ErrInterrupted: errors.Is(ErrTimeout,
// ErrInterrupted) is true.
var ErrTimeout error = waitErrorKindTimeout

// TimedOut reports whether err identifies a wait-owned timeout.
//
// TimedOut is a classification helper over errors.Is(err, ErrTimeout).
//
// TimedOut intentionally does not treat raw context.DeadlineExceeded as a wait
// timeout. A condition can return context.DeadlineExceeded for its own operation,
// and callers must not classify that as a wait-loop timeout unless the wait loop
// explicitly wrapped it with NewTimeoutError.
func TimedOut(err error) bool {
	return errors.Is(err, ErrTimeout)
}

// NewTimeoutError returns an error classified as ErrTimeout.
//
// Timeout is modeled as a specialized interruption. Errors returned by
// NewTimeoutError match both ErrTimeout and ErrInterrupted through errors.Is.
//
// cause describes why the wait timed out or which lower-level operation observed
// the timeout. The returned error unwraps to cause, so callers can still use
// errors.Is and errors.As for the underlying reason.
//
// If cause is nil, NewTimeoutError still returns a non-nil error classified as
// both ErrTimeout and ErrInterrupted.
//
// If cause is already classified as ErrTimeout, it is returned unchanged to avoid
// adding redundant timeout wrappers.
func NewTimeoutError(cause error) error {
	if cause == nil {
		return timeoutError{}
	}
	if errors.Is(cause, ErrTimeout) {
		return cause
	}

	return timeoutError{cause: cause}
}

// timeoutError marks a wait operation as timed out while preserving the
// lower-level cause.
//
// timeoutError is intentionally private. Public callers should classify timeout
// with TimedOut or errors.Is(err, ErrTimeout), not by depending on concrete error
// types.
type timeoutError struct {
	// cause is the lower-level reason the wait operation timed out.
	//
	// The value MAY be nil when the wait layer knows that the operation timed out
	// but has no more specific cause to expose.
	cause error
}

// Error returns the timeout message.
//
// When a concrete cause is available, the message includes both the wait-level
// classification and the underlying cause. Error classification must still use
// errors.Is, not string matching.
func (e timeoutError) Error() string {
	return waitErrorMessage(ErrTimeout, e.cause)
}

// Unwrap returns the lower-level timeout cause.
//
// A nil result means the timeout has no lower-level cause. This still remains a
// valid ErrTimeout and ErrInterrupted classification through timeoutError.Is.
func (e timeoutError) Unwrap() error {
	return e.cause
}

// Is reports whether target matches the wait timeout or the more general wait
// interruption classification.
//
// Timeout is a specialized interruption, so timeoutError matches both ErrTimeout
// and ErrInterrupted.
func (e timeoutError) Is(target error) bool {
	return target == ErrTimeout || target == ErrInterrupted
}

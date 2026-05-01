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

const (
	// errInterruptedMessage is the stable diagnostic text for wait-owned
	// interruption errors.
	//
	// The value is intentionally shared by the interruption sentinel and by
	// concrete interruption wrappers so callers see one consistent wait-level
	// classification message.
	errInterruptedMessage = "wait: interrupted"

	// errTimeoutMessage is the stable diagnostic text for wait-owned timeout
	// errors.
	//
	// The value is intentionally shared by the timeout sentinel and by concrete
	// timeout wrappers so callers see one consistent wait-level classification
	// message.
	errTimeoutMessage = "wait: timeout"
)

// waitErrorKind is the private concrete type behind wait classification
// sentinels.
//
// The type exists because the wait error model has hierarchy: timeout is a more
// specific form of interruption. A plain errors.New sentinel cannot express that
// hierarchy by itself, while a private type with Is can.
type waitErrorKind uint8

const (
	// waitErrorKindInterrupted identifies the general wait interruption sentinel.
	//
	// It is kept private so public callers classify interruptions through
	// ErrInterrupted, Interrupted, or errors.Is rather than depending on concrete
	// sentinel implementation details.
	waitErrorKindInterrupted waitErrorKind = iota + 1

	// waitErrorKindTimeout identifies the wait timeout sentinel.
	//
	// It is kept private so public callers classify timeouts through ErrTimeout,
	// TimedOut, or errors.Is rather than depending on concrete sentinel
	// implementation details.
	waitErrorKindTimeout
)

// Error returns the stable wait classification message for k.
//
// Unknown values are not expected in normal package use, but the method remains
// total so a zero or corrupted private value still produces a diagnostic string
// instead of an empty message.
func (k waitErrorKind) Error() string {
	switch k {
	case waitErrorKindInterrupted:
		return errInterruptedMessage
	case waitErrorKindTimeout:
		return errTimeoutMessage
	default:
		return "wait: unknown error"
	}
}

// Is reports whether target matches the classification represented by k.
//
// Timeout intentionally matches both ErrTimeout and ErrInterrupted because every
// wait-owned timeout is also a wait-owned interruption. The reverse is not true:
// a generic interruption is not necessarily a timeout.
func (k waitErrorKind) Is(target error) bool {
	switch k {
	case waitErrorKindInterrupted:
		return target == ErrInterrupted
	case waitErrorKindTimeout:
		return target == ErrTimeout || target == ErrInterrupted
	default:
		return false
	}
}

// waitErrorMessage builds a stable wait-level error message while preserving
// cause text for diagnostics.
//
// The returned string is for humans and logs only. Callers MUST classify wait
// errors with errors.Is, Interrupted, or TimedOut instead of matching the message.
func waitErrorMessage(kind error, cause error) string {
	if cause == nil {
		return kind.Error()
	}

	return kind.Error() + ": " + cause.Error()
}

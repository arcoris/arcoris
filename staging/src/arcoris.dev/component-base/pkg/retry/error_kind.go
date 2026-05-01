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

const (
	// errExhaustedMessage is the stable diagnostic text for retry-owned
	// exhaustion errors.
	//
	// The value is intentionally shared by the exhaustion sentinel and by
	// concrete exhaustion wrappers so callers see one consistent retry-level
	// classification message.
	errExhaustedMessage = "retry: exhausted"

	// errInterruptedMessage is the stable diagnostic text for retry-owned
	// interruption errors.
	//
	// The value is intentionally shared by the interruption sentinel and by
	// concrete interruption wrappers so callers see one consistent retry-level
	// classification message.
	errInterruptedMessage = "retry: interrupted"
)

// retryErrorKind is the private concrete type behind retry classification
// sentinels.
//
// The type exists so public callers can classify retry-owned failures through
// ErrExhausted, ErrInterrupted, Exhausted, Interrupted, or errors.Is without
// depending on concrete implementation details.
type retryErrorKind uint8

const (
	// retryErrorKindExhausted identifies retry-owned exhaustion.
	//
	// Exhaustion means retry execution stopped because a retry-owned boundary was
	// reached: max attempts, max elapsed time, or finite backoff sequence
	// exhaustion.
	retryErrorKindExhausted retryErrorKind = iota + 1

	// retryErrorKindInterrupted identifies retry-owned interruption.
	//
	// Interruption means retry execution stopped because retry observed its owning
	// context stop at a retry boundary, such as before an attempt or while waiting
	// between attempts.
	retryErrorKindInterrupted
)

// Error returns the stable retry classification message for k.
//
// Unknown values are not expected in normal package use, but the method remains
// total so a zero or corrupted private value still produces a diagnostic string
// instead of an empty message.
func (k retryErrorKind) Error() string {
	switch k {
	case retryErrorKindExhausted:
		return errExhaustedMessage
	case retryErrorKindInterrupted:
		return errInterruptedMessage
	default:
		return "retry: unknown error"
	}
}

// Is reports whether target matches the classification represented by k.
func (k retryErrorKind) Is(target error) bool {
	switch k {
	case retryErrorKindExhausted:
		return target == ErrExhausted
	case retryErrorKindInterrupted:
		return target == ErrInterrupted
	default:
		return false
	}
}

// retryErrorMessage builds a stable retry-level error message while preserving
// cause text for diagnostics.
//
// The returned string is for humans and logs only. Callers MUST classify retry
// errors with errors.Is, Exhausted, or Interrupted instead of matching the
// message.
func retryErrorMessage(kind error, cause error) string {
	if cause == nil {
		return kind.Error()
	}

	return kind.Error() + ": " + cause.Error()
}

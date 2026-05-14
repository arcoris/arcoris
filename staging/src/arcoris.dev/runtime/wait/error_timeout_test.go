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

import (
	"context"
	"errors"
	"testing"
)

// TestTimedOutClassification verifies that TimedOut reports only wait-owned
// timeout classifications.
func TestTimedOutClassification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "interrupted sentinel",
			err:  ErrInterrupted,
			want: false,
		},
		{
			name: "interrupted wrapper",
			err:  NewInterruptedError(context.Canceled),
			want: false,
		},
		{
			name: "timeout sentinel",
			err:  ErrTimeout,
			want: true,
		},
		{
			name: "timeout wrapper",
			err:  NewTimeoutError(context.DeadlineExceeded),
			want: true,
		},
		{
			name: "raw context deadline exceeded",
			err:  context.DeadlineExceeded,
			want: false,
		},
		{
			name: "ordinary condition error",
			err:  errors.New("condition failed"),
			want: false,
		},
		{
			name: "joined ordinary and timeout",
			err:  errors.Join(errors.New("condition failed"), NewTimeoutError(context.DeadlineExceeded)),
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := TimedOut(tc.err); got != tc.want {
				t.Fatalf("TimedOut(err) = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestNewTimeoutErrorWithCause verifies timeout wrapping with a concrete
// lower-level cause.
func TestNewTimeoutErrorWithCause(t *testing.T) {
	t.Parallel()

	cause := context.DeadlineExceeded
	err := NewTimeoutError(cause)

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustMatch(t, err, cause)
	mustUnwrapTo(t, err, cause)
	mustHaveMessage(t, err, "wait: timeout: context deadline exceeded")
}

// TestNewTimeoutErrorWithNilCause verifies that a nil cause still produces a
// classified, non-nil timeout error.
func TestNewTimeoutErrorWithNilCause(t *testing.T) {
	t.Parallel()

	err := NewTimeoutError(nil)

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustUnwrapTo(t, err, nil)
	mustHaveMessage(t, err, "wait: timeout")
}

// TestNewTimeoutErrorPreservesAlreadyTimedOutErrors verifies idempotent wrapping
// for existing wait-owned timeouts.
func TestNewTimeoutErrorPreservesAlreadyTimedOutErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "timeout sentinel",
			err:  ErrTimeout,
		},
		{
			name: "timeout wrapper",
			err:  NewTimeoutError(context.DeadlineExceeded),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := NewTimeoutError(tc.err); got != tc.err {
				t.Fatal("NewTimeoutError wrapped an error that was already classified as timeout")
			}
		})
	}
}

// TestNewTimeoutErrorCanWrapInterruptedCause verifies that a timeout can expose
// an interruption as its lower-level cause while remaining a timeout.
func TestNewTimeoutErrorCanWrapInterruptedCause(t *testing.T) {
	t.Parallel()

	cause := NewInterruptedError(context.Canceled)
	err := NewTimeoutError(cause)

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustMatch(t, err, context.Canceled)
	mustUnwrapTo(t, err, cause)
	mustHaveMessage(t, err, "wait: timeout: wait: interrupted: context canceled")
}

// TestWrappedTimeoutErrorRemainsClassified verifies timeout classification
// through ordinary Go error wrapping.
func TestWrappedTimeoutErrorRemainsClassified(t *testing.T) {
	t.Parallel()

	err := wrapForTest(NewTimeoutError(context.DeadlineExceeded))

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
}

// TestTimeoutErrorPreservesAsTraversal verifies that timeout wrappers preserve
// errors.As traversal for typed causes.
func TestTimeoutErrorPreservesAsTraversal(t *testing.T) {
	t.Parallel()

	cause := typedCause{message: "typed cause"}
	err := NewTimeoutError(cause)

	mustAsTypedCause(t, err, cause)
}

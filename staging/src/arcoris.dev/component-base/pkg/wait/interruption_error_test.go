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

// TestInterruptedClassification verifies that Interrupted reports only
// wait-owned interruption classifications.
func TestInterruptedClassification(t *testing.T) {
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
			want: true,
		},
		{
			name: "interrupted wrapper",
			err:  NewInterruptedError(context.Canceled),
			want: true,
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
			name: "raw context canceled",
			err:  context.Canceled,
			want: false,
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
			name: "joined ordinary and interrupted",
			err:  errors.Join(errors.New("condition failed"), NewInterruptedError(context.Canceled)),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Interrupted(tt.err); got != tt.want {
				t.Fatalf("Interrupted(err) = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewInterruptedErrorWithCause verifies interruption wrapping with a concrete
// lower-level cause.
func TestNewInterruptedErrorWithCause(t *testing.T) {
	t.Parallel()

	cause := context.Canceled
	err := NewInterruptedError(cause)

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, cause)
	mustUnwrapTo(t, err, cause)
	mustHaveMessage(t, err, "wait: interrupted: context canceled")
}

// TestNewInterruptedErrorWithNilCause verifies that a nil cause still produces a
// classified, non-nil interruption error.
func TestNewInterruptedErrorWithNilCause(t *testing.T) {
	t.Parallel()

	err := NewInterruptedError(nil)

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustUnwrapTo(t, err, nil)
	mustHaveMessage(t, err, "wait: interrupted")
}

// TestNewInterruptedErrorPreservesAlreadyInterruptedErrors verifies idempotent
// wrapping for existing wait-owned interruptions.
func TestNewInterruptedErrorPreservesAlreadyInterruptedErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "interrupted sentinel",
			err:  ErrInterrupted,
		},
		{
			name: "interrupted wrapper",
			err:  NewInterruptedError(context.Canceled),
		},
		{
			name: "timeout sentinel",
			err:  ErrTimeout,
		},
		{
			name: "timeout wrapper",
			err:  NewTimeoutError(context.DeadlineExceeded),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewInterruptedError(tt.err); got != tt.err {
				t.Fatal("NewInterruptedError wrapped an error that was already classified as interrupted")
			}
		})
	}
}

// TestWrappedInterruptedErrorRemainsClassified verifies interruption
// classification through ordinary Go error wrapping.
func TestWrappedInterruptedErrorRemainsClassified(t *testing.T) {
	t.Parallel()

	err := wrapForTest(NewInterruptedError(context.Canceled))

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
}

// TestInterruptedErrorPreservesAsTraversal verifies that interruption wrappers
// preserve errors.As traversal for typed causes.
func TestInterruptedErrorPreservesAsTraversal(t *testing.T) {
	t.Parallel()

	cause := typedCause{message: "typed cause"}
	err := NewInterruptedError(cause)

	mustAsTypedCause(t, err, cause)
}

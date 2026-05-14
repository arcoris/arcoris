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

package lifecycle

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestWaitErrorErrorForTargetWait(t *testing.T) {
	t.Parallel()

	// HasTarget is explicit because StateNew is a valid zero value and cannot
	// double as an absent-target marker.
	err := &WaitError{Snapshot: Snapshot{State: StateStarting, Revision: 1}, Target: StateRunning, HasTarget: true, Err: ErrWaitTargetUnreachable}
	if got, want := err.Error(), "lifecycle: wait target unreachable: target running at starting@1"; got != want {
		t.Fatalf("WaitError.Error() = %q, want %q", got, want)
	}
}

func TestWaitErrorErrorForGenericPredicateWait(t *testing.T) {
	t.Parallel()

	err := &WaitError{Snapshot: Snapshot{State: StateStarting, Revision: 1}, Err: ErrInvalidWaitPredicate}
	if got, want := err.Error(), "lifecycle: invalid wait predicate: at starting@1"; got != want {
		t.Fatalf("WaitError.Error() = %q, want %q", got, want)
	}
}

func TestWaitErrorNilReceiver(t *testing.T) {
	t.Parallel()

	var err *WaitError
	if got := err.Error(); got != ErrWaitTargetUnreachable.Error() {
		t.Fatalf("nil Error() = %q, want %q", got, ErrWaitTargetUnreachable.Error())
	}
	mustMatch(t, err, ErrWaitTargetUnreachable)
}

func TestWaitErrorUnwrap(t *testing.T) {
	t.Parallel()

	err := &WaitError{Err: context.Canceled}
	if got := err.Unwrap(); got != context.Canceled {
		t.Fatalf("Unwrap = %v, want context.Canceled", got)
	}
}

func TestNewWaitErrorDefaultsNilCause(t *testing.T) {
	t.Parallel()

	err := newWaitError(Snapshot{}, nil)
	if err.Err != ErrWaitTargetUnreachable {
		t.Fatalf("Err = %v, want unreachable", err.Err)
	}
}

func TestNewWaitStateErrorSetsTarget(t *testing.T) {
	t.Parallel()

	err := newWaitStateError(Snapshot{}, StateNew, nil)
	if !err.HasTarget || err.Target != StateNew {
		t.Fatalf("target = %s has=%v, want StateNew with target", err.Target, err.HasTarget)
	}
	if err.Err != ErrWaitTargetUnreachable {
		t.Fatalf("Err = %v, want unreachable", err.Err)
	}
}

func TestWaitErrorCauseNilBehavior(t *testing.T) {
	t.Parallel()

	if got := waitErrorCause(nil); got != ErrWaitTargetUnreachable {
		t.Fatalf("nil waitErrorCause = %v, want unreachable", got)
	}
	if got := waitErrorCause(&WaitError{}); got != ErrWaitTargetUnreachable {
		t.Fatalf("empty waitErrorCause = %v, want unreachable", got)
	}
}

func TestWaitErrorMessagePrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		cause error
		want  string
	}{
		{ErrWaitTargetUnreachable, ErrWaitTargetUnreachable.Error()},
		{ErrInvalidWaitPredicate, ErrInvalidWaitPredicate.Error()},
		{ErrInvalidWaitTarget, ErrInvalidWaitTarget.Error()},
		{context.Canceled, "lifecycle: wait failed: context canceled"},
		{context.DeadlineExceeded, "lifecycle: wait failed: context deadline exceeded"},
	}

	for _, tc := range tests {
		if got := waitErrorMessagePrefix(tc.cause); got != tc.want {
			t.Fatalf("prefix = %q, want %q", got, tc.want)
		}
		if !strings.HasPrefix(waitErrorMessagePrefix(tc.cause), "lifecycle:") {
			t.Fatalf("prefix %q lacks lifecycle prefix", waitErrorMessagePrefix(tc.cause))
		}
	}
}

func TestWaitErrorMatchesContextAndLifecycleSentinels(t *testing.T) {
	t.Parallel()

	for _, cause := range []error{
		context.Canceled,
		context.DeadlineExceeded,
		ErrWaitTargetUnreachable,
		ErrInvalidWaitPredicate,
		ErrInvalidWaitTarget,
	} {
		err := &WaitError{Err: cause}
		if !errors.Is(err, cause) {
			t.Fatalf("errors.Is(%v, %v) = false, want true", err, cause)
		}
	}
}

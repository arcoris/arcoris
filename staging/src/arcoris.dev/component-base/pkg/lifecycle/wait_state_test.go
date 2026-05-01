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
	"testing"
)

func TestWaitStateImmediateSuccess(t *testing.T) {
	t.Parallel()

	snapshot, err := NewController().WaitState(context.Background(), StateNew)
	if err != nil {
		t.Fatalf("WaitState current = %v", err)
	}
	if snapshot.State != StateNew {
		t.Fatalf("snapshot.State = %s, want new", snapshot.State)
	}
}

func TestWaitStateRejectsInvalidTarget(t *testing.T) {
	t.Parallel()

	snapshot, err := NewController().WaitState(context.Background(), State(99))
	if err == nil {
		t.Fatal("WaitState invalid target err = nil, want error")
	}
	mustMatch(t, err, ErrInvalidWaitTarget)
	if snapshot.State != StateNew {
		t.Fatalf("snapshot.State = %s, want latest new", snapshot.State)
	}
}

func TestWaitStateRejectsUnreachableTargetBeforeBlocking(t *testing.T) {
	t.Parallel()

	// WaitState uses static graph reachability, not guard-dependent runtime
	// progress, to reject targets that can no longer be reached.
	controller := NewController()
	_, _ = controller.BeginStart()
	snapshot, err := controller.WaitState(context.Background(), StateNew)
	if err == nil {
		t.Fatal("WaitState backward err = nil, want unreachable")
	}
	mustMatch(t, err, ErrWaitTargetUnreachable)
	if snapshot.State != StateStarting {
		t.Fatalf("snapshot.State = %s, want starting", snapshot.State)
	}
}

func TestWaitStateWaitsUntilReachableTargetCommits(t *testing.T) {
	t.Parallel()

	controller := NewController()
	results := make(chan Snapshot, 1)
	errs := make(chan error, 1)
	go func() {
		snapshot, err := controller.WaitState(context.Background(), StateRunning)
		if err != nil {
			errs <- err
			return
		}
		results <- snapshot
	}()

	_, _ = controller.BeginStart()
	_, _ = controller.MarkRunning()
	select {
	case err := <-errs:
		t.Fatalf("WaitState err = %v, want nil", err)
	default:
	}
	if got := mustReceiveSnapshot(t, results); got.State != StateRunning {
		t.Fatalf("snapshot.State = %s, want running", got.State)
	}
}

func TestWaitStateTerminalBeforeTarget(t *testing.T) {
	t.Parallel()

	controller := NewController()
	_, _ = controller.BeginStop()
	snapshot, err := controller.WaitState(context.Background(), StateRunning)
	if err == nil {
		t.Fatal("WaitState err = nil, want unreachable")
	}
	mustMatch(t, err, ErrWaitTargetUnreachable)
	if snapshot.State != StateStopped {
		t.Fatalf("snapshot.State = %s, want stopped", snapshot.State)
	}
}

func TestWaitStateDoneBranchReturnsUnreachableForOtherTarget(t *testing.T) {
	t.Parallel()

	// If the terminal done snapshot is not the target, graph reachability cannot
	// improve afterward and the wait must fail as unreachable.
	changed := make(chan struct{})
	done := make(chan struct{})
	controller := &Controller{state: StateStopping, changed: changed, done: done}
	close(done)

	_, err := controller.WaitState(context.Background(), StateFailed)
	mustMatch(t, err, ErrWaitTargetUnreachable)
}

func TestWaitStateContextCancellationAndDeadline(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	snapshot, err := NewController().WaitState(ctx, StateRunning)
	if err == nil {
		t.Fatal("WaitState cancel err = nil, want canceled")
	}
	mustMatch(t, err, context.Canceled)
	if snapshot.State != StateNew {
		t.Fatalf("snapshot.State = %s, want new", snapshot.State)
	}

	deadlineCtx, deadlineCancel := context.WithDeadline(context.Background(), testTime)
	defer deadlineCancel()
	_, deadlineErr := NewController().WaitState(deadlineCtx, StateRunning)
	if deadlineErr == nil {
		t.Fatal("WaitState deadline err = nil, want deadline exceeded")
	}
	mustMatch(t, deadlineErr, context.DeadlineExceeded)
}

func TestWaitStateNilContextBehavesAsBackground(t *testing.T) {
	t.Parallel()

	snapshot, err := NewController().WaitState(nil, StateNew)
	if err != nil {
		t.Fatalf("WaitState nil context = %v, want nil", err)
	}
	if snapshot.State != StateNew {
		t.Fatalf("snapshot.State = %s, want new", snapshot.State)
	}
}

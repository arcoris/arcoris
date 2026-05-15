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

func TestWaitImmediateSuccess(t *testing.T) {
	t.Parallel()

	snap, err := NewController().Wait(context.Background(), func(snap Snapshot) bool {
		return snap.State == StateNew
	})
	if err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
	if snap.State != StateNew {
		t.Fatalf("snapshot.State = %s, want new", snap.State)
	}
}

func TestWaitNilContextBehavesAsBackground(t *testing.T) {
	t.Parallel()

	snap, err := NewController().Wait(nil, func(snap Snapshot) bool {
		return snap.State == StateNew
	})
	if err != nil {
		t.Fatalf("Wait nil context = %v, want nil", err)
	}
	if snap.State != StateNew {
		t.Fatalf("snapshot.State = %s, want new", snap.State)
	}
}

func TestWaitNilPredicate(t *testing.T) {
	t.Parallel()

	snap, err := NewController().Wait(context.Background(), nil)
	if err == nil {
		t.Fatal("Wait nil predicate err = nil, want error")
	}
	mustMatch(t, err, ErrInvalidWaitPredicate)
	if snap.State != StateNew {
		t.Fatalf("snapshot.State = %s, want latest new", snap.State)
	}
}

func TestWaitAcrossChangedSignal(t *testing.T) {
	t.Parallel()

	// Wait must re-evaluate predicates after a changed signal because only the
	// next snapshot can satisfy conditions based on committed lifecycle progress.
	controller := NewController()
	firstEval := make(chan struct{})
	results := make(chan Snapshot, 1)
	errs := make(chan error, 1)
	go func() {
		snap, err := controller.Wait(context.Background(), func(snap Snapshot) bool {
			if snap.State == StateNew {
				select {
				case firstEval <- struct{}{}:
				default:
				}
			}
			return snap.State == StateStarting
		})
		if err != nil {
			errs <- err
			return
		}
		results <- snap
	}()

	mustSignalClosed(t, firstEval)
	_, _ = controller.BeginStart()
	select {
	case err := <-errs:
		t.Fatalf("Wait err = %v, want nil", err)
	default:
	}
	if got := mustReceiveSnapshot(t, results); got.State != StateStarting {
		t.Fatalf("snapshot.State = %s, want starting", got.State)
	}
}

func TestWaitReturnsUnreachableAtTerminalBoundary(t *testing.T) {
	t.Parallel()

	// Terminal states have no outgoing transitions; if the predicate is false at
	// the terminal boundary, it cannot become true later.
	controller := NewController()
	_, _ = controller.BeginStop()
	snap, err := controller.Wait(context.Background(), func(Snapshot) bool { return false })
	if err == nil {
		t.Fatal("Wait err = nil, want unreachable")
	}
	mustMatch(t, err, ErrWaitTargetUnreachable)
	if snap.State != StateStopped {
		t.Fatalf("snapshot.State = %s, want stopped", snap.State)
	}
}

func TestWaitContextCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	snap, err := NewController().Wait(ctx, func(Snapshot) bool { return false })
	if err == nil {
		t.Fatal("Wait err = nil, want canceled")
	}
	mustMatch(t, err, context.Canceled)
	if snap.State != StateNew {
		t.Fatalf("snapshot.State = %s, want latest new", snap.State)
	}
}

func TestWaitContextDeadlineExceeded(t *testing.T) {
	t.Parallel()

	deadlineCtx, deadlineCancel := context.WithDeadline(context.Background(), testTime)
	defer deadlineCancel()
	deadlineSnapshot, deadlineErr := NewController().Wait(deadlineCtx, func(Snapshot) bool { return false })
	if deadlineErr == nil {
		t.Fatal("Wait deadline err = nil, want deadline exceeded")
	}
	mustMatch(t, deadlineErr, context.DeadlineExceeded)
	if deadlineSnapshot.State != StateNew {
		t.Fatalf("deadline snapshot.State = %s, want latest new", deadlineSnapshot.State)
	}
}

func TestWaitReturnsLatestSnapshotOnCancellation(t *testing.T) {
	t.Parallel()

	controller := NewController()
	_, _ = controller.BeginStart()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	snap, err := controller.Wait(ctx, func(Snapshot) bool { return false })
	if err == nil {
		t.Fatal("Wait err = nil, want canceled")
	}
	mustMatch(t, err, context.Canceled)
	if snap.State != StateStarting {
		t.Fatalf("snapshot.State = %s, want starting", snap.State)
	}
}

func TestWaitDoneBranchReturnsUnreachableWhenPredicateStillFalse(t *testing.T) {
	t.Parallel()

	// The done branch reports unreachable when the final terminal snapshot still
	// does not satisfy the predicate.
	changed := make(chan struct{})
	done := make(chan struct{})
	controller := &Controller{
		state:   StateStarting,
		changed: changed,
		done:    done,
	}
	firstEval := make(chan struct{})
	errs := make(chan error, 1)
	go func() {
		_, err := controller.Wait(context.Background(), func(snap Snapshot) bool {
			if snap.State == StateStarting {
				select {
				case firstEval <- struct{}{}:
				default:
				}
			}
			return false
		})
		errs <- err
	}()

	mustSignalClosed(t, firstEval)
	controller.mu.Lock()
	controller.state = StateStopped
	controller.revision = 1
	controller.lastTransition = Transition{From: StateStarting, To: StateStopped, Event: EventMarkStopped, Revision: 1, At: testTime}
	close(done)
	controller.mu.Unlock()

	err := mustReceiveError(t, errs)
	mustMatch(t, err, ErrWaitTargetUnreachable)
}

func TestWaitDoneBranchCanSucceedForTerminalSnapshot(t *testing.T) {
	t.Parallel()

	// The done signal is a terminal boundary, but Wait still gives the predicate
	// one final evaluation so terminal predicates can succeed.
	changed := make(chan struct{})
	done := make(chan struct{})
	controller := &Controller{
		state:   StateStarting,
		changed: changed,
		done:    done,
	}
	results := make(chan Snapshot, 1)
	errs := make(chan error, 1)
	firstEval := make(chan struct{})
	go func() {
		snap, err := controller.Wait(context.Background(), func(snap Snapshot) bool {
			if snap.State == StateStarting {
				select {
				case firstEval <- struct{}{}:
				default:
				}
			}
			return snap.State == StateStopped
		})
		if err != nil {
			errs <- err
			return
		}
		results <- snap
	}()

	mustSignalClosed(t, firstEval)
	controller.mu.Lock()
	controller.state = StateStopped
	controller.revision = 1
	controller.lastTransition = Transition{From: StateStarting, To: StateStopped, Event: EventMarkStopped, Revision: 1, At: testTime}
	close(done)
	controller.mu.Unlock()

	select {
	case err := <-errs:
		t.Fatalf("Wait err = %v, want nil", err)
	default:
	}
	if got := mustReceiveSnapshot(t, results); got.State != StateStopped {
		t.Fatalf("snapshot.State = %s, want stopped", got.State)
	}
}

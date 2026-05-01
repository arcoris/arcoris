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

func TestControllerSignalZeroValueInitializesThroughPublicOperations(t *testing.T) {
	t.Parallel()

	// Zero-value controllers lazily create signals on first operation so embedded
	// lifecycle fields remain usable without constructor-only setup.
	var controller Controller
	done := controller.Done()
	if done == nil {
		t.Fatal("Done returned nil for zero-value Controller")
	}
	_, _ = controller.BeginStop()
	mustSignalClosed(t, done)
}

func TestControllerSignalDoneChannelStable(t *testing.T) {
	t.Parallel()

	controller := NewController()
	first := controller.Done()
	second := controller.Done()
	if first != second {
		t.Fatal("Done returned different channels")
	}
}

func TestControllerSignalDoneClosesExactlyOnceOnTerminalTransition(t *testing.T) {
	t.Parallel()

	controller := NewController()
	done := controller.Done()
	_, _ = controller.BeginStop()
	mustSignalClosed(t, done)
	mustSignalClosed(t, controller.Done())
}

func TestControllerSignalChangedWakesWaitersAfterNonTerminalTransition(t *testing.T) {
	t.Parallel()

	controller := NewController()
	results := make(chan Snapshot, 1)
	go func() {
		snapshot, _ := controller.WaitState(context.Background(), StateStarting)
		results <- snapshot
	}()

	_, _ = controller.BeginStart()
	if got := mustReceiveSnapshot(t, results); got.State != StateStarting {
		t.Fatalf("waited snapshot state = %s, want starting", got.State)
	}
}

func TestControllerSignalChangedRotatesForNonTerminalTransition(t *testing.T) {
	t.Parallel()

	// Each non-terminal commit closes the current changed channel and installs a
	// fresh one so later waiters observe later revisions independently.
	controller := NewController()
	_, firstChanged, firstDone := controller.waitSnapshot()
	_, _ = controller.BeginStart()
	_, secondChanged, secondDone := controller.waitSnapshot()
	if firstChanged == secondChanged {
		t.Fatal("changed channel was not rotated after non-terminal transition")
	}
	if firstDone != secondDone {
		t.Fatal("done channel rotated before terminal transition")
	}
	mustSignalClosed(t, firstChanged)
	mustNotSignalClosed(t, secondChanged)
}

func TestControllerSignalTerminalClosesChangedAndDoneWithoutRotation(t *testing.T) {
	t.Parallel()

	controller := NewController()
	_, changed, done := controller.waitSnapshot()
	_, _ = controller.BeginStop()
	_, afterChanged, afterDone := controller.waitSnapshot()
	if changed != afterChanged {
		t.Fatal("changed channel rotated after terminal transition")
	}
	if done != afterDone {
		t.Fatal("done channel rotated after terminal transition")
	}
	mustSignalClosed(t, changed)
	mustSignalClosed(t, done)
}

func TestControllerSignalCommitTimeFallsBackFromZeroClock(t *testing.T) {
	t.Parallel()

	controller := NewController(WithClock(testClock{}))
	transition, err := controller.BeginStart()
	if err != nil {
		t.Fatalf("BeginStart = %v", err)
	}
	if transition.At.IsZero() {
		t.Fatal("Transition.At is zero, want fallback time.Now value")
	}
}

func TestControllerSignalCommitTimeHandlesNilClock(t *testing.T) {
	t.Parallel()

	// commitTimeLocked has its own nil-clock fallback because tests and
	// zero-value paths may reach it without NewController normalization.
	var controller Controller
	controller.mu.Lock()
	at := controller.commitTimeLocked()
	controller.mu.Unlock()
	if at.IsZero() {
		t.Fatal("commitTimeLocked returned zero time")
	}
}

func TestControllerSignalDoneSignalZeroValueStable(t *testing.T) {
	t.Parallel()

	var controller Controller
	first := controller.doneSignal()
	second := controller.doneSignal()
	if first == nil {
		t.Fatal("doneSignal returned nil")
	}
	if first != second {
		t.Fatal("doneSignal returned different channels")
	}
}

func TestControllerSignalWaitSnapshotReturnsSignals(t *testing.T) {
	t.Parallel()

	controller := NewController()
	snapshot, changed, done := controller.waitSnapshot()
	if !snapshot.IsValid() {
		t.Fatalf("waitSnapshot snapshot is invalid: %+v", snapshot)
	}
	if changed == nil || done == nil {
		t.Fatalf("waitSnapshot returned changed=%v done=%v, want non-nil", changed, done)
	}
}

func TestControllerSignalCommitTimeUsesClockValue(t *testing.T) {
	t.Parallel()

	clock := &countingClock{now: testTime}
	controller := NewController(WithClock(clock))
	transition, err := controller.BeginStart()
	if err != nil {
		t.Fatalf("BeginStart = %v", err)
	}
	if !transition.At.Equal(testTime) {
		t.Fatalf("Transition.At = %v, want %v", transition.At, testTime)
	}
	if clock.calls != 1 {
		t.Fatalf("clock calls = %d, want 1", clock.calls)
	}
}

func TestControllerSignalClosedTerminalZeroValueWaitSnapshot(t *testing.T) {
	t.Parallel()

	var controller Controller
	_, _ = controller.BeginStop()
	_, changed, done := controller.waitSnapshot()
	mustSignalClosed(t, changed)
	mustSignalClosed(t, done)
}

func TestControllerSignalEnsureInitializedClosesSignalsForTerminalState(t *testing.T) {
	t.Parallel()

	// Lazy initialization must respect an already-terminal state so zero-value
	// derived controllers cannot expose open signals after terminal state is set.
	controller := &Controller{state: StateStopped}
	_, changed, done := controller.waitSnapshot()
	mustSignalClosed(t, changed)
	mustSignalClosed(t, done)
}

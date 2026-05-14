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
	"errors"
	"sync"
	"testing"
)

func TestControllerApplyAllowedTransitions(t *testing.T) {
	t.Parallel()

	cause := errors.New("failed")
	tests := []struct {
		name string
		run  func(*Controller) (Transition, error)
		want Transition
	}{
		{"BeginStart", func(c *Controller) (Transition, error) { return c.BeginStart() }, Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}},
		{"MarkRunning", func(c *Controller) (Transition, error) {
			_, _ = c.BeginStart()
			return c.MarkRunning()
		}, Transition{From: StateStarting, To: StateRunning, Event: EventMarkRunning}},
		{"BeginStop from new", func(c *Controller) (Transition, error) { return c.BeginStop() }, Transition{From: StateNew, To: StateStopped, Event: EventBeginStop}},
		{"BeginStop from starting", func(c *Controller) (Transition, error) {
			_, _ = c.BeginStart()
			return c.BeginStop()
		}, Transition{From: StateStarting, To: StateStopping, Event: EventBeginStop}},
		{"BeginStop from running", func(c *Controller) (Transition, error) {
			_, _ = c.BeginStart()
			_, _ = c.MarkRunning()
			return c.BeginStop()
		}, Transition{From: StateRunning, To: StateStopping, Event: EventBeginStop}},
		{"MarkStopped", func(c *Controller) (Transition, error) {
			_, _ = c.BeginStart()
			_, _ = c.BeginStop()
			return c.MarkStopped()
		}, Transition{From: StateStopping, To: StateStopped, Event: EventMarkStopped}},
		{"MarkFailed from starting", func(c *Controller) (Transition, error) {
			_, _ = c.BeginStart()
			return c.MarkFailed(cause)
		}, Transition{From: StateStarting, To: StateFailed, Event: EventMarkFailed, Cause: cause}},
		{"MarkFailed from running", func(c *Controller) (Transition, error) {
			_, _ = c.BeginStart()
			_, _ = c.MarkRunning()
			return c.MarkFailed(cause)
		}, Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed, Cause: cause}},
		{"MarkFailed from stopping", func(c *Controller) (Transition, error) {
			_, _ = c.BeginStart()
			_, _ = c.BeginStop()
			return c.MarkFailed(cause)
		}, Transition{From: StateStopping, To: StateFailed, Event: EventMarkFailed, Cause: cause}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := NewController(WithClock(testClock{now: testTime}))
			transition, err := tc.run(controller)
			if err != nil {
				t.Fatalf("%s = %v, want nil", tc.name, err)
			}
			tc.want.Revision = transition.Revision
			tc.want.At = testTime
			assertTransitionEqual(t, transition, tc.want)
			if !controller.Snapshot().IsValid() {
				t.Fatalf("snapshot after %s is invalid: %+v", tc.name, controller.Snapshot())
			}
		})
	}
}

func TestControllerApplyDoesNotMutateOnInvalidTransition(t *testing.T) {
	t.Parallel()

	controller := NewController()
	snapshot, changed, done := controller.waitSnapshot()
	_, err := controller.MarkRunning()
	if err == nil {
		t.Fatal("MarkRunning err = nil, want invalid transition")
	}
	mustMatch(t, err, ErrInvalidTransition)
	assertSnapshotEqual(t, controller.Snapshot(), snapshot)
	mustNotSignalClosed(t, changed)
	mustNotSignalClosed(t, done)
}

func TestControllerApplyTerminalTransitionError(t *testing.T) {
	t.Parallel()

	controller := NewController()
	_, _ = controller.BeginStop()
	snapshot, changed, done := controller.waitSnapshot()

	_, err := controller.BeginStart()
	if err == nil {
		t.Fatal("BeginStart from terminal err = nil, want terminal")
	}
	mustMatch(t, err, ErrTerminalState)
	mustMatch(t, err, ErrInvalidTransition)
	assertSnapshotEqual(t, controller.Snapshot(), snapshot)
	mustSignalClosed(t, changed)
	mustSignalClosed(t, done)
}

func TestControllerApplyMarkFailedRequiresCause(t *testing.T) {
	t.Parallel()

	controller := NewController()
	_, _ = controller.BeginStart()
	before, changed, done := controller.waitSnapshot()

	_, err := controller.MarkFailed(nil)
	if err == nil {
		t.Fatal("MarkFailed(nil) err = nil, want failure cause required")
	}
	mustMatch(t, err, ErrFailureCauseRequired)
	assertSnapshotEqual(t, controller.Snapshot(), before)
	mustNotSignalClosed(t, changed)
	mustNotSignalClosed(t, done)
}

func TestControllerApplyGuardRejectionDoesNotCommit(t *testing.T) {
	t.Parallel()

	// Failed apply is atomic: guards reject before state, revision, signals, and
	// observers are mutated.
	rejection := errors.New("blocked")
	observerCalled := false
	controller := NewController(
		WithGuard(TransitionGuardFunc(func(Transition) error { return rejection })),
		WithObserver(ObserverFunc(func(Transition) { observerCalled = true })),
	)
	before, changed, done := controller.waitSnapshot()

	_, err := controller.BeginStart()
	if err == nil {
		t.Fatal("BeginStart err = nil, want guard error")
	}
	mustMatch(t, err, ErrGuardRejected)
	mustMatch(t, err, rejection)
	assertSnapshotEqual(t, controller.Snapshot(), before)
	mustNotSignalClosed(t, changed)
	mustNotSignalClosed(t, done)
	if observerCalled {
		t.Fatal("observer called after failed apply")
	}
}

func TestControllerApplyGuardsObserveCandidateBeforeCommit(t *testing.T) {
	t.Parallel()

	seen := make(chan Transition, 1)
	controller := NewController(WithGuard(TransitionGuardFunc(func(transition Transition) error {
		seen <- transition
		return nil
	})))

	if _, err := controller.BeginStart(); err != nil {
		t.Fatalf("BeginStart = %v", err)
	}
	transition := mustReceiveTransition(t, seen)
	if transition.Revision != 0 || !transition.At.IsZero() {
		t.Fatalf("guard saw committed metadata: %+v", transition)
	}
}

func TestControllerApplyObserversObserveCommittedTransition(t *testing.T) {
	t.Parallel()

	seen := make(chan Transition, 1)
	controller := NewController(WithObserver(ObserverFunc(func(transition Transition) {
		seen <- transition
	})))

	transition, err := controller.BeginStart()
	if err != nil {
		t.Fatalf("BeginStart = %v", err)
	}
	observed := mustReceiveTransition(t, seen)
	assertTransitionEqual(t, observed, transition)
	if !observed.IsCommitted() {
		t.Fatalf("observer saw uncommitted transition: %+v", observed)
	}
}

func TestControllerApplyRevisionIncrementsMonotonically(t *testing.T) {
	t.Parallel()

	controller := NewController()
	first, _ := controller.BeginStart()
	second, _ := controller.MarkRunning()
	third, _ := controller.BeginStop()
	if got, want := []uint64{first.Revision, second.Revision, third.Revision}, []uint64{1, 2, 3}; !equalUint64s(got, want) {
		t.Fatalf("revisions = %v, want %v", got, want)
	}
}

func TestControllerApplyStoresFailureCauseOnlyForFailedState(t *testing.T) {
	t.Parallel()

	cause := errors.New("failed")
	controller := NewController()
	_, _ = controller.BeginStart()
	if controller.Snapshot().FailureCause != nil {
		t.Fatal("non-failed snapshot carries failure cause")
	}
	_, _ = controller.MarkFailed(cause)
	if got := controller.Snapshot().FailureCause; got != cause {
		t.Fatalf("FailureCause = %v, want %v", got, cause)
	}
}

func TestControllerApplyConcurrentBeginStartOnlyOneSucceeds(t *testing.T) {
	t.Parallel()

	controller := NewController()
	var wg sync.WaitGroup
	errs := make(chan error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := controller.BeginStart()
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)

	successes := 0
	failures := 0
	for err := range errs {
		if err == nil {
			successes++
			continue
		}
		mustMatch(t, err, ErrInvalidTransition)
		failures++
	}
	if successes != 1 || failures != 1 {
		t.Fatalf("successes=%d failures=%d, want 1 and 1", successes, failures)
	}
}

func TestControllerApplyConcurrentReadsDuringTransitions(t *testing.T) {
	t.Parallel()

	// State and Snapshot must remain race-safe while transitions commit; the
	// reader exits through a channel rather than timing assumptions.
	controller := NewController()
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				_ = controller.State()
				_ = controller.Snapshot()
			}
		}
	}()

	_, _ = controller.BeginStart()
	_, _ = controller.MarkRunning()
	_, _ = controller.BeginStop()
	_, _ = controller.MarkStopped()
	close(done)
	wg.Wait()
}

func equalUint64s(got, want []uint64) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

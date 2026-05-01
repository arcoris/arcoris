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

func TestControllerZeroValueUsable(t *testing.T) {
	t.Parallel()

	var controller Controller
	transition, err := controller.BeginStart()
	if err != nil {
		t.Fatalf("BeginStart = %v, want nil", err)
	}
	if transition.To != StateStarting {
		t.Fatalf("transition.To = %s, want starting", transition.To)
	}
	if controller.State() != StateStarting {
		t.Fatalf("State = %s, want starting", controller.State())
	}
}

func TestNewControllerInitialSnapshotValid(t *testing.T) {
	t.Parallel()

	snapshot := NewController().Snapshot()
	if !snapshot.IsValid() {
		t.Fatalf("initial snapshot %+v is invalid", snapshot)
	}
}

func TestControllerTransitions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		run   func(*Controller) (Transition, error)
		state State
		done  bool
	}{
		{
			name:  "BeginStart",
			run:   func(c *Controller) (Transition, error) { return c.BeginStart() },
			state: StateStarting,
		},
		{
			name: "MarkRunning",
			run: func(c *Controller) (Transition, error) {
				if _, err := c.BeginStart(); err != nil {
					return Transition{}, err
				}
				return c.MarkRunning()
			},
			state: StateRunning,
		},
		{
			name:  "BeginStop from new",
			run:   func(c *Controller) (Transition, error) { return c.BeginStop() },
			state: StateStopped,
			done:  true,
		},
		{
			name: "BeginStop from running",
			run: func(c *Controller) (Transition, error) {
				if _, err := c.BeginStart(); err != nil {
					return Transition{}, err
				}
				if _, err := c.MarkRunning(); err != nil {
					return Transition{}, err
				}
				return c.BeginStop()
			},
			state: StateStopping,
		},
		{
			name: "MarkStopped",
			run: func(c *Controller) (Transition, error) {
				if _, err := c.BeginStart(); err != nil {
					return Transition{}, err
				}
				if _, err := c.BeginStop(); err != nil {
					return Transition{}, err
				}
				return c.MarkStopped()
			},
			state: StateStopped,
			done:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			controller := NewController(WithClock(testClock{now: testTime}))
			transition, err := tt.run(controller)
			if err != nil {
				t.Fatalf("%s returned %v, want nil", tt.name, err)
			}
			if transition.To != tt.state {
				t.Fatalf("transition.To = %s, want %s", transition.To, tt.state)
			}
			if transition.Revision == 0 || transition.At.IsZero() {
				t.Fatalf("transition metadata = revision %d at %v, want non-zero", transition.Revision, transition.At)
			}
			if got := controller.Snapshot(); got.State != tt.state || !got.IsValid() {
				t.Fatalf("snapshot = %+v, want valid state %s", got, tt.state)
			}
			if tt.done {
				mustSignalClosed(t, controller.Done())
			}
		})
	}
}

func TestControllerMarkFailed(t *testing.T) {
	t.Parallel()

	cause := errors.New("boom")
	controller := NewController()
	if _, err := controller.BeginStart(); err != nil {
		t.Fatalf("BeginStart = %v", err)
	}

	transition, err := controller.MarkFailed(cause)
	if err != nil {
		t.Fatalf("MarkFailed = %v, want nil", err)
	}
	assertTransitionEqual(t, transition, Transition{
		From:     StateStarting,
		To:       StateFailed,
		Event:    EventMarkFailed,
		Revision: transition.Revision,
		At:       transition.At,
		Cause:    cause,
	})
	snapshot := controller.Snapshot()
	assertSnapshotEqual(t, snapshot, Snapshot{
		State:          StateFailed,
		Revision:       transition.Revision,
		LastTransition: transition,
		FailureCause:   cause,
	})
	mustSignalClosed(t, controller.Done())
}

func TestControllerFailedOperationsDoNotMutate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ctrl func() *Controller
		call func(*Controller) (Transition, error)
		want error
	}{
		{
			name: "missing cause",
			ctrl: func() *Controller {
				c := NewController()
				_, _ = c.BeginStart()
				return c
			},
			call: func(c *Controller) (Transition, error) { return c.MarkFailed(nil) },
			want: ErrFailureCauseRequired,
		},
		{
			name: "invalid transition",
			ctrl: func() *Controller { return NewController() },
			call: func(c *Controller) (Transition, error) { return c.MarkRunning() },
			want: ErrInvalidTransition,
		},
		{
			name: "terminal transition",
			ctrl: func() *Controller {
				c := NewController()
				_, _ = c.BeginStop()
				return c
			},
			call: func(c *Controller) (Transition, error) { return c.BeginStart() },
			want: ErrTerminalState,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			controller := tt.ctrl()
			before := controller.Snapshot()
			_, err := tt.call(controller)
			if err == nil {
				t.Fatal("operation err = nil, want error")
			}
			mustMatch(t, err, tt.want)
			after := controller.Snapshot()
			assertSnapshotEqual(t, after, before)
		})
	}
}

func TestControllerGuardRejectionDoesNotCommit(t *testing.T) {
	t.Parallel()

	rejection := errors.New("blocked")
	observerCalled := false
	controller := NewController(
		WithGuard(TransitionGuardFunc(func(Transition) error { return rejection })),
		WithObserver(ObserverFunc(func(Transition) { observerCalled = true })),
	)
	before := controller.Snapshot()
	beforeChanged, changed, done := controller.waitSnapshot()
	assertSnapshotEqual(t, beforeChanged, before)

	_, err := controller.BeginStart()
	if err == nil {
		t.Fatal("BeginStart err = nil, want guard error")
	}
	mustMatch(t, err, ErrGuardRejected)
	mustMatch(t, err, rejection)
	after := controller.Snapshot()
	assertSnapshotEqual(t, after, before)
	mustNotSignalClosed(t, changed)
	mustNotSignalClosed(t, done)
	if observerCalled {
		t.Fatal("observer called after guard rejection")
	}
}

func TestControllerObserverAfterCommit(t *testing.T) {
	t.Parallel()

	observed := make(chan Transition, 1)
	controller := NewController(WithObserver(ObserverFunc(func(transition Transition) {
		observed <- transition
	})))

	transition, err := controller.BeginStart()
	if err != nil {
		t.Fatalf("BeginStart = %v", err)
	}

	got := <-observed
	assertTransitionEqual(t, got, transition)
	if !got.IsCommitted() {
		t.Fatalf("observer transition = %+v, want committed", got)
	}
}

func TestControllerConcurrentBeginStartOnlyOneSucceeds(t *testing.T) {
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

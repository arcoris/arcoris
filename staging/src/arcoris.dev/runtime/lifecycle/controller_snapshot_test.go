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

func TestControllerSnapshotStateReturnsCurrentState(t *testing.T) {
	t.Parallel()

	controller := NewController()
	if got := controller.State(); got != StateNew {
		t.Fatalf("State before transition = %s, want new", got)
	}
	_, _ = controller.BeginStart()
	if got := controller.State(); got != StateStarting {
		t.Fatalf("State after transition = %s, want starting", got)
	}
}

func TestControllerSnapshotConsistentPointInTimeView(t *testing.T) {
	t.Parallel()

	// Snapshot is copied while holding the controller lock, so State, Revision,
	// LastTransition, and FailureCause describe the same committed point in time.
	controller := NewController(WithClock(testClock{now: testTime}))
	transition, _ := controller.BeginStart()
	snapshot := controller.Snapshot()
	assertSnapshotEqual(t, snapshot, Snapshot{
		State:          StateStarting,
		Revision:       transition.Revision,
		LastTransition: transition,
	})
}

func TestControllerSnapshotValidAcrossTransitions(t *testing.T) {
	t.Parallel()

	controller := NewController()
	if !controller.Snapshot().IsValid() {
		t.Fatal("snapshot before first transition is invalid")
	}
	steps := []func() (Transition, error){
		controller.BeginStart,
		controller.MarkRunning,
		controller.BeginStop,
		controller.MarkStopped,
	}
	for _, step := range steps {
		if _, err := step(); err != nil {
			t.Fatalf("transition = %v", err)
		}
		if snapshot := controller.Snapshot(); !snapshot.IsValid() {
			t.Fatalf("snapshot after transition is invalid: %+v", snapshot)
		}
	}
}

func TestControllerSnapshotAfterFailedLifecycleIncludesCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("failed")
	controller := NewController()
	_, _ = controller.BeginStart()
	transition, _ := controller.MarkFailed(cause)
	snapshot := controller.Snapshot()
	assertSnapshotEqual(t, snapshot, Snapshot{
		State:          StateFailed,
		Revision:       transition.Revision,
		LastTransition: transition,
		FailureCause:   cause,
	})
}

func TestControllerSnapshotIsCopyable(t *testing.T) {
	t.Parallel()

	// Mutating a returned Snapshot must not mutate controller internals; callers
	// own their read-model copy.
	controller := NewController()
	_, _ = controller.BeginStart()
	snapshot := controller.Snapshot()
	snapshot.State = StateFailed
	snapshot.Revision = 99
	snapshot.LastTransition = Transition{}

	got := controller.Snapshot()
	if got.State != StateStarting || got.Revision != 1 || got.LastTransition.IsZero() {
		t.Fatalf("controller snapshot was mutated through returned copy: %+v", got)
	}
}

func TestControllerSnapshotConcurrentWithTransitions(t *testing.T) {
	t.Parallel()

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

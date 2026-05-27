// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lifecycle

import (
	"testing"
	"time"

	channelassert "arcoris.dev/testutil/channel"
)

func TestLifecycleObserverCanCallSnapshotWithoutDeadlock(t *testing.T) {
	t.Parallel()

	observed := make(chan Snapshot, 1)
	var controller *Controller
	controller = NewController(WithObserver(ObserverFunc(func(Transition) {
		observed <- controller.Snapshot()
	})))

	transition, err := controller.BeginStart()
	if err != nil {
		t.Fatalf("BeginStart() error = %v", err)
	}

	snap := channelassert.RequireReceive(t, observed, time.Second)
	if snap.State != transition.To {
		t.Fatalf("observer Snapshot().State = %s, want %s", snap.State, transition.To)
	}
}

func TestLifecycleObserverCannotRollbackCommittedTransition(t *testing.T) {
	t.Parallel()

	controller := NewController(WithObserver(ObserverFunc(func(Transition) {
		panic("observer failed")
	})))

	func() {
		defer func() {
			if recovered := recover(); recovered == nil {
				t.Fatal("observer panic was recovered, want panic to propagate")
			}
		}()
		_, _ = controller.BeginStart()
	}()

	snap := controller.Snapshot()
	if snap.State != StateStarting {
		t.Fatalf("Snapshot().State = %s, want %s", snap.State, StateStarting)
	}
	if snap.Revision != 1 {
		t.Fatalf("Snapshot().Revision = %d, want 1", snap.Revision)
	}
}

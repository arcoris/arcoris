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
	"testing"
)

func TestNewControllerInitialState(t *testing.T) {
	t.Parallel()

	controller := NewController()
	if got := controller.State(); got != StateNew {
		t.Fatalf("State = %s, want new", got)
	}
}

func TestNewControllerInitialRevision(t *testing.T) {
	t.Parallel()

	if got := NewController().Snapshot().Revision; got != 0 {
		t.Fatalf("initial revision = %d, want 0", got)
	}
}

func TestNewControllerInitialSnapshotValid(t *testing.T) {
	t.Parallel()

	snapshot := NewController().Snapshot()
	if !snapshot.IsValid() {
		t.Fatalf("initial snapshot %+v is invalid", snapshot)
	}
}

func TestNewControllerCopiesGuardAndObserverSlices(t *testing.T) {
	t.Parallel()

	// Controller copies config slices at construction so later caller-side slice
	// mutation cannot change guard or observer behavior behind the controller.
	rejection := errors.New("blocked")
	allow := TransitionGuardFunc(func(Transition) error { return nil })
	block := TransitionGuardFunc(func(Transition) error { return rejection })
	guards := []TransitionGuard{allow}
	observed := make(chan Transition, 1)
	observer := ObserverFunc(func(transition Transition) { observed <- transition })
	observers := []Observer{observer}

	controller := NewController(WithGuards(guards...), WithObservers(observers...))
	guards[0] = block
	observers[0] = ObserverFunc(func(Transition) {
		t.Fatal("mutated observer should not be retained")
	})

	transition, err := controller.BeginStart()
	if err != nil {
		t.Fatalf("BeginStart = %v, want nil", err)
	}
	assertTransitionEqual(t, mustReceiveTransition(t, observed), transition)
}

func TestNewControllerFallsBackWhenOptionClearsClock(t *testing.T) {
	t.Parallel()

	// NewController normalizes config after options run, so even a package-local
	// option that clears the clock cannot produce zero commit timestamps.
	controller := NewController(Option(func(config *controllerConfig) {
		config.now = nil
	}))
	transition, err := controller.BeginStart()
	if err != nil {
		t.Fatalf("BeginStart = %v, want nil", err)
	}
	if transition.At.IsZero() {
		t.Fatal("Transition.At is zero, want fallback time source")
	}
}

func TestControllerZeroValueUsableForLazyConstruction(t *testing.T) {
	t.Parallel()

	// The zero-value Controller contract lets embedded controllers work before a
	// constructor is called; detailed lazy signal behavior belongs to signal tests.
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

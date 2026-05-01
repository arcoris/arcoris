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

func TestGuardChainNilAndEmptyAllow(t *testing.T) {
	t.Parallel()

	// Missing guard chains are valid optional wiring and must behave the same as
	// an empty list of preconditions.
	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	if err := (GuardChain)(nil).Allow(transition); err != nil {
		t.Fatalf("nil GuardChain.Allow = %v, want nil", err)
	}
	if err := (GuardChain{}).Allow(transition); err != nil {
		t.Fatalf("empty GuardChain.Allow = %v, want nil", err)
	}
}

func TestGuardChainIgnoresNilEntries(t *testing.T) {
	t.Parallel()

	called := false
	chain := GuardChain{
		nil,
		TransitionGuardFunc(func(Transition) error {
			called = true
			return nil
		}),
	}

	if err := chain.Allow(Transition{}); err != nil {
		t.Fatalf("GuardChain.Allow = %v, want nil", err)
	}
	if !called {
		t.Fatal("non-nil guard was not called")
	}
}

func TestGuardChainOrderAndShortCircuit(t *testing.T) {
	t.Parallel()

	// Guard evaluation is ordered and short-circuited so the first rejected
	// precondition is the domain cause returned to Controller.
	rejection := errors.New("blocked")
	var order []string
	chain := GuardChain{
		TransitionGuardFunc(func(Transition) error {
			order = append(order, "first")
			return nil
		}),
		TransitionGuardFunc(func(Transition) error {
			order = append(order, "second")
			return rejection
		}),
		TransitionGuardFunc(func(Transition) error {
			order = append(order, "third")
			return nil
		}),
	}

	if err := chain.Allow(Transition{}); err != rejection {
		t.Fatalf("GuardChain.Allow err = %v, want %v", err, rejection)
	}
	assertDeepEqual(t, order, []string{"first", "second"})
}

func TestAllowTransitionDelegatesToGuardChain(t *testing.T) {
	t.Parallel()

	rejection := errors.New("blocked")
	err := allowTransition([]TransitionGuard{
		TransitionGuardFunc(func(Transition) error { return rejection }),
	}, Transition{})
	if err != rejection {
		t.Fatalf("allowTransition err = %v, want %v", err, rejection)
	}
}

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

func TestTransitionGuardFunc(t *testing.T) {
	t.Parallel()

	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}

	if err := (TransitionGuardFunc)(nil).Allow(transition); err != nil {
		t.Fatalf("nil TransitionGuardFunc.Allow = %v, want nil", err)
	}

	var got Transition
	guard := TransitionGuardFunc(func(transition Transition) error {
		got = transition
		return nil
	})

	if err := guard.Allow(transition); err != nil {
		t.Fatalf("TransitionGuardFunc.Allow = %v, want nil", err)
	}
	assertTransitionEqual(t, got, transition)
}

func TestGuardChain(t *testing.T) {
	t.Parallel()

	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	rejection := errors.New("blocked")

	if err := (GuardChain)(nil).Allow(transition); err != nil {
		t.Fatalf("nil GuardChain.Allow = %v, want nil", err)
	}
	if err := (GuardChain{}).Allow(transition); err != nil {
		t.Fatalf("empty GuardChain.Allow = %v, want nil", err)
	}

	var order []string
	chain := GuardChain{
		nil,
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

	err := chain.Allow(transition)
	if err != rejection {
		t.Fatalf("GuardChain.Allow err = %v, want %v", err, rejection)
	}
	if got, want := order, []string{"first", "second"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("guard order = %v, want %v", got, want)
	}
}

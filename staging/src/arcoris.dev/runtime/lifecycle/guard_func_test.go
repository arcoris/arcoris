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

func TestTransitionGuardFuncNilAllowsTransition(t *testing.T) {
	t.Parallel()

	// Nil function adapters make optional guard wiring safe: an absent callback
	// must not reject a transition or panic during controller construction.
	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	if err := (TransitionGuardFunc)(nil).Allow(transition); err != nil {
		t.Fatalf("nil TransitionGuardFunc.Allow = %v, want nil", err)
	}
}

func TestTransitionGuardFuncReceivesExactTransition(t *testing.T) {
	t.Parallel()

	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
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

func TestTransitionGuardFuncPreservesReturnedError(t *testing.T) {
	t.Parallel()

	rejection := errors.New("blocked")
	guard := TransitionGuardFunc(func(Transition) error {
		return rejection
	})

	if err := guard.Allow(Transition{}); err != rejection {
		t.Fatalf("TransitionGuardFunc.Allow err = %v, want %v", err, rejection)
	}
}

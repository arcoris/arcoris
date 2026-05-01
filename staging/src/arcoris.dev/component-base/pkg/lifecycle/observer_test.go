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

import "testing"

func TestObserverFunc(t *testing.T) {
	t.Parallel()

	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}

	(ObserverFunc)(nil).ObserveLifecycleTransition(transition)

	var got Transition
	ObserverFunc(func(transition Transition) {
		got = transition
	}).ObserveLifecycleTransition(transition)

	if got != transition {
		t.Fatalf("observer received %v, want %v", got, transition)
	}
}

func TestObserverChain(t *testing.T) {
	t.Parallel()

	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	var order []string
	var seen []Transition

	chain := ObserverChain{
		nil,
		ObserverFunc(func(transition Transition) {
			order = append(order, "first")
			seen = append(seen, transition)
		}),
		ObserverFunc(func(transition Transition) {
			order = append(order, "second")
			seen = append(seen, transition)
		}),
	}

	chain.ObserveLifecycleTransition(transition)

	if got, want := order, []string{"first", "second"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("observer order = %v, want %v", got, want)
	}
	for _, got := range seen {
		if got != transition {
			t.Fatalf("observer received %v, want %v", got, transition)
		}
	}
}

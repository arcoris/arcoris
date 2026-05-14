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

func TestObserverChainNilAndEmptyAreNoop(t *testing.T) {
	t.Parallel()

	// Observers are post-commit notifications, not validation hooks; an absent
	// chain must be a no-op and must never reject or mutate a transition.
	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	(ObserverChain)(nil).ObserveLifecycleTransition(transition)
	(ObserverChain{}).ObserveLifecycleTransition(transition)
}

func TestObserverChainIgnoresNilEntriesAndPreservesOrder(t *testing.T) {
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

	assertDeepEqual(t, order, []string{"first", "second"})
	if len(seen) != 2 {
		t.Fatalf("seen len = %d, want 2", len(seen))
	}
	for _, got := range seen {
		assertTransitionEqual(t, got, transition)
	}
}

func TestNotifyObserversDelegatesToObserverChain(t *testing.T) {
	t.Parallel()

	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	observed := make(chan Transition, 1)
	notifyObservers([]Observer{
		ObserverFunc(func(transition Transition) { observed <- transition }),
	}, transition)

	assertTransitionEqual(t, mustReceiveTransition(t, observed), transition)
}

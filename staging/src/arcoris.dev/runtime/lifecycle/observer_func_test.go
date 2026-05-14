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

func TestObserverFuncNilIgnoresTransition(t *testing.T) {
	t.Parallel()

	// Nil observer adapters make optional diagnostics wiring safe: absence of an
	// observer must not affect committed lifecycle transitions.
	(ObserverFunc)(nil).ObserveLifecycleTransition(Transition{
		From:  StateNew,
		To:    StateStarting,
		Event: EventBeginStart,
	})
}

func TestObserverFuncReceivesExactTransition(t *testing.T) {
	t.Parallel()

	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	var got Transition
	ObserverFunc(func(transition Transition) {
		got = transition
	}).ObserveLifecycleTransition(transition)

	assertTransitionEqual(t, got, transition)
}

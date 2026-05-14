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

func TestReduceTransitionAllowedPair(t *testing.T) {
	t.Parallel()

	// reduceTransition is pure: it derives an uncommitted candidate and does not
	// run guards, assign Revision/At, or signal observers.
	cause := errors.New("note")
	transition, ok := reduceTransition(StateNew, EventBeginStart, cause)
	if !ok {
		t.Fatal("reduceTransition ok = false, want true")
	}
	assertTransitionEqual(t, transition, Transition{
		From:  StateNew,
		To:    StateStarting,
		Event: EventBeginStart,
		Cause: cause,
	})
	if transition.Revision != 0 || !transition.At.IsZero() {
		t.Fatalf("candidate metadata = revision %d at %v, want zero", transition.Revision, transition.At)
	}
}

func TestReduceTransitionInvalidFallback(t *testing.T) {
	t.Parallel()

	cause := errors.New("preserved")
	transition, ok := reduceTransition(StateRunning, EventBeginStart, cause)
	if ok {
		t.Fatal("reduceTransition ok = true, want false")
	}
	assertTransitionEqual(t, transition, Transition{
		From:  StateRunning,
		To:    StateRunning,
		Event: EventBeginStart,
		Cause: cause,
	})
}

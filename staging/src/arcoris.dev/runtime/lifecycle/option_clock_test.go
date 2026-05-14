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
	"testing"
	"time"
)

func TestWithClockNilLeavesClockUnchanged(t *testing.T) {
	t.Parallel()

	config := controllerConfig{now: func() time.Time { return testTime }}
	WithClock(nil)(&config)
	if got := config.now(); !got.Equal(testTime) {
		t.Fatalf("config.now() = %v, want %v", got, testTime)
	}
}

func TestWithClockCustomSourceUsedForCommittedTransition(t *testing.T) {
	t.Parallel()

	// lifecycle only depends on the minimal Now method; richer clock behavior is
	// outside the controller contract.
	controller := NewController(WithClock(testClock{now: testTime}))
	transition, err := controller.BeginStart()
	if err != nil {
		t.Fatalf("BeginStart = %v", err)
	}
	if !transition.At.Equal(testTime) {
		t.Fatalf("Transition.At = %v, want %v", transition.At, testTime)
	}
}

func TestWithClockDoesNotCallNowDuringConfigConstruction(t *testing.T) {
	t.Parallel()

	clock := &countingClock{now: testTime}
	_ = newControllerConfig(WithClock(clock))
	if clock.calls != 0 {
		t.Fatalf("clock calls during config construction = %d, want 0", clock.calls)
	}
}

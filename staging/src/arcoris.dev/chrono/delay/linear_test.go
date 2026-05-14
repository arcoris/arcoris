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

package delay

import (
	"testing"
	"time"
)

func TestLinearRejectsNegativeInput(t *testing.T) {
	mustPanicWith(t, errNegativeLinearInitialDelay, func() {
		Linear(-time.Nanosecond, time.Second)
	})
	mustPanicWith(t, errNegativeLinearStep, func() {
		Linear(0, -time.Nanosecond)
	})
}

func TestLinearAllowsZeroInitialDelay(t *testing.T) {
	sequence := Linear(0, time.Second).NewSequence()

	mustNext(t, sequence, 0)
	mustNext(t, sequence, time.Second)
}

func TestLinearSequenceGrowsByStep(t *testing.T) {
	sequence := Linear(time.Second, 500*time.Millisecond).NewSequence()

	mustNext(t, sequence, time.Second)
	mustNext(t, sequence, 1500*time.Millisecond)
	mustNext(t, sequence, 2*time.Second)
}

func TestLinearZeroStepBehavesLikeFixed(t *testing.T) {
	sequence := Linear(time.Second, 0).NewSequence()

	mustNext(t, sequence, time.Second)
	mustNext(t, sequence, time.Second)
}

func TestLinearSequencesHaveIndependentIndexes(t *testing.T) {
	schedule := Linear(time.Second, time.Second)

	l := schedule.NewSequence()
	r := schedule.NewSequence()

	mustNext(t, l, time.Second)
	mustNext(t, l, 2*time.Second)
	mustNext(t, r, time.Second)
}

func TestLinearDelaySaturates(t *testing.T) {
	if got := linearDelay(maxDuration, time.Nanosecond, 1); got != maxDuration {
		t.Fatalf("linearDelay() = %s, want %s", got, maxDuration)
	}
	if got := linearDelay(time.Second, maxDuration, 2); got != maxDuration {
		t.Fatalf("linearDelay() = %s, want %s", got, maxDuration)
	}
}

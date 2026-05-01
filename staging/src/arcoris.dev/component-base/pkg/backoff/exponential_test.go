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

package backoff

import (
	"math"
	"testing"
	"time"
)

func TestExponentialRejectsInvalidInput(t *testing.T) {
	mustPanicWith(t, errNonPositiveExponentialInitialDelay, func() {
		Exponential(0, 2)
	})
	mustPanicWith(t, errNonPositiveExponentialInitialDelay, func() {
		Exponential(-time.Nanosecond, 2)
	})
	mustPanicWith(t, errInvalidExponentialMultiplier, func() {
		Exponential(time.Second, 1)
	})
	mustPanicWith(t, errInvalidExponentialMultiplier, func() {
		Exponential(time.Second, math.NaN())
	})
	mustPanicWith(t, errInvalidExponentialMultiplier, func() {
		Exponential(time.Second, math.Inf(1))
	})
}

func TestExponentialSequenceGrowsByMultiplier(t *testing.T) {
	sequence := Exponential(time.Second, 2).NewSequence()

	mustNext(t, sequence, time.Second)
	mustNext(t, sequence, 2*time.Second)
	mustNext(t, sequence, 4*time.Second)
}

func TestExponentialKeepsFloatStateBetweenIntegerDurations(t *testing.T) {
	sequence := Exponential(2*time.Nanosecond, 1.5).NewSequence()

	mustNext(t, sequence, 2*time.Nanosecond)
	mustNext(t, sequence, 3*time.Nanosecond)
	mustNext(t, sequence, 4*time.Nanosecond)
	mustNext(t, sequence, 6*time.Nanosecond)
}

func TestExponentialSequencesHaveIndependentState(t *testing.T) {
	schedule := Exponential(time.Second, 2)

	left := schedule.NewSequence()
	right := schedule.NewSequence()

	mustNext(t, left, time.Second)
	mustNext(t, left, 2*time.Second)
	mustNext(t, right, time.Second)
}

func TestExponentialSequenceSaturates(t *testing.T) {
	sequence := (&exponentialSequence{next: maxDurationFloat, multiplier: 2})

	mustNext(t, sequence, maxDuration)
	mustNext(t, sequence, maxDuration)
}

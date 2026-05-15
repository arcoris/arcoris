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
	seq := Exponential(time.Second, 2).NewSequence()

	mustNext(t, seq, time.Second)
	mustNext(t, seq, 2*time.Second)
	mustNext(t, seq, 4*time.Second)
}

func TestExponentialKeepsFloatStateBetweenIntegerDurations(t *testing.T) {
	seq := Exponential(2*time.Nanosecond, 1.5).NewSequence()

	mustNext(t, seq, 2*time.Nanosecond)
	mustNext(t, seq, 3*time.Nanosecond)
	mustNext(t, seq, 4*time.Nanosecond)
	mustNext(t, seq, 6*time.Nanosecond)
}

func TestExponentialSequencesHaveIndependentState(t *testing.T) {
	sched := Exponential(time.Second, 2)

	l := sched.NewSequence()
	r := sched.NewSequence()

	mustNext(t, l, time.Second)
	mustNext(t, l, 2*time.Second)
	mustNext(t, r, time.Second)
}

func TestExponentialSequenceSaturates(t *testing.T) {
	seq := &exponentialSequence{next: maxDurationFloat, multiplier: 2}

	mustNext(t, seq, maxDuration)
	mustNext(t, seq, maxDuration)
}

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

func TestPositiveJitterRejectsInvalidInput(t *testing.T) {
	mustPanicWith(t, errNilJitterSchedule, func() {
		PositiveJitter(nil, 0.1)
	})
	mustPanicWith(t, errInvalidJitterFactor, func() {
		PositiveJitter(Fixed(time.Second), -0.1)
	})
	mustPanicWith(t, errInvalidJitterFactor, func() {
		PositiveJitter(Fixed(time.Second), math.NaN())
	})
	mustPanicWith(t, errInvalidJitterFactor, func() {
		PositiveJitter(Fixed(time.Second), math.Inf(1))
	})
}

func TestPositiveJitterReturnsValueInsidePositiveRange(t *testing.T) {
	sequence := PositiveJitter(Fixed(10*time.Second), 0.5, WithRandom(fixedRandom(5*time.Second))).NewSequence()

	mustNext(t, sequence, 15*time.Second)
}

func TestPositiveJitterFactorZeroReturnsBaseDelay(t *testing.T) {
	sequence := PositiveJitter(Fixed(10*time.Second), 0, WithRandom(fixedRandom(5*time.Second))).NewSequence()

	mustNext(t, sequence, 10*time.Second)
}

func TestPositiveJitterLeavesZeroBaseDelayAtZero(t *testing.T) {
	sequence := PositiveJitter(Fixed(0), 0.5, WithRandom(fixedRandom(99))).NewSequence()

	mustNext(t, sequence, 0)
}

func TestPositiveJitterPreservesChildExhaustion(t *testing.T) {
	sequence := PositiveJitter(Delays(), 0.5, WithRandom(fixedRandom(0))).NewSequence()

	mustExhausted(t, sequence)
}

func TestPositiveJitterTransformSaturates(t *testing.T) {
	transform := positiveJitterTransform(1)

	if got := transform(maxDuration, fixedRandom(1)); got != maxDuration {
		t.Fatalf("positiveJitterTransform() = %s, want %s", got, maxDuration)
	}
}

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

package jitter

import (
	"math"
	"testing"
	"time"

	"arcoris.dev/chrono/delay"
)

func TestProportionalRejectsInvalidInput(t *testing.T) {
	mustPanicWith(t, errNilJitterSchedule, func() {
		Proportional(nil, 0.1)
	})
	mustPanicWith(t, errInvalidJitterRatio, func() {
		Proportional(delay.Fixed(time.Second), -0.1)
	})
	mustPanicWith(t, errInvalidJitterRatio, func() {
		Proportional(delay.Fixed(time.Second), 1.1)
	})
	mustPanicWith(t, errInvalidJitterRatio, func() {
		Proportional(delay.Fixed(time.Second), math.NaN())
	})
	mustPanicWith(t, errInvalidJitterRatio, func() {
		Proportional(delay.Fixed(time.Second), math.Inf(1))
	})
}

func TestProportionalReturnsValueInsideSymmetricRange(t *testing.T) {
	seq := Proportional(delay.Fixed(10*time.Second), 0.5, WithRandom(fixedRandom(10*time.Second))).NewSequence()

	mustNext(t, seq, 15*time.Second)
}

func TestProportionalCanReturnLowerBound(t *testing.T) {
	seq := Proportional(delay.Fixed(10*time.Second), 0.5, WithRandom(fixedRandom(0))).NewSequence()

	mustNext(t, seq, 5*time.Second)
}

func TestProportionalRatioZeroReturnsBaseDelay(t *testing.T) {
	seq := Proportional(delay.Fixed(10*time.Second), 0, WithRandom(fixedRandom(99))).NewSequence()

	mustNext(t, seq, 10*time.Second)
}

func TestProportionalLeavesZeroBaseDelayAtZero(t *testing.T) {
	seq := Proportional(delay.Fixed(0), 0.5, WithRandom(fixedRandom(99))).NewSequence()

	mustNext(t, seq, 0)
}

func TestProportionalPreservesChildExhaustion(t *testing.T) {
	seq := Proportional(delay.Delays(), 0.5, WithRandom(fixedRandom(0))).NewSequence()

	mustExhausted(t, seq)
}

func TestProportionalTransformSaturates(t *testing.T) {
	transform := proportionalJitterTransform(1)

	if got := transform(maxDuration, fixedRandom(0)); got != 0 {
		t.Fatalf("lower-bound proportional transform = %s, want 0", got)
	}
	if got := transform(maxDuration, fixedRandom(int64(maxDuration))); got != maxDuration {
		t.Fatalf("upper-bound proportional transform = %s, want %s", got, maxDuration)
	}
}

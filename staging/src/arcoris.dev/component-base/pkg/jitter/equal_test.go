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
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/delay"
)

func TestEqualRejectsNilSchedule(t *testing.T) {
	mustPanicWith(t, errNilJitterSchedule, func() {
		Equal(nil)
	})
}

func TestEqualReturnsValueInsideUpperHalfRange(t *testing.T) {
	sequence := Equal(delay.Fixed(10*time.Second), WithRandom(fixedRandom(5*time.Second))).NewSequence()

	mustNext(t, sequence, 10*time.Second)
}

func TestEqualCanReturnLowerBound(t *testing.T) {
	sequence := Equal(delay.Fixed(10*time.Second), WithRandom(fixedRandom(0))).NewSequence()

	mustNext(t, sequence, 5*time.Second)
}

func TestEqualLeavesZeroBaseDelayAtZero(t *testing.T) {
	sequence := Equal(delay.Fixed(0), WithRandom(fixedRandom(99))).NewSequence()

	mustNext(t, sequence, 0)
}

func TestEqualPreservesChildExhaustion(t *testing.T) {
	sequence := Equal(delay.Delays(), WithRandom(fixedRandom(0))).NewSequence()

	mustExhausted(t, sequence)
}

func TestEqualTransformUsesIntegerHalfLowerBound(t *testing.T) {
	if got := equalJitterTransform(5*time.Nanosecond, fixedRandom(3)); got != 5*time.Nanosecond {
		t.Fatalf("equalJitterTransform() = %s, want 5ns", got)
	}
}

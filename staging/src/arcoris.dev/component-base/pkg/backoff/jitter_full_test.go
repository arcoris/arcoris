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
	"testing"
	"time"
)

func TestFullJitterRejectsNilSchedule(t *testing.T) {
	mustPanicWith(t, errNilJitterSchedule, func() {
		FullJitter(nil)
	})
}

func TestFullJitterReturnsValueInsideFullRange(t *testing.T) {
	sequence := FullJitter(Fixed(10*time.Second), WithRandom(fixedRandom(10*time.Second))).NewSequence()

	mustNext(t, sequence, 10*time.Second)
}

func TestFullJitterCanReturnZero(t *testing.T) {
	sequence := FullJitter(Fixed(10*time.Second), WithRandom(fixedRandom(0))).NewSequence()

	mustNext(t, sequence, 0)
}

func TestFullJitterLeavesZeroBaseDelayAtZero(t *testing.T) {
	sequence := FullJitter(Fixed(0), WithRandom(fixedRandom(99))).NewSequence()

	mustNext(t, sequence, 0)
}

func TestFullJitterPreservesChildExhaustion(t *testing.T) {
	sequence := FullJitter(Delays(), WithRandom(fixedRandom(0))).NewSequence()

	mustExhausted(t, sequence)
}

func TestFullJitterTransformRange(t *testing.T) {
	if got := fullJitterTransform(5*time.Nanosecond, fixedRandom(3)); got != 3*time.Nanosecond {
		t.Fatalf("fullJitterTransform() = %s, want 3ns", got)
	}
}

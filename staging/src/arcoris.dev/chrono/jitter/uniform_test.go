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
)

func TestUniformRejectsInvalidBounds(t *testing.T) {
	mustPanicWith(t, errNegativeUniformMinDelay, func() {
		Uniform(-time.Nanosecond, time.Second)
	})
	mustPanicWith(t, errUniformMaxDelayBeforeMinDelay, func() {
		Uniform(time.Second, time.Millisecond)
	})
}

func TestUniformRejectsInvalidOptionsAndSource(t *testing.T) {
	mustPanicWith(t, errNilRandomOption, func() {
		Uniform(0, time.Second, nil)
	})
	mustPanicWith(t, errNilUniformSource, func() {
		randomWithSource(0, time.Second, nil)
	})
}

func TestUniformReturnsValuesInsideInclusiveRange(t *testing.T) {
	lower := Uniform(10*time.Second, 20*time.Second, WithRandom(fixedRandom(0))).NewSequence()
	upper := Uniform(10*time.Second, 20*time.Second, WithRandom(fixedRandom(int64(10*time.Second)))).NewSequence()
	middle := Uniform(10*time.Second, 20*time.Second, WithRandom(fixedRandom(int64(5*time.Second)))).NewSequence()

	mustNext(t, lower, 10*time.Second)
	mustNext(t, middle, 15*time.Second)
	mustNext(t, upper, 20*time.Second)
}

func TestUniformEqualBoundsReturnFixedDelay(t *testing.T) {
	sequence := Uniform(5*time.Second, 5*time.Second, WithRandom(fixedRandom(99))).NewSequence()

	mustNext(t, sequence, 5*time.Second)
	mustNext(t, sequence, 5*time.Second)
}

func TestUniformScheduleRequestsGeneratorPerSequence(t *testing.T) {
	source := &countingRandomSource{}
	schedule := Uniform(0, time.Second, WithRandomSource(source))

	left := schedule.NewSequence()
	right := schedule.NewSequence()

	if source.calls != 2 {
		t.Fatalf("NewRandom calls = %d, want 2", source.calls)
	}
	mustNext(t, left, 0)
	mustNext(t, right, time.Nanosecond)
}

func TestUniformSequencesWithSameSeedAreIndependentAndDeterministic(t *testing.T) {
	schedule := Uniform(time.Second, 10*time.Second, WithSeed(42))

	left := schedule.NewSequence()
	right := schedule.NewSequence()

	for i := 0; i < 5; i++ {
		leftDelay, leftOK := left.Next()
		rightDelay, rightOK := right.Next()
		if !leftOK || !rightOK {
			t.Fatalf("sequence %d exhausted: left=%v right=%v", i, leftOK, rightOK)
		}
		if leftDelay != rightDelay {
			t.Fatalf("sequence %d mismatch: left=%s right=%s", i, leftDelay, rightDelay)
		}
	}
}

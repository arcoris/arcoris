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

func TestDecorrelatedJitterRejectsInvalidInput(t *testing.T) {
	mustPanicWith(t, errNonPositiveDecorrelatedInitialDelay, func() {
		DecorrelatedJitter(0, time.Second, 2)
	})
	mustPanicWith(t, errNonPositiveDecorrelatedInitialDelay, func() {
		DecorrelatedJitter(-time.Nanosecond, time.Second, 2)
	})
	mustPanicWith(t, errDecorrelatedMaxDelayBeforeInitialDelay, func() {
		DecorrelatedJitter(time.Second, time.Millisecond, 2)
	})
	mustPanicWith(t, errInvalidDecorrelatedMultiplier, func() {
		DecorrelatedJitter(time.Second, 2*time.Second, 1)
	})
	mustPanicWith(t, errInvalidDecorrelatedMultiplier, func() {
		DecorrelatedJitter(time.Second, 2*time.Second, math.NaN())
	})
	mustPanicWith(t, errInvalidDecorrelatedMultiplier, func() {
		DecorrelatedJitter(time.Second, 2*time.Second, math.Inf(1))
	})
	mustPanicWith(t, errNilRandomOption, func() {
		DecorrelatedJitter(time.Second, 2*time.Second, 2, nil)
	})
	mustPanicWith(t, errNilDecorrelatedJitterSource, func() {
		decorrelatedJitterWithSource(time.Second, 2*time.Second, 2, nil)
	})
}

func TestDecorrelatedJitterReturnsValuesInsideExpectedRanges(t *testing.T) {
	source := RandomSourceFunc(func() RandomGenerator {
		return &sequenceRandom{values: []int64{
			int64(time.Second),
			int64(2 * time.Second),
			int64(3 * time.Second),
		}}
	})
	sequence := DecorrelatedJitter(time.Second, 5*time.Second, 2, WithRandomSource(source)).NewSequence()

	first := mustNextInRange(t, sequence, time.Second, 2*time.Second)
	second := mustNextInRange(t, sequence, time.Second, 4*time.Second)
	third := mustNextInRange(t, sequence, time.Second, 5*time.Second)

	if first != 2*time.Second || second != 3*time.Second || third != 4*time.Second {
		t.Fatalf("delays = %s, %s, %s; want 2s, 3s, 4s", first, second, third)
	}
}

func TestDecorrelatedJitterSequencesHaveIndependentPreviousState(t *testing.T) {
	source := RandomSourceFunc(func() RandomGenerator {
		return &sequenceRandom{values: []int64{int64(time.Second), int64(2 * time.Second)}}
	})
	schedule := DecorrelatedJitter(time.Second, 5*time.Second, 2, WithRandomSource(source))

	left := schedule.NewSequence()
	right := schedule.NewSequence()

	mustNext(t, left, 2*time.Second)
	mustNext(t, left, 3*time.Second)
	mustNext(t, right, 2*time.Second)
}

func TestDecorrelatedJitterUpperBoundUsesMaxDelay(t *testing.T) {
	sequence := &decorrelatedJitterSequence{
		initial:    time.Second,
		maxDelay:   3 * time.Second,
		multiplier: 10,
		previous:   2 * time.Second,
		random:     fixedRandom(0),
	}

	if got := sequence.upperBound(); got != 3*time.Second {
		t.Fatalf("upperBound() = %s, want 3s", got)
	}
}

func TestDecorrelatedJitterUpperBoundSaturates(t *testing.T) {
	sequence := &decorrelatedJitterSequence{
		initial:    time.Second,
		maxDelay:   maxDuration,
		multiplier: 2,
		previous:   maxDuration,
		random:     fixedRandom(0),
	}

	if got := sequence.upperBound(); got != maxDuration {
		t.Fatalf("upperBound() = %s, want %s", got, maxDuration)
	}
}

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
)

func TestDecorrelatedRejectsInvalidInput(t *testing.T) {
	mustPanicWith(t, errNonPositiveDecorrelatedInitialDelay, func() {
		Decorrelated(0, time.Second, 2)
	})
	mustPanicWith(t, errNonPositiveDecorrelatedInitialDelay, func() {
		Decorrelated(-time.Nanosecond, time.Second, 2)
	})
	mustPanicWith(t, errDecorrelatedMaxDelayBeforeInitialDelay, func() {
		Decorrelated(time.Second, time.Millisecond, 2)
	})
	mustPanicWith(t, errInvalidDecorrelatedMultiplier, func() {
		Decorrelated(time.Second, 2*time.Second, 1)
	})
	mustPanicWith(t, errInvalidDecorrelatedMultiplier, func() {
		Decorrelated(time.Second, 2*time.Second, math.NaN())
	})
	mustPanicWith(t, errInvalidDecorrelatedMultiplier, func() {
		Decorrelated(time.Second, 2*time.Second, math.Inf(1))
	})
	mustPanicWith(t, errNilRandomOption, func() {
		Decorrelated(time.Second, 2*time.Second, 2, nil)
	})
	mustPanicWith(t, errNilDecorrelatedSource, func() {
		decorrelatedJitterWithSource(time.Second, 2*time.Second, 2, nil)
	})
}

func TestDecorrelatedReturnsValuesInsideExpectedRanges(t *testing.T) {
	source := RandomSourceFunc(func() RandomGenerator {
		return &sequenceRandom{values: []int64{
			int64(time.Second),
			int64(2 * time.Second),
			int64(3 * time.Second),
		}}
	})
	seq := Decorrelated(time.Second, 5*time.Second, 2, WithRandomSource(source)).NewSequence()

	first := mustNextInRange(t, seq, time.Second, 2*time.Second)
	second := mustNextInRange(t, seq, time.Second, 4*time.Second)
	third := mustNextInRange(t, seq, time.Second, 5*time.Second)

	if first != 2*time.Second || second != 3*time.Second || third != 4*time.Second {
		t.Fatalf("delays = %s, %s, %s; want 2s, 3s, 4s", first, second, third)
	}
}

func TestDecorrelatedSequencesHaveIndependentPreviousState(t *testing.T) {
	source := RandomSourceFunc(func() RandomGenerator {
		return &sequenceRandom{values: []int64{int64(time.Second), int64(2 * time.Second)}}
	})
	sched := Decorrelated(time.Second, 5*time.Second, 2, WithRandomSource(source))

	l := sched.NewSequence()
	r := sched.NewSequence()

	mustNext(t, l, 2*time.Second)
	mustNext(t, l, 3*time.Second)
	mustNext(t, r, 2*time.Second)
}

func TestDecorrelatedUpperBoundUsesMaxDelay(t *testing.T) {
	seq := &decorrelatedJitterSequence{
		initial:    time.Second,
		maxDelay:   3 * time.Second,
		multiplier: 10,
		previous:   2 * time.Second,
		random:     fixedRandom(0),
	}

	if got := seq.upperBound(); got != 3*time.Second {
		t.Fatalf("upperBound() = %s, want 3s", got)
	}
}

func TestDecorrelatedUpperBoundSaturates(t *testing.T) {
	seq := &decorrelatedJitterSequence{
		initial:    time.Second,
		maxDelay:   maxDuration,
		multiplier: 2,
		previous:   maxDuration,
		random:     fixedRandom(0),
	}

	if got := seq.upperBound(); got != maxDuration {
		t.Fatalf("upperBound() = %s, want %s", got, maxDuration)
	}
}

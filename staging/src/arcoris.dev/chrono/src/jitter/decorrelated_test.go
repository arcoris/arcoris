// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jitter

import (
	panicassert "arcoris.dev/testutil/panic"
	"math"
	"testing"
	"time"
)

func TestDecorrelatedRejectsInvalidInput(t *testing.T) {
	panicassert.RequireValue(t, errNonPositiveDecorrelatedInitialDelay, func() {
		Decorrelated(0, time.Second, 2)
	})
	panicassert.RequireValue(t, errNonPositiveDecorrelatedInitialDelay, func() {
		Decorrelated(-time.Nanosecond, time.Second, 2)
	})
	panicassert.RequireValue(t, errDecorrelatedMaxDelayBeforeInitialDelay, func() {
		Decorrelated(time.Second, time.Millisecond, 2)
	})
	panicassert.RequireValue(t, errInvalidDecorrelatedMultiplier, func() {
		Decorrelated(time.Second, 2*time.Second, 1)
	})
	panicassert.RequireValue(t, errInvalidDecorrelatedMultiplier, func() {
		Decorrelated(time.Second, 2*time.Second, 0.5)
	})
	panicassert.RequireValue(t, errInvalidDecorrelatedMultiplier, func() {
		Decorrelated(time.Second, 2*time.Second, math.NaN())
	})
	panicassert.RequireValue(t, errInvalidDecorrelatedMultiplier, func() {
		Decorrelated(time.Second, 2*time.Second, math.Inf(1))
	})
	panicassert.RequireValue(t, errInvalidDecorrelatedMultiplier, func() {
		Decorrelated(time.Second, 2*time.Second, math.Inf(-1))
	})
	panicassert.RequireValue(t, errNilRandomOption, func() {
		Decorrelated(time.Second, 2*time.Second, 2, nil)
	})
	panicassert.RequireValue(t, errNilDecorrelatedSource, func() {
		decorrelatedJitterWithSource(time.Second, 2*time.Second, 2, nil)
	})
}

func TestDecorrelatedRejectsNilRandomGenerator(t *testing.T) {
	panicassert.RequireValue(t, errNilRandom, func() {
		decorrelatedJitterWithSource(time.Second, 2*time.Second, 2, nilRandomSource{}).NewSequence()
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

func TestDecorrelatedCanReturnRangeBounds(t *testing.T) {
	source := RandomSourceFunc(func() RandomGenerator {
		return &sequenceRandom{values: []int64{0, int64(time.Second)}}
	})
	seq := Decorrelated(time.Second, 5*time.Second, 2, WithRandomSource(source)).NewSequence()

	mustNext(t, seq, time.Second)
	mustNext(t, seq, 2*time.Second)
}

func TestDecorrelatedEqualBoundsReturnInitialWithoutDrawing(t *testing.T) {
	seq := Decorrelated(time.Second, time.Second, 2, WithRandomFunc(func() int64 {
		return -1
	})).NewSequence()

	mustNext(t, seq, time.Second)
	mustNext(t, seq, time.Second)
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

func TestDecorrelatedMultiplierTruncatesFractionalNanoseconds(t *testing.T) {
	seq := Decorrelated(5*time.Nanosecond, 10*time.Nanosecond, 1.5, WithRandom(fixedRandom(2))).NewSequence()

	mustNext(t, seq, 7*time.Nanosecond)
}

func TestDecorrelatedNeverReturnsAboveMaxDelay(t *testing.T) {
	seq := Decorrelated(time.Second, 3*time.Second, 10, WithRandom(fixedRandom(2*time.Second))).NewSequence()

	for i := 0; i < 16; i++ {
		got := mustNextInRange(t, seq, time.Second, 3*time.Second)
		if got > 3*time.Second {
			t.Fatalf("Next() delay = %s, want <= 3s", got)
		}
	}
}

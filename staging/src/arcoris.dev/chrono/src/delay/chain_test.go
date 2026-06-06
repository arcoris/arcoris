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

package delay

import (
	panicassert "arcoris.dev/testutil/panic"
	"testing"
	"time"
)

func TestChainEmptyIsExhaustedImmediately(t *testing.T) {
	seq := Chain().NewSequence()

	mustExhausted(t, seq)
	mustExhausted(t, seq)
}

func TestChainSingleChildBehavesLikeChild(t *testing.T) {
	seq := Chain(Delays(time.Second, 2*time.Second)).NewSequence()

	mustNext(t, seq, time.Second)
	mustNext(t, seq, 2*time.Second)
	mustExhausted(t, seq)
}

func TestChainConsumesFiniteChildrenInOrder(t *testing.T) {
	seq := Chain(
		Delays(0, time.Second),
		Delays(2*time.Second),
		Delays(3*time.Second, 4*time.Second),
	).NewSequence()

	mustNext(t, seq, 0)
	mustNext(t, seq, time.Second)
	mustNext(t, seq, 2*time.Second)
	mustNext(t, seq, 3*time.Second)
	mustNext(t, seq, 4*time.Second)
	mustExhausted(t, seq)
}

func TestChainFinitePrefixThenInfiniteTail(t *testing.T) {
	seq := Chain(
		Delays(0, 10*time.Millisecond),
		Exponential(100*time.Millisecond, 2),
	).NewSequence()

	mustNext(t, seq, 0)
	mustNext(t, seq, 10*time.Millisecond)
	mustNext(t, seq, 100*time.Millisecond)
	mustNext(t, seq, 200*time.Millisecond)
	mustNext(t, seq, 400*time.Millisecond)
}

func TestChainInfiniteChildPreventsLaterConsumption(t *testing.T) {
	later := &countingSequence{values: []time.Duration{time.Hour}}
	seq := Chain(
		Fixed(time.Second),
		ScheduleFunc(func() Sequence { return later }),
	).NewSequence()

	mustNext(t, seq, time.Second)
	mustNext(t, seq, time.Second)
	mustNext(t, seq, time.Second)

	if later.calls != 0 {
		t.Fatalf("later child calls = %d, want 0", later.calls)
	}
}

func TestChainNewSequenceCreatesIndependentChildSequences(t *testing.T) {
	sched := Chain(
		Delays(time.Second, 2*time.Second),
		Fixed(3*time.Second),
	)

	left := sched.NewSequence()
	right := sched.NewSequence()

	mustNext(t, left, time.Second)
	mustNext(t, left, 2*time.Second)
	mustNext(t, left, 3*time.Second)

	mustNext(t, right, time.Second)
	mustNext(t, right, 2*time.Second)
}

func TestChainCopiesCallerScheduleSlice(t *testing.T) {
	children := []Schedule{
		Delays(time.Second),
		Delays(2 * time.Second),
		Fixed(3 * time.Second),
	}
	sched := Chain(children...)
	children[0] = Fixed(time.Hour)
	children[2] = Fixed(time.Hour)

	seq := sched.NewSequence()

	mustNext(t, seq, time.Second)
	mustNext(t, seq, 2*time.Second)
	mustNext(t, seq, 3*time.Second)
}

func TestChainRejectsNilChildSchedule(t *testing.T) {
	panicassert.RequireValue(t, errNilChainSchedule, func() {
		Chain(Delays(time.Second), nil)
	})
}

func TestChainRejectsNilChildSequence(t *testing.T) {
	sched := Chain(Delays(time.Second), nilSequenceSchedule{})

	panicassert.RequireValue(t, errChainScheduleReturnedNilSequence, func() {
		sched.NewSequence()
	})
}

func TestChainRejectsNegativeAvailableChildDelay(t *testing.T) {
	seq := Chain(
		ScheduleFunc(func() Sequence { return negativeDelaySequence{} }),
	).NewSequence()

	panicassert.RequireValue(t, errChainScheduleReturnedNegativeDelay, func() {
		seq.Next()
	})
}

func TestChainIgnoresNegativeDelayAfterChildExhaustion(t *testing.T) {
	seq := Chain(
		ScheduleFunc(func() Sequence { return exhaustedNegativeSequence{} }),
		Delays(time.Second),
	).NewSequence()

	mustNext(t, seq, time.Second)
	mustExhausted(t, seq)
}

func TestChainPreservesZeroDelays(t *testing.T) {
	seq := Chain(
		Delays(0),
		Fixed(0),
	).NewSequence()

	mustNext(t, seq, 0)
	mustNext(t, seq, 0)
	mustNext(t, seq, 0)
}

func TestChainDoesNotCallExhaustedChildrenAgain(t *testing.T) {
	first := &countingSequence{values: []time.Duration{time.Second}}
	second := &countingSequence{values: []time.Duration{2 * time.Second}}
	seq := Chain(
		ScheduleFunc(func() Sequence { return first }),
		ScheduleFunc(func() Sequence { return second }),
	).NewSequence()

	mustNext(t, seq, time.Second)
	mustNext(t, seq, 2*time.Second)
	mustExhausted(t, seq)
	mustExhausted(t, seq)

	if first.calls != 2 {
		t.Fatalf("first child calls = %d, want 2", first.calls)
	}
	if second.calls != 2 {
		t.Fatalf("second child calls = %d, want 2", second.calls)
	}
}

func TestSingleChildChainDoesNotCallChildAfterExhaustion(t *testing.T) {
	child := &countingSequence{values: []time.Duration{time.Second}}
	seq := Chain(
		ScheduleFunc(func() Sequence { return child }),
	).NewSequence()

	mustNext(t, seq, time.Second)
	mustExhausted(t, seq)
	mustExhausted(t, seq)

	if child.calls != 2 {
		t.Fatalf("child calls = %d, want 2", child.calls)
	}
}

func TestChainSkipsEmptyChildrenBeforeNonEmptyChild(t *testing.T) {
	seq := Chain(
		Delays(),
		Delays(),
		Delays(time.Second),
	).NewSequence()

	mustNext(t, seq, time.Second)
	mustExhausted(t, seq)
}

func TestChainComposesWithCap(t *testing.T) {
	seq := Cap(
		Chain(
			Delays(5*time.Second),
			Fixed(10*time.Second),
		),
		time.Second,
	).NewSequence()

	mustNext(t, seq, time.Second)
	mustNext(t, seq, time.Second)
}

func TestChainComposesWithLimit(t *testing.T) {
	seq := Limit(
		Chain(
			Delays(0),
			Fixed(time.Second),
		),
		3,
	).NewSequence()

	mustNext(t, seq, 0)
	mustNext(t, seq, time.Second)
	mustNext(t, seq, time.Second)
	mustExhausted(t, seq)
}

func TestChainComposesWithScheduleFuncAndSequenceFunc(t *testing.T) {
	calls := 0
	sched := Chain(
		ScheduleFunc(func() Sequence {
			return SequenceFunc(func() (time.Duration, bool) {
				if calls > 0 {
					return 0, false
				}
				calls++
				return time.Second, true
			})
		}),
		Delays(2*time.Second),
	)

	seq := sched.NewSequence()

	mustNext(t, seq, time.Second)
	mustNext(t, seq, 2*time.Second)
	mustExhausted(t, seq)
}

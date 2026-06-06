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

func TestLimitRejectsInvalidInput(t *testing.T) {
	panicassert.RequireValue(t, errNilLimitSchedule, func() {
		Limit(nil, 1)
	})
	panicassert.RequireValue(t, errNegativeLimitMaxDelays, func() {
		Limit(Fixed(time.Second), -1)
	})
}

func TestLimitExposesOnlyConfiguredNumberOfValues(t *testing.T) {
	seq := Limit(Fixed(time.Second), 2).NewSequence()

	mustNext(t, seq, time.Second)
	mustNext(t, seq, time.Second)
	mustExhausted(t, seq)
}

func TestLimitSequencesHaveIndependentRemainingCounts(t *testing.T) {
	sched := Limit(Fixed(time.Second), 2)

	left := sched.NewSequence()
	right := sched.NewSequence()

	mustNext(t, left, time.Second)
	mustNext(t, left, time.Second)
	mustExhausted(t, left)

	mustNext(t, right, time.Second)
	mustNext(t, right, time.Second)
}

func TestLimitPreservesEarlyChildExhaustion(t *testing.T) {
	child := &countingSequence{values: []time.Duration{time.Second}}
	seq := Limit(ScheduleFunc(func() Sequence { return child }), 3).NewSequence()

	mustNext(t, seq, time.Second)
	mustExhausted(t, seq)
	mustExhausted(t, seq)
	if child.calls != 2 {
		t.Fatalf("child calls = %d, want 2", child.calls)
	}
}

func TestLimitZeroDoesNotCreateChildSequence(t *testing.T) {
	called := false
	seq := Limit(ScheduleFunc(func() Sequence {
		called = true
		return Immediate().NewSequence()
	}), 0).NewSequence()

	mustExhausted(t, seq)
	if called {
		t.Fatal("zero limit created child sequence")
	}
}

func TestLimitDoesNotCallChildAfterExhaustion(t *testing.T) {
	child := &countingSequence{values: []time.Duration{time.Second, 2 * time.Second}}
	seq := Limit(ScheduleFunc(func() Sequence { return child }), 1).NewSequence()

	mustNext(t, seq, time.Second)
	mustExhausted(t, seq)
	mustExhausted(t, seq)
	if child.calls != 1 {
		t.Fatalf("child calls = %d, want 1", child.calls)
	}
}

func TestLimitRejectsNilChildSequence(t *testing.T) {
	panicassert.RequireValue(t, errLimitScheduleReturnedNilSequence, func() {
		Limit(nilSequenceSchedule{}, 1).NewSequence()
	})
}

func TestLimitRejectsNegativeChildDelay(t *testing.T) {
	seq := Limit(ScheduleFunc(func() Sequence { return negativeDelaySequence{} }), 1).NewSequence()

	panicassert.RequireValue(t, errLimitScheduleReturnedNegativeDelay, func() {
		seq.Next()
	})
}

func TestLimitIgnoresNegativeDelayAfterChildExhaustion(t *testing.T) {
	seq := Limit(ScheduleFunc(func() Sequence { return exhaustedNegativeSequence{} }), 1).NewSequence()

	mustExhausted(t, seq)
}

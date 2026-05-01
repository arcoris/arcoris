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

func TestLimitRejectsInvalidInput(t *testing.T) {
	mustPanicWith(t, errNilLimitSchedule, func() {
		Limit(nil, 1)
	})
	mustPanicWith(t, errNegativeLimitMaxDelays, func() {
		Limit(Fixed(time.Second), -1)
	})
}

func TestLimitReturnsPointerSequence(t *testing.T) {
	schedule := Limit(Fixed(time.Second), 1)

	if _, ok := schedule.NewSequence().(*limitSequence); !ok {
		t.Fatalf("NewSequence() = %T, want *limitSequence", schedule.NewSequence())
	}
}

func TestLimitExposesOnlyConfiguredNumberOfValues(t *testing.T) {
	sequence := Limit(Fixed(time.Second), 2).NewSequence()

	mustNext(t, sequence, time.Second)
	mustNext(t, sequence, time.Second)
	mustExhausted(t, sequence)
}

func TestLimitPreservesEarlyChildExhaustion(t *testing.T) {
	sequence := Limit(Delays(time.Second), 3).NewSequence()

	mustNext(t, sequence, time.Second)
	mustExhausted(t, sequence)
}

func TestLimitZeroDoesNotCreateChildSequence(t *testing.T) {
	called := false
	sequence := Limit(ScheduleFunc(func() Sequence {
		called = true
		return Immediate().NewSequence()
	}), 0).NewSequence()

	mustExhausted(t, sequence)
	if called {
		t.Fatal("zero limit created child sequence")
	}
}

func TestLimitDoesNotCallChildAfterExhaustion(t *testing.T) {
	child := &countingSequence{values: []time.Duration{time.Second, 2 * time.Second}}
	sequence := Limit(ScheduleFunc(func() Sequence { return child }), 1).NewSequence()

	mustNext(t, sequence, time.Second)
	mustExhausted(t, sequence)
	mustExhausted(t, sequence)
	if child.calls != 1 {
		t.Fatalf("child calls = %d, want 1", child.calls)
	}
}

func TestLimitRejectsNilChildSequence(t *testing.T) {
	mustPanicWith(t, errLimitScheduleReturnedNilSequence, func() {
		Limit(nilSequenceSchedule{}, 1).NewSequence()
	})
}

func TestLimitRejectsNegativeChildDelay(t *testing.T) {
	sequence := Limit(ScheduleFunc(func() Sequence { return negativeDelaySequence{} }), 1).NewSequence()

	mustPanicWith(t, errLimitScheduleReturnedNegativeDelay, func() {
		sequence.Next()
	})
}

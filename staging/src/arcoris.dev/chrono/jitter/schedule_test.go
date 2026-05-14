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

	"arcoris.dev/chrono/delay"
)

func TestNewJitterScheduleRejectsInvalidInput(t *testing.T) {
	mustPanicWith(t, errNilJitterSchedule, func() {
		newJitterSchedule(nil, fullJitterTransform)
	})
	mustPanicWith(t, errNilJitterTransform, func() {
		newJitterSchedule(delay.Fixed(time.Second), nil)
	})
	mustPanicWith(t, errNilRandomOption, func() {
		newJitterSchedule(delay.Fixed(time.Second), fullJitterTransform, nil)
	})
	mustPanicWith(t, errNilRandomSource, func() {
		newJitterScheduleWithSource(delay.Fixed(time.Second), fullJitterTransform, nil)
	})
}

func TestJitterPreservesChildExhaustion(t *testing.T) {
	sequence := newJitterSchedule(delay.Delays(), fullJitterTransform, WithRandom(fixedRandom(0))).NewSequence()

	mustExhausted(t, sequence)
}

func TestJitterRejectsNilChildSequence(t *testing.T) {
	mustPanicWith(t, errJitterScheduleReturnedNilSequence, func() {
		newJitterSchedule(nilSequenceSchedule{}, fullJitterTransform).NewSequence()
	})
}

func TestJitterRejectsNegativeChildDelay(t *testing.T) {
	sequence := newJitterSchedule(delay.ScheduleFunc(func() delay.Sequence { return negativeDelaySequence{} }), fullJitterTransform).NewSequence()

	mustPanicWith(t, errJitterScheduleReturnedNegativeDelay, func() {
		sequence.Next()
	})
}

func TestJitterRejectsNegativeTransformOutput(t *testing.T) {
	sequence := newJitterSchedule(delay.Fixed(time.Second), func(time.Duration, RandomGenerator) time.Duration {
		return -time.Nanosecond
	}).NewSequence()

	mustPanicWith(t, errJitterTransformReturnedNegativeDelay, func() {
		sequence.Next()
	})
}

func TestJitterRequestsRandomGeneratorPerSequence(t *testing.T) {
	source := &countingRandomSource{}
	schedule := newJitterScheduleWithSource(delay.Fixed(time.Second), fullJitterTransform, source)

	left := schedule.NewSequence()
	right := schedule.NewSequence()

	if source.calls != 2 {
		t.Fatalf("NewRandom calls = %d, want 2", source.calls)
	}
	mustNext(t, left, 0)
	mustNext(t, right, time.Nanosecond)
}

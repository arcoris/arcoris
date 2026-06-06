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
	"fmt"
	"sync"
	"testing"
	"time"

	"arcoris.dev/chrono/delay"
)

func TestNewJitterScheduleRejectsInvalidInput(t *testing.T) {
	panicassert.RequireValue(t, errNilJitterSchedule, func() {
		newJitterSchedule(nil, fullJitterTransform)
	})
	panicassert.RequireValue(t, errNilJitterTransform, func() {
		newJitterSchedule(delay.Fixed(time.Second), nil)
	})
	panicassert.RequireValue(t, errNilRandomOption, func() {
		newJitterSchedule(delay.Fixed(time.Second), fullJitterTransform, nil)
	})
	panicassert.RequireValue(t, errNilRandomSource, func() {
		newJitterScheduleWithSource(delay.Fixed(time.Second), fullJitterTransform, nil)
	})
}

func TestJitterPreservesChildExhaustion(t *testing.T) {
	seq := newJitterSchedule(delay.Delays(), fullJitterTransform, WithRandom(fixedRandom(0))).NewSequence()

	mustExhausted(t, seq)
}

func TestJitterRejectsNilChildSequence(t *testing.T) {
	panicassert.RequireValue(t, errJitterScheduleReturnedNilSequence, func() {
		newJitterSchedule(nilSequenceSchedule{}, fullJitterTransform).NewSequence()
	})
}

func TestJitterRejectsNilRandomGenerator(t *testing.T) {
	panicassert.RequireValue(t, errNilRandom, func() {
		newJitterScheduleWithSource(delay.Fixed(time.Second), fullJitterTransform, nilRandomSource{}).NewSequence()
	})
}

func TestJitterRejectsNegativeChildDelay(t *testing.T) {
	seq := newJitterSchedule(delay.ScheduleFunc(func() delay.Sequence { return negativeDelaySequence{} }), fullJitterTransform).NewSequence()

	panicassert.RequireValue(t, errJitterScheduleReturnedNegativeDelay, func() {
		seq.Next()
	})
}

func TestJitterIgnoresNegativeDelayAfterChildExhaustion(t *testing.T) {
	seq := newJitterSchedule(delay.ScheduleFunc(func() delay.Sequence {
		return exhaustedNegativeSequence{}
	}), fullJitterTransform).NewSequence()

	mustExhausted(t, seq)
}

func TestJitterRejectsNegativeRandomOutputWhenDrawn(t *testing.T) {
	seq := Uniform(0, time.Nanosecond, WithRandomFunc(func() int64 {
		return -1
	})).NewSequence()

	panicassert.RequireValue(t, errRandomReturnedNegative, func() {
		seq.Next()
	})
}

func TestJitterDoesNotDrawRandomWhenNoDrawIsNeeded(t *testing.T) {
	negativeRandom := WithRandomFunc(func() int64 {
		return -1
	})

	mustNext(t, Uniform(time.Second, time.Second, negativeRandom).NewSequence(), time.Second)
	mustNext(t, Full(delay.Fixed(0), negativeRandom).NewSequence(), 0)
	mustExhausted(t, Full(delay.Delays(), negativeRandom).NewSequence())
}

func TestJitterRejectsNegativeTransformOutput(t *testing.T) {
	seq := newJitterSchedule(delay.Fixed(time.Second), func(time.Duration, RandomGenerator) time.Duration {
		return -time.Nanosecond
	}).NewSequence()

	panicassert.RequireValue(t, errJitterTransformReturnedNegativeDelay, func() {
		seq.Next()
	})
}

func TestJitterRequestsRandomGeneratorPerSequence(t *testing.T) {
	source := &countingRandomSource{}
	sched := newJitterScheduleWithSource(delay.Fixed(time.Second), fullJitterTransform, source)

	l := sched.NewSequence()
	r := sched.NewSequence()

	if source.calls != 2 {
		t.Fatalf("NewRandom calls = %d, want 2", source.calls)
	}
	mustNext(t, l, 0)
	mustNext(t, r, time.Nanosecond)
}

func TestJitterSchedulesNeverReturnNegativeDelaysForValidInputs(t *testing.T) {
	tests := []struct {
		name  string
		sched delay.Schedule
	}{
		{name: "full", sched: Full(delay.Fixed(time.Second), WithRandom(fixedRandom(0)))},
		{name: "equal", sched: Equal(delay.Fixed(time.Second), WithRandom(fixedRandom(0)))},
		{name: "positive", sched: Positive(delay.Fixed(time.Second), 0.5, WithRandom(fixedRandom(0)))},
		{name: "proportional", sched: Proportional(delay.Fixed(time.Second), 0.5, WithRandom(fixedRandom(0)))},
		{name: "uniform", sched: Uniform(0, time.Second, WithRandom(fixedRandom(0)))},
		{name: "decorrelated", sched: Decorrelated(time.Second, 5*time.Second, 2, WithRandom(fixedRandom(0)))},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			seq := tc.sched.NewSequence()
			for i := 0; i < 3; i++ {
				got, ok := seq.Next()
				if !ok {
					t.Fatal("Next() exhausted, want available delay")
				}
				if got < 0 {
					t.Fatalf("Next() delay = %s, want non-negative", got)
				}
			}
		})
	}
}

func TestJitterSchedulesSupportConcurrentNewSequence(t *testing.T) {
	tests := []struct {
		name  string
		sched delay.Schedule
	}{
		{name: "full", sched: Full(delay.Fixed(time.Second), WithSeed(1))},
		{name: "equal", sched: Equal(delay.Fixed(time.Second), WithSeed(1))},
		{name: "positive", sched: Positive(delay.Fixed(time.Second), 0.2, WithSeed(1))},
		{name: "proportional", sched: Proportional(delay.Fixed(time.Second), 0.2, WithSeed(1))},
		{name: "uniform", sched: Uniform(0, time.Second, WithSeed(1))},
		{name: "decorrelated", sched: Decorrelated(time.Second, 10*time.Second, 2, WithSeed(1))},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			want := nextPrefix(t, tc.sched.NewSequence(), 4)

			const (
				workers            = 16
				sequencesPerWorker = 32
			)

			errs := make(chan string, workers)

			var wg sync.WaitGroup
			wg.Add(workers)

			for worker := 0; worker < workers; worker++ {
				go func() {
					defer wg.Done()

					for i := 0; i < sequencesPerWorker; i++ {
						seq := tc.sched.NewSequence()
						for step, wantDelay := range want {
							got, ok := seq.Next()
							if !ok || got != wantDelay {
								errs <- fmt.Sprintf(
									"step %d: Next() = %s, %t; want %s, true",
									step,
									got,
									ok,
									wantDelay,
								)
								return
							}
						}
					}
				}()
			}

			wg.Wait()
			close(errs)

			for err := range errs {
				t.Error(err)
			}
		})
	}
}

func nextPrefix(t *testing.T, seq delay.Sequence, n int) []time.Duration {
	t.Helper()

	values := make([]time.Duration, n)
	for i := range values {
		got, ok := seq.Next()
		if !ok {
			t.Fatalf("Next() exhausted at %d, want available delay", i)
		}
		values[i] = got
	}

	return values
}

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
	"sync"
	"testing"
	"time"
)

func TestSchedulesCreateIndependentSequences(t *testing.T) {
	tests := []struct {
		name     string
		schedule Schedule
		first    time.Duration
		second   time.Duration
	}{
		{
			name:     "immediate",
			schedule: Immediate(),
			first:    0,
			second:   0,
		},
		{
			name:     "fixed",
			schedule: Fixed(time.Second),
			first:    time.Second,
			second:   time.Second,
		},
		{
			name:     "delays",
			schedule: Delays(time.Second, 2*time.Second),
			first:    time.Second,
			second:   2 * time.Second,
		},
		{
			name:     "linear",
			schedule: Linear(time.Second, time.Second),
			first:    time.Second,
			second:   2 * time.Second,
		},
		{
			name:     "exponential",
			schedule: Exponential(time.Second, 2),
			first:    time.Second,
			second:   2 * time.Second,
		},
		{
			name:     "fibonacci",
			schedule: Fibonacci(time.Second),
			first:    time.Second,
			second:   time.Second,
		},
		{
			name:     "cap",
			schedule: Cap(Linear(time.Second, time.Second), 5*time.Second),
			first:    time.Second,
			second:   2 * time.Second,
		},
		{
			name:     "limit",
			schedule: Limit(Fixed(time.Second), 2),
			first:    time.Second,
			second:   time.Second,
		},
		{
			name:     "chain",
			schedule: Chain(Delays(time.Second), Fixed(2*time.Second)),
			first:    time.Second,
			second:   2 * time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			left := tc.schedule.NewSequence()
			right := tc.schedule.NewSequence()

			mustNext(t, left, tc.first)
			mustNext(t, left, tc.second)
			mustNext(t, right, tc.first)
		})
	}
}

func TestBuiltInSchedulesSupportConcurrentNewSequence(t *testing.T) {
	tests := []struct {
		name     string
		schedule Schedule
		want     []time.Duration
	}{
		{
			name:     "immediate",
			schedule: Immediate(),
			want:     []time.Duration{0, 0},
		},
		{
			name:     "fixed",
			schedule: Fixed(time.Second),
			want:     []time.Duration{time.Second, time.Second},
		},
		{
			name:     "delays",
			schedule: Delays(0, time.Millisecond, time.Second),
			want:     []time.Duration{0, time.Millisecond, time.Second},
		},
		{
			name:     "linear",
			schedule: Linear(0, time.Millisecond),
			want:     []time.Duration{0, time.Millisecond, 2 * time.Millisecond},
		},
		{
			name:     "exponential",
			schedule: Exponential(time.Millisecond, 2),
			want:     []time.Duration{time.Millisecond, 2 * time.Millisecond},
		},
		{
			name:     "fibonacci",
			schedule: Fibonacci(time.Millisecond),
			want:     []time.Duration{time.Millisecond, time.Millisecond, 2 * time.Millisecond},
		},
		{
			name:     "cap",
			schedule: Cap(Exponential(time.Millisecond, 2), 2*time.Millisecond),
			want:     []time.Duration{time.Millisecond, 2 * time.Millisecond, 2 * time.Millisecond},
		},
		{
			name:     "limit",
			schedule: Limit(Fixed(time.Second), 3),
			want:     []time.Duration{time.Second, time.Second, time.Second},
		},
		{
			name:     "chain",
			schedule: Chain(Delays(0), Exponential(time.Millisecond, 2)),
			want:     []time.Duration{0, time.Millisecond, 2 * time.Millisecond},
		},
		{
			name: "large chain",
			schedule: Chain(
				Delays(0),
				Delays(time.Millisecond),
				Fixed(2*time.Millisecond),
			),
			want: []time.Duration{0, time.Millisecond, 2 * time.Millisecond},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var wg sync.WaitGroup

			for worker := 0; worker < 16; worker++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					for i := 0; i < 64; i++ {
						seq := tc.schedule.NewSequence()
						for _, want := range tc.want {
							got, ok := seq.Next()
							if !ok || got != want {
								t.Errorf("Next() = %s, %t; want %s, true", got, ok, want)
								return
							}
						}
					}
				}()
			}

			wg.Wait()
		})
	}
}

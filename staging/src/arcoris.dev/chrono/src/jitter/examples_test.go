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

package jitter_test

import (
	"fmt"
	"time"

	"arcoris.dev/chrono/delay"
	"arcoris.dev/chrono/jitter"
)

func ExampleFull() {
	schedule := jitter.Full(
		delay.Fixed(time.Second),
		jitter.WithRandomFunc(func() int64 { return 0 }),
	)

	printNext(schedule.NewSequence(), 1)

	// Output:
	// 0s
}

func ExampleEqual() {
	schedule := jitter.Equal(
		delay.Fixed(time.Second),
		jitter.WithRandomFunc(func() int64 { return int64(250 * time.Millisecond) }),
	)

	printNext(schedule.NewSequence(), 1)

	// Output:
	// 750ms
}

func ExamplePositive() {
	schedule := jitter.Positive(
		delay.Fixed(time.Second),
		0.2,
		jitter.WithRandomFunc(func() int64 { return int64(100 * time.Millisecond) }),
	)

	printNext(schedule.NewSequence(), 1)

	// Output:
	// 1.1s
}

func ExampleProportional() {
	schedule := jitter.Proportional(
		delay.Fixed(time.Second),
		0.2,
		jitter.WithRandomFunc(func() int64 { return int64(100 * time.Millisecond) }),
	)

	printNext(schedule.NewSequence(), 1)

	// Output:
	// 900ms
}

func ExampleUniform() {
	schedule := jitter.Uniform(
		100*time.Millisecond,
		500*time.Millisecond,
		jitter.WithRandomFunc(func() int64 { return int64(200 * time.Millisecond) }),
	)

	printNext(schedule.NewSequence(), 1)

	// Output:
	// 300ms
}

func ExampleDecorrelated() {
	schedule := jitter.Decorrelated(
		100*time.Millisecond,
		2*time.Second,
		3,
		jitter.WithRandomSource(sequenceSource(
			50*time.Millisecond,
			200*time.Millisecond,
			0,
		)),
	)

	printNext(schedule.NewSequence(), 3)

	// Output:
	// 150ms
	// 300ms
	// 100ms
}

func ExampleWithSeed() {
	schedule := jitter.Uniform(0, time.Second, jitter.WithSeed(42))

	left := schedule.NewSequence()
	right := schedule.NewSequence()

	leftDelay, _ := left.Next()
	rightDelay, _ := right.Next()

	fmt.Println(leftDelay == rightDelay)

	// Output:
	// true
}

func ExampleWithRandomSource() {
	source := jitter.RandomSourceFunc(func() jitter.RandomGenerator {
		return jitter.RandomFunc(func() int64 {
			return int64(250 * time.Millisecond)
		})
	})

	schedule := jitter.Uniform(0, time.Second, jitter.WithRandomSource(source))

	printNext(schedule.NewSequence(), 1)

	// Output:
	// 250ms
}

func ExampleFull_cappedExponential() {
	schedule := jitter.Full(
		delay.Cap(
			delay.Exponential(100*time.Millisecond, 2),
			2*time.Second,
		),
		jitter.WithRandomSource(sequenceSource(
			50*time.Millisecond,
			100*time.Millisecond,
			200*time.Millisecond,
		)),
	)

	printNext(schedule.NewSequence(), 3)

	// Output:
	// 50ms
	// 100ms
	// 200ms
}

func ExampleDecorrelated_limited() {
	schedule := delay.Limit(
		jitter.Decorrelated(
			100*time.Millisecond,
			2*time.Second,
			3,
			jitter.WithRandomSource(sequenceSource(
				50*time.Millisecond,
				200*time.Millisecond,
				0,
			)),
		),
		3,
	)

	printNext(schedule.NewSequence(), 4)

	// Output:
	// 150ms
	// 300ms
	// 100ms
	// exhausted
}

func printNext(seq delay.Sequence, n int) {
	for i := 0; i < n; i++ {
		d, ok := seq.Next()
		if !ok {
			fmt.Println("exhausted")
			return
		}

		fmt.Println(d)
	}
}

func sequenceSource(values ...time.Duration) jitter.RandomSource {
	return jitter.RandomSourceFunc(func() jitter.RandomGenerator {
		next := 0

		return jitter.RandomFunc(func() int64 {
			if len(values) == 0 {
				return 0
			}

			value := values[next%len(values)]
			next++

			return int64(value)
		})
	})
}

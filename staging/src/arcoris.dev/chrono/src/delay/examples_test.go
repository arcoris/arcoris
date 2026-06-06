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

package delay_test

import (
	"fmt"
	"time"

	"arcoris.dev/chrono/delay"
)

func ExampleImmediate() {
	sequence := delay.Immediate().NewSequence()

	printNext(sequence)
	printNext(sequence)

	// Output:
	// 0s true
	// 0s true
}

func ExampleFixed() {
	sequence := delay.Fixed(250 * time.Millisecond).NewSequence()

	printNext(sequence)
	printNext(sequence)
	printNext(sequence)

	// Output:
	// 250ms true
	// 250ms true
	// 250ms true
}

func ExampleDelays() {
	sequence := delay.Delays(0, 50*time.Millisecond, 200*time.Millisecond).NewSequence()

	printNext(sequence)
	printNext(sequence)
	printNext(sequence)
	printNext(sequence)

	// Output:
	// 0s true
	// 50ms true
	// 200ms true
	// 0s false
}

func ExampleLinear() {
	sequence := delay.Linear(100*time.Millisecond, 50*time.Millisecond).NewSequence()

	printNext(sequence)
	printNext(sequence)
	printNext(sequence)

	// Output:
	// 100ms true
	// 150ms true
	// 200ms true
}

func ExampleExponential() {
	sequence := delay.Exponential(100*time.Millisecond, 2).NewSequence()

	printNext(sequence)
	printNext(sequence)
	printNext(sequence)
	printNext(sequence)

	// Output:
	// 100ms true
	// 200ms true
	// 400ms true
	// 800ms true
}

func ExampleFibonacci() {
	sequence := delay.Fibonacci(100 * time.Millisecond).NewSequence()

	printNext(sequence)
	printNext(sequence)
	printNext(sequence)
	printNext(sequence)
	printNext(sequence)

	// Output:
	// 100ms true
	// 100ms true
	// 200ms true
	// 300ms true
	// 500ms true
}

func ExampleCap() {
	sequence := delay.Cap(
		delay.Delays(100*time.Millisecond, 2*time.Second),
		time.Second,
	).NewSequence()

	printNext(sequence)
	printNext(sequence)
	printNext(sequence)

	// Output:
	// 100ms true
	// 1s true
	// 0s false
}

func ExampleLimit() {
	sequence := delay.Limit(
		delay.Fixed(250*time.Millisecond),
		2,
	).NewSequence()

	printNext(sequence)
	printNext(sequence)
	printNext(sequence)

	// Output:
	// 250ms true
	// 250ms true
	// 0s false
}

func ExampleChain() {
	schedule := delay.Limit(
		delay.Cap(
			delay.Chain(
				delay.Delays(0),
				delay.Exponential(100*time.Millisecond, 2),
			),
			time.Second,
		),
		5,
	)
	sequence := schedule.NewSequence()

	printNext(sequence)
	printNext(sequence)
	printNext(sequence)
	printNext(sequence)
	printNext(sequence)
	printNext(sequence)

	// Output:
	// 0s true
	// 100ms true
	// 200ms true
	// 400ms true
	// 800ms true
	// 0s false
}

func ExampleScheduleFunc() {
	schedule := delay.ScheduleFunc(func() delay.Sequence {
		return delay.Delays(time.Second).NewSequence()
	})
	sequence := schedule.NewSequence()

	printNext(sequence)
	printNext(sequence)

	// Output:
	// 1s true
	// 0s false
}

func ExampleSequenceFunc() {
	delays := []time.Duration{0, 50 * time.Millisecond}
	sequence := delay.SequenceFunc(func() (time.Duration, bool) {
		if len(delays) == 0 {
			return 0, false
		}

		next := delays[0]
		delays = delays[1:]

		return next, true
	})

	printNext(sequence)
	printNext(sequence)
	printNext(sequence)

	// Output:
	// 0s true
	// 50ms true
	// 0s false
}

func printNext(sequence delay.Sequence) {
	d, ok := sequence.Next()
	fmt.Printf("%s %t\n", d, ok)
}

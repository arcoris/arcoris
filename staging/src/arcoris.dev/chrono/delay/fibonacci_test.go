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

package delay

import (
	"testing"
	"time"
)

func TestFibonacciRejectsNonPositiveBaseDelay(t *testing.T) {
	mustPanicWith(t, errNonPositiveFibonacciBaseDelay, func() {
		Fibonacci(0)
	})
	mustPanicWith(t, errNonPositiveFibonacciBaseDelay, func() {
		Fibonacci(-time.Nanosecond)
	})
}

func TestFibonacciSequenceUsesFibonacciGrowth(t *testing.T) {
	sequence := Fibonacci(time.Second).NewSequence()

	mustNext(t, sequence, time.Second)
	mustNext(t, sequence, time.Second)
	mustNext(t, sequence, 2*time.Second)
	mustNext(t, sequence, 3*time.Second)
	mustNext(t, sequence, 5*time.Second)
}

func TestFibonacciSequencesHaveIndependentState(t *testing.T) {
	schedule := Fibonacci(time.Second)

	l := schedule.NewSequence()
	r := schedule.NewSequence()

	mustNext(t, l, time.Second)
	mustNext(t, l, time.Second)
	mustNext(t, r, time.Second)
}

func TestFibonacciSequenceSaturates(t *testing.T) {
	sequence := &fibonacciSequence{previous: maxDuration, current: time.Nanosecond}

	mustNext(t, sequence, time.Nanosecond)
	mustNext(t, sequence, maxDuration)
	mustNext(t, sequence, maxDuration)
}

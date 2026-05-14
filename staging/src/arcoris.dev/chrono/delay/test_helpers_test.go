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

// mustPanicWith asserts that fn panics with exactly want.
//
// Delay constructors and adapters use stable package-local diagnostic strings
// for invalid programmer input. Tests compare exact panic payloads so accidental
// diagnostic drift is caught by the package test suite.
func mustPanicWith(t *testing.T, want any, fn func()) {
	t.Helper()

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("expected panic %v, got no panic", want)
		}
		if got != want {
			t.Fatalf("panic value = %v, want %v", got, want)
		}
	}()

	fn()
}

// mustNext asserts that seq returns want with ok=true.
func mustNext(t *testing.T, seq Sequence, want time.Duration) {
	t.Helper()

	got, ok := seq.Next()
	if !ok {
		t.Fatalf("Next() exhausted, want %s", want)
	}
	if got != want {
		t.Fatalf("Next() delay = %s, want %s", got, want)
	}
}

// mustExhausted asserts that seq reports finite exhaustion.
func mustExhausted(t *testing.T, seq Sequence) {
	t.Helper()

	got, ok := seq.Next()
	if ok {
		t.Fatalf("Next() = %s, true; want exhausted", got)
	}
}

// countingSequence returns configured values in order and records call count.
type countingSequence struct {
	// values is the remaining delay list returned by Next.
	values []time.Duration

	// calls is incremented for every Next invocation, including exhaustion.
	calls int
}

// Next returns the next configured delay and records that the sequence was used.
func (s *countingSequence) Next() (time.Duration, bool) {
	s.calls++
	if len(s.values) == 0 {
		return 0, false
	}
	v := s.values[0]
	s.values = s.values[1:]
	return v, true
}

// negativeDelaySequence violates the Sequence contract for wrapper tests.
type negativeDelaySequence struct{}

// Next returns a negative delay with ok=true.
func (negativeDelaySequence) Next() (time.Duration, bool) {
	return -time.Nanosecond, true
}

// nilSequenceSchedule violates the Schedule contract for wrapper tests.
type nilSequenceSchedule struct{}

// NewSequence returns nil.
func (nilSequenceSchedule) NewSequence() Sequence {
	return nil
}

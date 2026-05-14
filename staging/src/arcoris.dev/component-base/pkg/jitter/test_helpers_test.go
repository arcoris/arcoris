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

	"arcoris.dev/component-base/pkg/delay"
)

// mustPanicWith asserts that fn panics with exactly want.
//
// Jitter constructor failures are programming errors and use stable
// package-local diagnostic strings. Tests must compare exact panic payloads so
// accidental diagnostic drift is caught by the package test suite.
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
func mustNext(t *testing.T, seq delay.Sequence, want time.Duration) {
	t.Helper()

	got, ok := seq.Next()
	if !ok {
		t.Fatalf("Next() exhausted, want %s", want)
	}
	if got != want {
		t.Fatalf("Next() delay = %s, want %s", got, want)
	}
}

// mustNextInRange asserts that seq returns a delay in [lo, hi].
func mustNextInRange(t *testing.T, seq delay.Sequence, lo, hi time.Duration) time.Duration {
	t.Helper()

	got, ok := seq.Next()
	if !ok {
		t.Fatalf("Next() exhausted, want delay in [%s, %s]", lo, hi)
	}
	if got < lo || got > hi {
		t.Fatalf("Next() delay = %s, want in [%s, %s]", got, lo, hi)
	}
	return got
}

// mustExhausted asserts that seq reports finite exhaustion.
func mustExhausted(t *testing.T, seq delay.Sequence) {
	t.Helper()

	got, ok := seq.Next()
	if ok {
		t.Fatalf("Next() = %s, true; want exhausted", got)
	}
}

// fixedRandom is a deterministic RandomGenerator used for boundary tests.
type fixedRandom int64

// Int63 returns the configured fixed pseudo-random value.
func (r fixedRandom) Int63() int64 {
	return int64(r)
}

// sequenceRandom returns configured values cyclically.
type sequenceRandom struct {
	values []int64
	next   int
}

// Int63 returns the next configured value and advances the cursor.
func (r *sequenceRandom) Int63() int64 {
	if len(r.values) == 0 {
		return 0
	}
	value := r.values[r.next%len(r.values)]
	r.next++
	return value
}

// countingSequence records how many times Next is called.
type countingSequence struct {
	values []time.Duration
	calls  int
}

// Next returns configured values in order and records each call.
func (s *countingSequence) Next() (time.Duration, bool) {
	s.calls++
	if len(s.values) == 0 {
		return 0, false
	}
	value := s.values[0]
	s.values = s.values[1:]
	return value, true
}

// countingRandomSource records how many generators were requested.
type countingRandomSource struct {
	calls int
}

// NewRandom returns a deterministic generator based on the call index.
func (s *countingRandomSource) NewRandom() RandomGenerator {
	value := int64(s.calls)
	s.calls++
	return fixedRandom(value)
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
func (nilSequenceSchedule) NewSequence() delay.Sequence {
	return nil
}

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

package retry

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/chrono/delay"
)

// expectPanic asserts that fn panics with want or an error matching want.
//
// Retry exposes both stable string diagnostics and sentinel-wrapping errors.
// This helper supports exact panic values and errors.Is matching without
// weakening tests that expect a non-error panic payload.
func expectPanic(t *testing.T, want any, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("function did not panic")
		}
		if !errors.Is(asError(recovered), asError(want)) && recovered != want {
			t.Fatalf("panic = %v, want %v", recovered, want)
		}
	}()

	fn()
}

// asError returns v when it is an error and nil otherwise.
func asError(v any) error {
	err, ok := v.(error)
	if !ok {
		return nil
	}
	return err
}

// retryTestAttempt returns deterministic attempt metadata for n.
func retryTestAttempt(n uint) Attempt {
	return Attempt{
		Number:    n,
		StartedAt: time.Unix(int64(n), 0),
	}
}

// retryTestStartedAt returns the shared deterministic start timestamp.
func retryTestStartedAt() time.Time {
	return time.Unix(100, 0)
}

// retryTestFinishedAt returns the shared deterministic finish timestamp.
func retryTestFinishedAt() time.Time {
	return time.Unix(101, 0)
}

// retryTestSuccessOutcome returns an expected success Outcome for n attempts.
func retryTestSuccessOutcome(n uint) Outcome {
	return Outcome{
		Attempts:   n,
		StartedAt:  retryTestStartedAt(),
		FinishedAt: retryTestFinishedAt(),
		Reason:     StopReasonSucceeded,
	}
}

// retryTestFailureOutcome returns an expected failure Outcome for n attempts.
func retryTestFailureOutcome(n uint, reason StopReason, err error) Outcome {
	return Outcome{
		Attempts:   n,
		StartedAt:  retryTestStartedAt(),
		FinishedAt: retryTestFinishedAt(),
		LastErr:    err,
		Reason:     reason,
	}
}

// retryTestSchedule returns a preconfigured sequence for runtime tests.
type retryTestSchedule struct {
	// sequence is returned directly from NewSequence.
	//
	// Tests use this helper to exercise retry boundary validation without
	// creating additional delay package fixtures.
	sequence delay.Sequence
}

// NewSequence returns the configured test sequence.
func (s retryTestSchedule) NewSequence() delay.Sequence {
	return s.sequence
}

// retryTestSequence returns configured delays in order and then exhausts.
type retryTestSequence struct {
	// delays is the remaining delay list returned by Next.
	delays []time.Duration
}

// Next returns the next configured delay or finite exhaustion.
func (s *retryTestSequence) Next() (time.Duration, bool) {
	if len(s.delays) == 0 {
		return 0, false
	}

	d := s.delays[0]
	s.delays = s.delays[1:]
	return d, true
}

// retryObserverRecorder records observer-visible events and call order.
type retryObserverRecorder struct {
	// events is the ordered list of observed retry events.
	events []Event

	// order records the recorder name for each ObserveRetry call.
	order []string

	// name is appended to order when non-empty.
	name string
}

// ObserveRetry records observer calls without changing retry behavior.
//
// Tests use it when they need to assert observer-visible ordering or terminal
// metadata while preserving the production rule that observers are synchronous
// notifications and cannot reject a retry decision.
func (r *retryObserverRecorder) ObserveRetry(_ context.Context, event Event) {
	r.events = append(r.events, event)
	if r.name != "" {
		r.order = append(r.order, r.name)
	}
}

// retryTimerSignalClock exposes the point where retry has entered a real delay.
//
// The signal lets tests cancel context after retry has created the timer, which
// exercises the "context stopped during retry-owned delay" path without sleeps
// or scheduler timing assumptions.
type retryTimerSignalClock struct {
	// Clock is the wrapped clock used for all real clock behavior.
	clock.Clock

	// once guarantees timerCreated is closed at most once.
	once sync.Once

	// timerCreated is closed after the first delegated NewTimer call.
	timerCreated chan struct{}
}

// newRetryTimerSignalClock wraps an existing Clock for a single deterministic
// delay-observation test. It does not change timer behavior; it only reports
// that NewTimer was called.
func newRetryTimerSignalClock(base clock.Clock) *retryTimerSignalClock {
	return &retryTimerSignalClock{
		Clock:        base,
		timerCreated: make(chan struct{}),
	}
}

// NewTimer delegates to the wrapped clock and closes timerCreated once.
//
// Closing after delegation means the retry goroutine is inside waitDelay's timer
// setup path before the test cancels context.
func (c *retryTimerSignalClock) NewTimer(d time.Duration) clock.Timer {
	timer := c.Clock.NewTimer(d)
	c.once.Do(func() {
		close(c.timerCreated)
	})
	return timer
}

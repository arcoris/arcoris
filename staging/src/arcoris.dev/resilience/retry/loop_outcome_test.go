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
	"testing"

	"arcoris.dev/chrono/clock"
)

func TestRetryExecutionStopBuildsOutcomeAndEmitsStopEvent(t *testing.T) {
	errBoom := errors.New("boom")
	recorder := &retryObserverRecorder{}
	execution := retryOutcomeTestExecution(recorder)

	outcome := execution.stop(context.Background(), StopReasonNonRetryable, errBoom)

	want := retryTestFailureOutcome(1, StopReasonNonRetryable, errBoom)
	if outcome != want {
		t.Fatalf("outcome = %+v, want %+v", outcome, want)
	}
	if len(recorder.events) != 1 {
		t.Fatalf("events len = %d, want 1", len(recorder.events))
	}

	event := recorder.events[0]
	if event.Kind != EventRetryStop {
		t.Fatalf("event.Kind = %s, want %s", event.Kind, EventRetryStop)
	}
	if event.Attempt != retryTestAttempt(1) {
		t.Fatalf("event.Attempt = %+v, want %+v", event.Attempt, retryTestAttempt(1))
	}
	if event.Err != errBoom {
		t.Fatalf("event.Err = %v, want %v", event.Err, errBoom)
	}
	if event.Outcome != want {
		t.Fatalf("event.Outcome = %+v, want %+v", event.Outcome, want)
	}
	if !event.IsValid() {
		t.Fatalf("stop event is invalid: %+v", event)
	}
}

func TestRetryExecutionSucceededClearsLastError(t *testing.T) {
	errBoom := errors.New("boom")
	recorder := &retryObserverRecorder{}
	execution := retryOutcomeTestExecution(recorder)
	execution.lastErr = errBoom

	execution.succeeded(context.Background())

	if len(recorder.events) != 1 {
		t.Fatalf("events len = %d, want 1", len(recorder.events))
	}
	outcome := recorder.events[0].Outcome
	if outcome.Reason != StopReasonSucceeded {
		t.Fatalf("outcome.Reason = %s, want %s", outcome.Reason, StopReasonSucceeded)
	}
	if outcome.LastErr != nil {
		t.Fatalf("outcome.LastErr = %v, want nil", outcome.LastErr)
	}
}

func TestRetryExecutionNonRetryableReturnsOriginalError(t *testing.T) {
	errBoom := errors.New("boom")
	execution := retryOutcomeTestExecution()

	err := execution.nonRetryable(context.Background(), errBoom)

	if !errors.Is(err, errBoom) {
		t.Fatalf("nonRetryable error = %v, want %v", err, errBoom)
	}
}

func TestRetryExecutionExhaustedReturnsExhaustedError(t *testing.T) {
	errBoom := errors.New("boom")
	execution := retryOutcomeTestExecution()
	execution.lastErr = errBoom

	err := execution.exhausted(context.Background(), StopReasonMaxAttempts)

	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("exhausted error = %v, want ErrExhausted", err)
	}
	if !errors.Is(err, errBoom) {
		t.Fatalf("exhausted error = %v, want %v", err, errBoom)
	}

	outcome, ok := ExhaustedOutcome(err)
	if !ok {
		t.Fatalf("ExhaustedOutcome returned ok=false")
	}
	if outcome != retryTestFailureOutcome(1, StopReasonMaxAttempts, errBoom) {
		t.Fatalf("outcome = %+v, want %+v", outcome, retryTestFailureOutcome(1, StopReasonMaxAttempts, errBoom))
	}
}

func TestRetryExecutionInterruptedReturnsOriginalError(t *testing.T) {
	errBoom := errors.New("boom")
	ctxErr := NewInterruptedError(context.Canceled)
	recorder := &retryObserverRecorder{}
	execution := retryOutcomeTestExecution(recorder)
	execution.lastErr = errBoom

	err := execution.interrupted(context.Background(), ctxErr)

	if err != ctxErr {
		t.Fatalf("interrupted error = %v, want %v", err, ctxErr)
	}
	if len(recorder.events) != 1 {
		t.Fatalf("events len = %d, want 1", len(recorder.events))
	}
	if recorder.events[0].Outcome.LastErr != errBoom {
		t.Fatalf("outcome.LastErr = %v, want %v", recorder.events[0].Outcome.LastErr, errBoom)
	}
	if recorder.events[0].Outcome.Reason != StopReasonInterrupted {
		t.Fatalf("outcome.Reason = %s, want %s", recorder.events[0].Outcome.Reason, StopReasonInterrupted)
	}
}

func TestStopEventOmitsAttemptBeforeFirstAttempt(t *testing.T) {
	outcome := Outcome{
		StartedAt:  retryTestStartedAt(),
		FinishedAt: retryTestFinishedAt(),
		Reason:     StopReasonInterrupted,
	}

	event := stopEvent(outcome, Attempt{})

	if event.Kind != EventRetryStop {
		t.Fatalf("event.Kind = %s, want %s", event.Kind, EventRetryStop)
	}
	if event.Attempt != (Attempt{}) {
		t.Fatalf("event.Attempt = %+v, want zero attempt", event.Attempt)
	}
	if event.Err != nil {
		t.Fatalf("event.Err = %v, want nil", event.Err)
	}
	if event.Outcome != outcome {
		t.Fatalf("event.Outcome = %+v, want %+v", event.Outcome, outcome)
	}
	if !event.IsValid() {
		t.Fatalf("stop event is invalid: %+v", event)
	}
}

func TestStopEventIncludesAttemptAfterAttempt(t *testing.T) {
	errBoom := errors.New("boom")
	outcome := retryTestFailureOutcome(1, StopReasonNonRetryable, errBoom)
	attempt := retryTestAttempt(1)

	event := stopEvent(outcome, attempt)

	if event.Kind != EventRetryStop {
		t.Fatalf("event.Kind = %s, want %s", event.Kind, EventRetryStop)
	}
	if event.Attempt != attempt {
		t.Fatalf("event.Attempt = %+v, want %+v", event.Attempt, attempt)
	}
	if event.Err != errBoom {
		t.Fatalf("event.Err = %v, want %v", event.Err, errBoom)
	}
	if event.Outcome != outcome {
		t.Fatalf("event.Outcome = %+v, want %+v", event.Outcome, outcome)
	}
	if !event.IsValid() {
		t.Fatalf("stop event is invalid: %+v", event)
	}
}

func retryOutcomeTestExecution(observers ...Observer) *retryExecution {
	config := configOf(WithClock(clock.NewFakeClock(retryTestFinishedAt())))
	config.observers = append(config.observers, observers...)

	return &retryExecution{
		config:      config,
		startedAt:   retryTestStartedAt(),
		attempts:    1,
		lastAttempt: retryTestAttempt(1),
	}
}

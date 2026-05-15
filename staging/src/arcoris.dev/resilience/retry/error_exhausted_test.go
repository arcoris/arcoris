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
	"errors"
	"testing"
	"time"
)

func TestErrExhaustedSentinel(t *testing.T) {
	if ErrExhausted == nil {
		t.Fatalf("ErrExhausted is nil")
	}
	if ErrExhausted.Error() != errExhaustedMessage {
		t.Fatalf("ErrExhausted.Error() = %q, want %q", ErrExhausted.Error(), errExhaustedMessage)
	}
	if !errors.Is(ErrExhausted, ErrExhausted) {
		t.Fatalf("ErrExhausted does not match itself")
	}
}

func TestExhausted(t *testing.T) {
	errBoom := errors.New("boom")
	err := NewExhaustedError(retryTestExhaustedOutcome(StopReasonMaxAttempts, errBoom))

	if !Exhausted(err) {
		t.Fatalf("Exhausted(exhausted error) = false, want true")
	}
	if !Exhausted(ErrExhausted) {
		t.Fatalf("Exhausted(ErrExhausted) = false, want true")
	}
	if Exhausted(errBoom) {
		t.Fatalf("Exhausted(non-exhausted error) = true, want false")
	}
	if Exhausted(nil) {
		t.Fatalf("Exhausted(nil) = true, want false")
	}
}

func TestNewExhaustedError(t *testing.T) {
	errBoom := errors.New("boom")
	outcome := retryTestExhaustedOutcome(StopReasonMaxAttempts, errBoom)

	err := NewExhaustedError(outcome)
	if err == nil {
		t.Fatalf("NewExhaustedError returned nil")
	}
	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("NewExhaustedError does not match ErrExhausted")
	}
	if !errors.Is(err, errBoom) {
		t.Fatalf("NewExhaustedError does not unwrap to last error")
	}

	wantMessage := errExhaustedMessage + ": " + errBoom.Error()
	if err.Error() != wantMessage {
		t.Fatalf("error message = %q, want %q", err.Error(), wantMessage)
	}

	got, ok := ExhaustedOutcome(err)
	if !ok {
		t.Fatalf("ExhaustedOutcome returned ok=false")
	}
	if got != outcome {
		t.Fatalf("ExhaustedOutcome = %+v, want %+v", got, outcome)
	}
}

func TestNewExhaustedErrorReasons(t *testing.T) {
	errBoom := errors.New("boom")

	reasons := []StopReason{
		StopReasonMaxAttempts,
		StopReasonMaxElapsed,
		StopReasonDeadline,
		StopReasonDelayExhausted,
	}

	for _, reason := range reasons {
		t.Run(reason.String(), func(t *testing.T) {
			err := NewExhaustedError(retryTestExhaustedOutcome(reason, errBoom))
			if !errors.Is(err, ErrExhausted) {
				t.Fatalf("NewExhaustedError(%s) does not match ErrExhausted", reason)
			}

			outcome, ok := ExhaustedOutcome(err)
			if !ok {
				t.Fatalf("ExhaustedOutcome returned ok=false")
			}
			if outcome.Reason != reason {
				t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, reason)
			}
		})
	}
}

func TestNewExhaustedErrorPanicsOnInvalidOutcome(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("NewExhaustedError did not panic")
		}
		if recovered != panicInvalidExhaustedOutcome {
			t.Fatalf("panic = %v, want %q", recovered, panicInvalidExhaustedOutcome)
		}
	}()

	_ = NewExhaustedError(Outcome{})
}

func TestNewExhaustedErrorPanicsOnNonExhaustedReason(t *testing.T) {
	outcome := Outcome{
		Attempts:   1,
		StartedAt:  time.Unix(1, 0),
		FinishedAt: time.Unix(2, 0),
		Reason:     StopReasonSucceeded,
	}

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("NewExhaustedError did not panic")
		}
		if recovered != panicNonExhaustedOutcomeReason {
			t.Fatalf("panic = %v, want %q", recovered, panicNonExhaustedOutcomeReason)
		}
	}()

	_ = NewExhaustedError(outcome)
}

func TestExhaustedOutcome(t *testing.T) {
	errBoom := errors.New("boom")
	outcome := retryTestExhaustedOutcome(StopReasonMaxAttempts, errBoom)
	err := NewExhaustedError(outcome)

	got, ok := ExhaustedOutcome(err)
	if !ok {
		t.Fatalf("ExhaustedOutcome(exhausted error) ok=false, want true")
	}
	if got != outcome {
		t.Fatalf("ExhaustedOutcome = %+v, want %+v", got, outcome)
	}

	wrapped := errors.Join(errors.New("outer"), err)
	got, ok = ExhaustedOutcome(wrapped)
	if !ok {
		t.Fatalf("ExhaustedOutcome(joined exhausted error) ok=false, want true")
	}
	if got != outcome {
		t.Fatalf("ExhaustedOutcome(joined) = %+v, want %+v", got, outcome)
	}

	got, ok = ExhaustedOutcome(ErrExhausted)
	if ok {
		t.Fatalf("ExhaustedOutcome(ErrExhausted) ok=true, want false")
	}
	if !got.IsZero() {
		t.Fatalf("ExhaustedOutcome(ErrExhausted) outcome = %+v, want zero", got)
	}

	got, ok = ExhaustedOutcome(nil)
	if ok {
		t.Fatalf("ExhaustedOutcome(nil) ok=true, want false")
	}
	if !got.IsZero() {
		t.Fatalf("ExhaustedOutcome(nil) outcome = %+v, want zero", got)
	}
}

func retryTestExhaustedOutcome(reason StopReason, err error) Outcome {
	return Outcome{
		Attempts:   2,
		StartedAt:  time.Unix(1, 0),
		FinishedAt: time.Unix(2, 0),
		LastErr:    err,
		Reason:     reason,
	}
}

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

func TestEventZeroValueIsInvalid(t *testing.T) {
	var event Event

	if !event.IsZero() {
		t.Fatalf("zero Event IsZero() = false, want true")
	}
	if event.IsValid() {
		t.Fatalf("zero Event IsValid() = true, want false")
	}
}

func TestEventIsValidAttemptStart(t *testing.T) {
	attempt := retryTestAttempt(1)
	errBoom := errors.New("boom")
	outcome := retryTestSuccessOutcome(1)

	tests := []struct {
		name  string
		event Event
		want  bool
	}{
		{
			name: "valid",
			event: Event{
				Kind:    EventAttemptStart,
				Attempt: attempt,
			},
			want: true,
		},
		{
			name: "invalid attempt",
			event: Event{
				Kind: EventAttemptStart,
			},
			want: false,
		},
		{
			name: "delay set",
			event: Event{
				Kind:    EventAttemptStart,
				Attempt: attempt,
				Delay:   time.Millisecond,
			},
			want: false,
		},
		{
			name: "error set",
			event: Event{
				Kind:    EventAttemptStart,
				Attempt: attempt,
				Err:     errBoom,
			},
			want: false,
		},
		{
			name: "outcome set",
			event: Event{
				Kind:    EventAttemptStart,
				Attempt: attempt,
				Outcome: outcome,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.event.IsValid(); got != tt.want {
				t.Fatalf("Event.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventIsValidAttemptFailure(t *testing.T) {
	attempt := retryTestAttempt(1)
	errBoom := errors.New("boom")
	outcome := retryTestSuccessOutcome(1)

	tests := []struct {
		name  string
		event Event
		want  bool
	}{
		{
			name: "valid",
			event: Event{
				Kind:    EventAttemptFailure,
				Attempt: attempt,
				Err:     errBoom,
			},
			want: true,
		},
		{
			name: "missing error",
			event: Event{
				Kind:    EventAttemptFailure,
				Attempt: attempt,
			},
			want: false,
		},
		{
			name: "invalid attempt",
			event: Event{
				Kind: EventAttemptFailure,
				Err:  errBoom,
			},
			want: false,
		},
		{
			name: "delay set",
			event: Event{
				Kind:    EventAttemptFailure,
				Attempt: attempt,
				Delay:   time.Millisecond,
				Err:     errBoom,
			},
			want: false,
		},
		{
			name: "outcome set",
			event: Event{
				Kind:    EventAttemptFailure,
				Attempt: attempt,
				Err:     errBoom,
				Outcome: outcome,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.event.IsValid(); got != tt.want {
				t.Fatalf("Event.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventIsValidRetryDelay(t *testing.T) {
	attempt := retryTestAttempt(1)
	errBoom := errors.New("boom")
	outcome := retryTestSuccessOutcome(1)

	tests := []struct {
		name  string
		event Event
		want  bool
	}{
		{
			name: "valid positive delay",
			event: Event{
				Kind:    EventRetryDelay,
				Attempt: attempt,
				Delay:   time.Second,
				Err:     errBoom,
			},
			want: true,
		},
		{
			name: "valid zero delay",
			event: Event{
				Kind:    EventRetryDelay,
				Attempt: attempt,
				Err:     errBoom,
			},
			want: true,
		},
		{
			name: "negative delay",
			event: Event{
				Kind:    EventRetryDelay,
				Attempt: attempt,
				Delay:   -time.Nanosecond,
				Err:     errBoom,
			},
			want: false,
		},
		{
			name: "missing error",
			event: Event{
				Kind:    EventRetryDelay,
				Attempt: attempt,
				Delay:   time.Second,
			},
			want: false,
		},
		{
			name: "invalid attempt",
			event: Event{
				Kind:  EventRetryDelay,
				Delay: time.Second,
				Err:   errBoom,
			},
			want: false,
		},
		{
			name: "outcome set",
			event: Event{
				Kind:    EventRetryDelay,
				Attempt: attempt,
				Delay:   time.Second,
				Err:     errBoom,
				Outcome: outcome,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.event.IsValid(); got != tt.want {
				t.Fatalf("Event.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventIsValidRetryStop(t *testing.T) {
	errBoom := errors.New("boom")

	successOutcome := retryTestSuccessOutcome(2)
	nonRetryableOutcome := retryTestFailureOutcome(1, StopReasonNonRetryable, errBoom)
	interruptedBeforeAttempt := Outcome{
		StartedAt:  retryTestStartedAt(),
		FinishedAt: retryTestFinishedAt(),
		Reason:     StopReasonInterrupted,
	}

	tests := []struct {
		name  string
		event Event
		want  bool
	}{
		{
			name: "valid success",
			event: Event{
				Kind:    EventRetryStop,
				Attempt: retryTestAttempt(2),
				Outcome: successOutcome,
			},
			want: true,
		},
		{
			name: "valid non retryable failure",
			event: Event{
				Kind:    EventRetryStop,
				Attempt: retryTestAttempt(1),
				Err:     errBoom,
				Outcome: nonRetryableOutcome,
			},
			want: true,
		},
		{
			name: "valid interrupted before first attempt",
			event: Event{
				Kind:    EventRetryStop,
				Outcome: interruptedBeforeAttempt,
			},
			want: true,
		},
		{
			name: "invalid outcome",
			event: Event{
				Kind: EventRetryStop,
			},
			want: false,
		},
		{
			name: "delay set",
			event: Event{
				Kind:    EventRetryStop,
				Attempt: retryTestAttempt(2),
				Delay:   time.Millisecond,
				Outcome: successOutcome,
			},
			want: false,
		},
		{
			name: "missing attempt after operation calls",
			event: Event{
				Kind:    EventRetryStop,
				Outcome: successOutcome,
			},
			want: false,
		},
		{
			name: "attempt number does not match outcome attempts",
			event: Event{
				Kind:    EventRetryStop,
				Attempt: retryTestAttempt(1),
				Outcome: successOutcome,
			},
			want: false,
		},
		{
			name: "success with event error",
			event: Event{
				Kind:    EventRetryStop,
				Attempt: retryTestAttempt(2),
				Err:     errBoom,
				Outcome: successOutcome,
			},
			want: false,
		},
		{
			name: "failure without event error",
			event: Event{
				Kind:    EventRetryStop,
				Attempt: retryTestAttempt(1),
				Outcome: nonRetryableOutcome,
			},
			want: false,
		},
		{
			name: "interrupted before first attempt with attempt",
			event: Event{
				Kind:    EventRetryStop,
				Attempt: retryTestAttempt(1),
				Outcome: interruptedBeforeAttempt,
			},
			want: false,
		},
		{
			name: "interrupted before first attempt with error",
			event: Event{
				Kind:    EventRetryStop,
				Err:     errBoom,
				Outcome: interruptedBeforeAttempt,
			},
			want: false,
		},
		{
			name: "non-stop event with stop payload",
			event: Event{
				Kind:    EventAttemptStart,
				Attempt: retryTestAttempt(1),
				Outcome: successOutcome,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.event.IsValid(); got != tt.want {
				t.Fatalf("Event.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventIsZero(t *testing.T) {
	errBoom := errors.New("boom")

	tests := []struct {
		name  string
		event Event
		want  bool
	}{
		{
			name:  "zero",
			event: Event{},
			want:  true,
		},
		{
			name: "kind set",
			event: Event{
				Kind: EventAttemptStart,
			},
			want: false,
		},
		{
			name: "attempt set",
			event: Event{
				Attempt: retryTestAttempt(1),
			},
			want: false,
		},
		{
			name: "delay set",
			event: Event{
				Delay: time.Nanosecond,
			},
			want: false,
		},
		{
			name: "error set",
			event: Event{
				Err: errBoom,
			},
			want: false,
		},
		{
			name: "outcome set",
			event: Event{
				Outcome: retryTestSuccessOutcome(1),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.event.IsZero(); got != tt.want {
				t.Fatalf("Event.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func retryTestAttempt(number uint) Attempt {
	return Attempt{
		Number:    number,
		StartedAt: time.Unix(int64(number), 0),
	}
}

func retryTestStartedAt() time.Time {
	return time.Unix(100, 0)
}

func retryTestFinishedAt() time.Time {
	return time.Unix(101, 0)
}

func retryTestSuccessOutcome(attempts uint) Outcome {
	return Outcome{
		Attempts:   attempts,
		StartedAt:  retryTestStartedAt(),
		FinishedAt: retryTestFinishedAt(),
		Reason:     StopReasonSucceeded,
	}
}

func retryTestFailureOutcome(attempts uint, reason StopReason, err error) Outcome {
	return Outcome{
		Attempts:   attempts,
		StartedAt:  retryTestStartedAt(),
		FinishedAt: retryTestFinishedAt(),
		LastErr:    err,
		Reason:     reason,
	}
}

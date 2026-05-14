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

func TestOutcomeZeroValueIsInvalid(t *testing.T) {
	var outcome Outcome

	if !outcome.IsZero() {
		t.Fatalf("zero Outcome IsZero() = false, want true")
	}
	if outcome.IsValid() {
		t.Fatalf("zero Outcome IsValid() = true, want false")
	}
	if outcome.Succeeded() {
		t.Fatalf("zero Outcome Succeeded() = true, want false")
	}
	if outcome.Failed() {
		t.Fatalf("zero Outcome Failed() = true, want false")
	}
	if outcome.Exhausted() {
		t.Fatalf("zero Outcome Exhausted() = true, want false")
	}
	if outcome.Interrupted() {
		t.Fatalf("zero Outcome Interrupted() = true, want false")
	}
	if got := outcome.Duration(); got != 0 {
		t.Fatalf("zero Outcome Duration() = %s, want 0", got)
	}
}

func TestOutcomeIsValid(t *testing.T) {
	started := time.Unix(10, 0)
	finished := time.Unix(12, 0)
	errBoom := errors.New("boom")

	tests := []struct {
		name    string
		outcome Outcome
		want    bool
	}{
		{
			name: "succeeded",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonSucceeded,
			},
			want: true,
		},
		{
			name: "non retryable",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonNonRetryable,
			},
			want: true,
		},
		{
			name: "max attempts",
			outcome: Outcome{
				Attempts:   3,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonMaxAttempts,
			},
			want: true,
		},
		{
			name: "max elapsed",
			outcome: Outcome{
				Attempts:   2,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonMaxElapsed,
			},
			want: true,
		},
		{
			name: "delay exhausted",
			outcome: Outcome{
				Attempts:   2,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonDelayExhausted,
			},
			want: true,
		},
		{
			name: "interrupted before first attempt",
			outcome: Outcome{
				Attempts:   0,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonInterrupted,
			},
			want: true,
		},
		{
			name: "interrupted after failed attempt",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonInterrupted,
			},
			want: true,
		},
		{
			name: "invalid reason",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     0,
			},
			want: false,
		},
		{
			name: "zero started at",
			outcome: Outcome{
				Attempts:   1,
				FinishedAt: finished,
				Reason:     StopReasonSucceeded,
			},
			want: false,
		},
		{
			name: "zero finished at",
			outcome: Outcome{
				Attempts:  1,
				StartedAt: started,
				Reason:    StopReasonSucceeded,
			},
			want: false,
		},
		{
			name: "finished before started",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  finished,
				FinishedAt: started,
				Reason:     StopReasonSucceeded,
			},
			want: false,
		},
		{
			name: "success with last error",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonSucceeded,
			},
			want: false,
		},
		{
			name: "success with zero attempts",
			outcome: Outcome{
				Attempts:   0,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonSucceeded,
			},
			want: false,
		},
		{
			name: "non retryable without last error",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonNonRetryable,
			},
			want: false,
		},
		{
			name: "max attempts without last error",
			outcome: Outcome{
				Attempts:   3,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonMaxAttempts,
			},
			want: false,
		},
		{
			name: "max elapsed without last error",
			outcome: Outcome{
				Attempts:   2,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonMaxElapsed,
			},
			want: false,
		},
		{
			name: "delay exhausted without last error",
			outcome: Outcome{
				Attempts:   2,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonDelayExhausted,
			},
			want: false,
		},
		{
			name: "non interrupted failure with zero attempts",
			outcome: Outcome{
				Attempts:   0,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonMaxAttempts,
			},
			want: false,
		},
		{
			name: "interrupted before first attempt with last error",
			outcome: Outcome{
				Attempts:   0,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonInterrupted,
			},
			want: false,
		},
		{
			name: "interrupted after attempt without last error",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonInterrupted,
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.outcome.IsValid(); got != tc.want {
				t.Fatalf("Outcome.IsValid() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestOutcomePredicatesRequireValidOutcome(t *testing.T) {
	started := time.Unix(10, 0)
	finished := time.Unix(12, 0)
	errBoom := errors.New("boom")

	tests := []struct {
		name            string
		outcome         Outcome
		wantSucceeded   bool
		wantFailed      bool
		wantExhausted   bool
		wantInterrupted bool
	}{
		{
			name: "succeeded",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonSucceeded,
			},
			wantSucceeded: true,
		},
		{
			name: "non retryable",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonNonRetryable,
			},
			wantFailed: true,
		},
		{
			name: "max attempts",
			outcome: Outcome{
				Attempts:   3,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonMaxAttempts,
			},
			wantFailed:    true,
			wantExhausted: true,
		},
		{
			name: "max elapsed",
			outcome: Outcome{
				Attempts:   2,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonMaxElapsed,
			},
			wantFailed:    true,
			wantExhausted: true,
		},
		{
			name: "delay exhausted",
			outcome: Outcome{
				Attempts:   2,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonDelayExhausted,
			},
			wantFailed:    true,
			wantExhausted: true,
		},
		{
			name: "interrupted before first attempt",
			outcome: Outcome{
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonInterrupted,
			},
			wantFailed:      true,
			wantInterrupted: true,
		},
		{
			name: "interrupted after failed attempt",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonInterrupted,
			},
			wantFailed:      true,
			wantInterrupted: true,
		},
		{
			name: "invalid success with last error",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				LastErr:    errBoom,
				Reason:     StopReasonSucceeded,
			},
		},
		{
			name: "invalid exhausted without last error",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     StopReasonMaxAttempts,
			},
		},
		{
			name: "invalid reason",
			outcome: Outcome{
				Attempts:   1,
				StartedAt:  started,
				FinishedAt: finished,
				Reason:     0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.outcome.Succeeded(); got != tc.wantSucceeded {
				t.Fatalf("Outcome.Succeeded() = %v, want %v", got, tc.wantSucceeded)
			}
			if got := tc.outcome.Failed(); got != tc.wantFailed {
				t.Fatalf("Outcome.Failed() = %v, want %v", got, tc.wantFailed)
			}
			if got := tc.outcome.Exhausted(); got != tc.wantExhausted {
				t.Fatalf("Outcome.Exhausted() = %v, want %v", got, tc.wantExhausted)
			}
			if got := tc.outcome.Interrupted(); got != tc.wantInterrupted {
				t.Fatalf("Outcome.Interrupted() = %v, want %v", got, tc.wantInterrupted)
			}
		})
	}
}

func TestOutcomeDuration(t *testing.T) {
	started := time.Unix(10, 0)
	finished := time.Unix(12, int64(500*time.Millisecond))

	tests := []struct {
		name    string
		outcome Outcome
		want    time.Duration
	}{
		{
			name: "valid timestamps",
			outcome: Outcome{
				StartedAt:  started,
				FinishedAt: finished,
			},
			want: 2500 * time.Millisecond,
		},
		{
			name: "same timestamp",
			outcome: Outcome{
				StartedAt:  started,
				FinishedAt: started,
			},
			want: 0,
		},
		{
			name: "zero started at",
			outcome: Outcome{
				FinishedAt: finished,
			},
			want: 0,
		},
		{
			name: "zero finished at",
			outcome: Outcome{
				StartedAt: started,
			},
			want: 0,
		},
		{
			name: "finished before started",
			outcome: Outcome{
				StartedAt:  finished,
				FinishedAt: started,
			},
			want: 0,
		},
		{
			name:    "zero outcome",
			outcome: Outcome{},
			want:    0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.outcome.Duration(); got != tc.want {
				t.Fatalf("Outcome.Duration() = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestOutcomeIsZero(t *testing.T) {
	started := time.Unix(10, 0)
	errBoom := errors.New("boom")

	tests := []struct {
		name    string
		outcome Outcome
		want    bool
	}{
		{
			name:    "zero",
			outcome: Outcome{},
			want:    true,
		},
		{
			name: "attempts set",
			outcome: Outcome{
				Attempts: 1,
			},
			want: false,
		},
		{
			name: "started at set",
			outcome: Outcome{
				StartedAt: started,
			},
			want: false,
		},
		{
			name: "finished at set",
			outcome: Outcome{
				FinishedAt: started,
			},
			want: false,
		},
		{
			name: "last error set",
			outcome: Outcome{
				LastErr: errBoom,
			},
			want: false,
		},
		{
			name: "reason set",
			outcome: Outcome{
				Reason: StopReasonSucceeded,
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.outcome.IsZero(); got != tc.want {
				t.Fatalf("Outcome.IsZero() = %v, want %v", got, tc.want)
			}
		})
	}
}

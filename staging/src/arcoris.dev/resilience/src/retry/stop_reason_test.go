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

import "testing"

func TestStopReasonStringAndPredicates(t *testing.T) {
	tests := []struct {
		name            string
		reason          StopReason
		wantString      string
		wantValid       bool
		wantSucceeded   bool
		wantFailed      bool
		wantExhausted   bool
		wantInterrupted bool
	}{
		{
			name:          "succeeded",
			reason:        StopReasonSucceeded,
			wantString:    "succeeded",
			wantValid:     true,
			wantSucceeded: true,
		},
		{
			name:       "non retryable",
			reason:     StopReasonNonRetryable,
			wantString: "non_retryable",
			wantValid:  true,
			wantFailed: true,
		},
		{
			name:          "max attempts",
			reason:        StopReasonMaxAttempts,
			wantString:    "max_attempts",
			wantValid:     true,
			wantFailed:    true,
			wantExhausted: true,
		},
		{
			name:          "max elapsed",
			reason:        StopReasonMaxElapsed,
			wantString:    "max_elapsed",
			wantValid:     true,
			wantFailed:    true,
			wantExhausted: true,
		},
		{
			name:          "deadline",
			reason:        StopReasonDeadline,
			wantString:    "deadline",
			wantValid:     true,
			wantFailed:    true,
			wantExhausted: true,
		},
		{
			name:          "delay exhausted",
			reason:        StopReasonDelayExhausted,
			wantString:    "delay_exhausted",
			wantValid:     true,
			wantFailed:    true,
			wantExhausted: true,
		},
		{
			name:            "interrupted",
			reason:          StopReasonInterrupted,
			wantString:      "interrupted",
			wantValid:       true,
			wantFailed:      true,
			wantInterrupted: true,
		},
		{
			name:       "invalid",
			reason:     0,
			wantString: "invalid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.reason.String(); got != tc.wantString {
				t.Fatalf("String() = %q, want %q", got, tc.wantString)
			}
			if got := tc.reason.IsValid(); got != tc.wantValid {
				t.Fatalf("IsValid() = %v, want %v", got, tc.wantValid)
			}
			if got := tc.reason.Succeeded(); got != tc.wantSucceeded {
				t.Fatalf("Succeeded() = %v, want %v", got, tc.wantSucceeded)
			}
			if got := tc.reason.Failed(); got != tc.wantFailed {
				t.Fatalf("Failed() = %v, want %v", got, tc.wantFailed)
			}
			if got := tc.reason.Exhausted(); got != tc.wantExhausted {
				t.Fatalf("Exhausted() = %v, want %v", got, tc.wantExhausted)
			}
			if got := tc.reason.Interrupted(); got != tc.wantInterrupted {
				t.Fatalf("Interrupted() = %v, want %v", got, tc.wantInterrupted)
			}
		})
	}
}

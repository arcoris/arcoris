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

package health

import (
	"strings"
	"testing"
)

func TestReasonString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		reason Reason
		want   string
	}{
		{ReasonNone, "none"},
		{ReasonTimeout, "timeout"},
		{Reason("custom_reason"), "custom_reason"},
		{Reason("Bad"), "invalid"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.want, func(t *testing.T) {
			t.Parallel()

			if got := test.reason.String(); got != test.want {
				t.Fatalf("String() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestReasonValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		reason Reason
		want   bool
	}{
		{ReasonNone, true},
		{Reason("custom_reason_1"), true},
		{Reason("1custom"), false},
		{Reason("_custom"), false},
		{Reason("custom_"), false},
		{Reason("custom__reason"), false},
		{Reason("custom-reason"), false},
		{Reason(strings.Repeat("a", maxReasonLength+1)), false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.reason.String(), func(t *testing.T) {
			t.Parallel()

			if got := test.reason.IsValid(); got != test.want {
				t.Fatalf("IsValid() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestReasonClassification(t *testing.T) {
	t.Parallel()

	if !ReasonNone.IsNone() {
		t.Fatal("ReasonNone.IsNone() = false, want true")
	}
	if !ReasonTimeout.IsBuiltin() || Reason("custom_reason").IsBuiltin() {
		t.Fatal("builtin classification mismatch")
	}
	if !ReasonTimeout.IsExecutionReason() || ReasonFatal.IsExecutionReason() {
		t.Fatal("execution classification mismatch")
	}
	if !ReasonShuttingDown.IsLifecycleReason() || ReasonOverloaded.IsLifecycleReason() {
		t.Fatal("lifecycle classification mismatch")
	}
	if !ReasonAdmissionClosed.IsControlReason() || ReasonTimeout.IsControlReason() {
		t.Fatal("control classification mismatch")
	}
}

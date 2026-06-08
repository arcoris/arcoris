// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deadlineadmission

import (
	"testing"
	"time"

	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
	"arcoris.dev/resilience/deadline"
)

func TestAdmissionResultMapsValidDecisions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		decision deadline.Decision
		want     admission.Decision
		admitted bool
		denied   bool
		metadata deadline.Decision
	}{
		{
			name: "allowed budget",
			decision: deadline.Decision{
				Allowed:   true,
				Remaining: time.Second,
				Reason:    deadline.ReasonAllowed,
			},
			want:     admission.AdmitDecision(admission.ReasonAdmitted),
			admitted: true,
			metadata: deadline.Decision{
				Allowed:   true,
				Remaining: time.Second,
				Reason:    deadline.ReasonAllowed,
			},
		},
		{
			name: "allowed no deadline",
			decision: deadline.Decision{
				Allowed: true,
				Reason:  deadline.ReasonNoDeadline,
			},
			want:     admission.AdmitDecision(admission.ReasonAdmitted),
			admitted: true,
			metadata: deadline.Decision{
				Allowed: true,
				Reason:  deadline.ReasonNoDeadline,
			},
		},
		{
			name: "expired",
			decision: deadline.Decision{
				Reason: deadline.ReasonExpired,
			},
			want:   admission.DenyDecision(admissionbuiltin.ReasonDeadlineExceeded),
			denied: true,
			metadata: deadline.Decision{
				Reason: deadline.ReasonExpired,
			},
		},
		{
			name: "insufficient budget",
			decision: deadline.Decision{
				Remaining: time.Second,
				Reason:    deadline.ReasonInsufficientBudget,
			},
			want:   admission.DenyDecision(admissionbuiltin.ReasonDeadlineExceeded),
			denied: true,
			metadata: deadline.Decision{
				Remaining: time.Second,
				Reason:    deadline.ReasonInsufficientBudget,
			},
		},
		{
			name: "context done without deadline budget",
			decision: deadline.Decision{
				Reason: deadline.ReasonContextDone,
			},
			want:   admission.DenyDecision(admissionbuiltin.ReasonCanceled),
			denied: true,
			metadata: deadline.Decision{
				Reason: deadline.ReasonContextDone,
			},
		},
		{
			name: "context done with future deadline budget",
			decision: deadline.Decision{
				Remaining: time.Second,
				Reason:    deadline.ReasonContextDone,
			},
			want:   admission.DenyDecision(admissionbuiltin.ReasonCanceled),
			denied: true,
			metadata: deadline.Decision{
				Remaining: time.Second,
				Reason:    deadline.ReasonContextDone,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := AdmissionResult(test.decision)
			if !result.IsValid() {
				t.Fatalf("AdmissionResult is invalid: %+v", result.Decision())
			}
			if got := result.Decision(); got != test.want {
				t.Fatalf("decision = %+v, want %+v", got, test.want)
			}
			if got := result.Decision().IsAdmitted(); got != test.admitted {
				t.Fatalf("IsAdmitted = %v, want %v", got, test.admitted)
			}
			if got := result.Decision().IsDenied(); got != test.denied {
				t.Fatalf("IsDenied = %v, want %v", got, test.denied)
			}
			if result.Decision().HasSideEffect() {
				t.Fatal("AdmissionResult has side effect, want none")
			}
			if result.HasGrant() {
				t.Fatal("AdmissionResult has grant, want none")
			}
			if _, ok := result.Grant(); ok {
				t.Fatal("GrantDecision() ok=true, want false")
			}
			if !result.HasMetadata() {
				t.Fatal("AdmissionResult has no metadata")
			}
			if metadata, ok := result.Metadata(); !ok || metadata != test.metadata {
				t.Fatalf("metadata = (%+v, %t), want (%+v, true)", metadata, ok, test.metadata)
			}
		})
	}
}

func TestAdmissionResultInvalidDecisionStaysInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		decision deadline.Decision
	}{
		{
			name:     "zero",
			decision: deadline.Decision{},
		},
		{
			name: "invalid allowed",
			decision: deadline.Decision{
				Allowed:   true,
				Remaining: time.Second,
				Reason:    deadline.ReasonExpired,
			},
		},
		{
			name: "invalid denied",
			decision: deadline.Decision{
				Reason: deadline.ReasonAllowed,
			},
		},
		{
			name: "expired with remaining budget",
			decision: deadline.Decision{
				Remaining: time.Second,
				Reason:    deadline.ReasonExpired,
			},
		},
		{
			name: "insufficient budget without remaining budget",
			decision: deadline.Decision{
				Reason: deadline.ReasonInsufficientBudget,
			},
		},
		{
			name: "insufficient budget with negative remaining",
			decision: deadline.Decision{
				Remaining: -time.Nanosecond,
				Reason:    deadline.ReasonInsufficientBudget,
			},
		},
		{
			name: "unknown reason",
			decision: deadline.Decision{
				Reason: deadline.Reason(255),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := AdmissionResult(test.decision)
			if result.IsValid() {
				t.Fatalf("AdmissionResult().IsValid() = true, want false: %+v", result.Decision())
			}
			if result.HasGrant() {
				t.Fatal("invalid AdmissionResult has grant, want none")
			}
			if _, ok := result.Grant(); ok {
				t.Fatal("invalid AdmissionResult GrantDecision() ok=true, want false")
			}
			if !result.HasMetadata() {
				t.Fatal("invalid AdmissionResult should still preserve metadata")
			}
			if metadata, ok := result.Metadata(); !ok || metadata != test.decision {
				t.Fatalf("metadata = (%+v, %t), want (%+v, true)", metadata, ok, test.decision)
			}
		})
	}
}

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

package deadline

import (
	"testing"
	"time"

	"arcoris.dev/admission"
)

func TestDecisionAdmissionResultMapsValidDecisions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		decision Decision
		want     admission.Decision
		admitted bool
		denied   bool
		metadata Decision
	}{
		{
			name: "allowed budget",
			decision: Decision{
				Allowed:   true,
				Remaining: time.Second,
				Reason:    ReasonAllowed,
			},
			want:     admission.Admit(admission.ReasonAdmitted),
			admitted: true,
			metadata: Decision{
				Allowed:   true,
				Remaining: time.Second,
				Reason:    ReasonAllowed,
			},
		},
		{
			name: "allowed no deadline",
			decision: Decision{
				Allowed: true,
				Reason:  ReasonNoDeadline,
			},
			want:     admission.Admit(admission.ReasonAdmitted),
			admitted: true,
			metadata: Decision{
				Allowed: true,
				Reason:  ReasonNoDeadline,
			},
		},
		{
			name: "expired",
			decision: Decision{
				Reason: ReasonExpired,
			},
			want:   admission.Deny(admission.ReasonDeadlineExceeded),
			denied: true,
			metadata: Decision{
				Reason: ReasonExpired,
			},
		},
		{
			name: "insufficient budget",
			decision: Decision{
				Remaining: time.Second,
				Reason:    ReasonInsufficientBudget,
			},
			want:   admission.Deny(admission.ReasonDeadlineExceeded),
			denied: true,
			metadata: Decision{
				Remaining: time.Second,
				Reason:    ReasonInsufficientBudget,
			},
		},
		{
			name: "context done without deadline budget",
			decision: Decision{
				Reason: ReasonContextDone,
			},
			want:   admission.Deny(admission.ReasonCanceled),
			denied: true,
			metadata: Decision{
				Reason: ReasonContextDone,
			},
		},
		{
			name: "context done with future deadline budget",
			decision: Decision{
				Remaining: time.Second,
				Reason:    ReasonContextDone,
			},
			want:   admission.Deny(admission.ReasonCanceled),
			denied: true,
			metadata: Decision{
				Remaining: time.Second,
				Reason:    ReasonContextDone,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := test.decision.AdmissionResult()
			if !result.IsValid() {
				t.Fatalf("AdmissionResult is invalid: %+v", result.Decision())
			}
			if got := result.Decision(); got != test.want {
				t.Fatalf("decision = %+v, want %+v", got, test.want)
			}
			if got := result.IsAdmitted(); got != test.admitted {
				t.Fatalf("IsAdmitted = %v, want %v", got, test.admitted)
			}
			if got := result.IsDenied(); got != test.denied {
				t.Fatalf("IsDenied = %v, want %v", got, test.denied)
			}
			if result.HasSideEffect() {
				t.Fatal("AdmissionResult has side effect, want none")
			}
			if result.HasGrant() {
				t.Fatal("AdmissionResult has grant, want none")
			}
			if _, ok := result.Grant(); ok {
				t.Fatal("Grant() ok=true, want false")
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

func TestDecisionAdmissionResultInvalidDecisionStaysInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		decision Decision
	}{
		{
			name:     "zero",
			decision: Decision{},
		},
		{
			name: "invalid allowed",
			decision: Decision{
				Allowed:   true,
				Remaining: time.Second,
				Reason:    ReasonExpired,
			},
		},
		{
			name: "invalid denied",
			decision: Decision{
				Reason: ReasonAllowed,
			},
		},
		{
			name: "expired with remaining budget",
			decision: Decision{
				Remaining: time.Second,
				Reason:    ReasonExpired,
			},
		},
		{
			name: "insufficient budget without remaining budget",
			decision: Decision{
				Reason: ReasonInsufficientBudget,
			},
		},
		{
			name: "insufficient budget with negative remaining",
			decision: Decision{
				Remaining: -time.Nanosecond,
				Reason:    ReasonInsufficientBudget,
			},
		},
		{
			name: "unknown reason",
			decision: Decision{
				Reason: Reason(255),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := test.decision.AdmissionResult()
			if result.IsValid() {
				t.Fatalf("AdmissionResult().IsValid() = true, want false: %+v", result.Decision())
			}
			if result.HasGrant() {
				t.Fatal("invalid AdmissionResult has grant, want none")
			}
			if _, ok := result.Grant(); ok {
				t.Fatal("invalid AdmissionResult Grant() ok=true, want false")
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

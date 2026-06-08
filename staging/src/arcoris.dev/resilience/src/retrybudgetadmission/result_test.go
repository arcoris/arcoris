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

package retrybudgetadmission

import (
	"testing"

	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
	"arcoris.dev/resilience/retrybudget"
)

func TestAdmissionResultAllowedMapsToCommitted(t *testing.T) {
	t.Parallel()

	decision := validAllowedDecision()
	result := AdmissionResult(decision)

	if !result.IsValid() || !result.Decision().IsAdmitted() {
		t.Fatalf("AdmissionResult is invalid or not admitted: %+v", result.Decision())
	}
	if got, want := result.Decision(), admission.CommitDecision(admission.ReasonAdmitted); got != want {
		t.Fatalf("Decision() = %+v, want %+v", got, want)
	}
	if result.HasGrant() {
		t.Fatal("admitted retry-budget result carried a grant")
	}
	metadata, ok := result.Metadata()
	if !ok {
		t.Fatal("AdmissionResult has no metadata")
	}
	if metadata != decision {
		t.Fatalf("metadata = %+v, want original decision %+v", metadata, decision)
	}
}

func TestAdmissionResultDeniedMapsToBudgetExhausted(t *testing.T) {
	t.Parallel()

	decision := validDeniedDecision()
	result := AdmissionResult(decision)

	if !result.IsValid() || !result.Decision().IsDenied() {
		t.Fatalf("AdmissionResult is invalid or not denied: %+v", result.Decision())
	}
	if got, want := result.Decision(), admission.DenyDecision(admissionbuiltin.ReasonBudgetExhausted); got != want {
		t.Fatalf("Decision() = %+v, want %+v", got, want)
	}
	if result.HasGrant() {
		t.Fatal("denied retry-budget result carried a grant")
	}
	metadata, ok := result.Metadata()
	if !ok {
		t.Fatal("AdmissionResult has no metadata")
	}
	if metadata != decision {
		t.Fatalf("metadata = %+v, want original decision %+v", metadata, decision)
	}
}

func TestAdmissionResultInvalidDecisionStaysInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		decision retrybudget.Decision
	}{
		{
			name:     "invalid allowed decision",
			decision: retrybudget.Decision{Allowed: true, Reason: retrybudget.ReasonExhausted},
		},
		{
			name:     "invalid denied decision",
			decision: retrybudget.Decision{Allowed: false, Reason: retrybudget.ReasonAllowed},
		},
		{
			name:     "zero decision",
			decision: retrybudget.Decision{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AdmissionResult(tt.decision)
			if result.IsValid() {
				t.Fatalf("AdmissionResult().IsValid() = true, want false: %+v", result.Decision())
			}
			metadata, ok := result.Metadata()
			if !ok {
				t.Fatal("AdmissionResult has no metadata")
			}
			if metadata != tt.decision {
				t.Fatalf("metadata = %+v, want original decision %+v", metadata, tt.decision)
			}
		})
	}
}

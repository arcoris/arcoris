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

package retrybudget

import (
	"testing"

	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
	"arcoris.dev/snapshot"
)

func TestDecisionAdmissionResultAllowedMapsToCommitted(t *testing.T) {
	snap := validGenericSnapshot()
	decision := Decision{
		Allowed:  true,
		Reason:   ReasonAllowed,
		Snapshot: snap,
	}

	result := decision.AdmissionResult()
	if !result.IsValid() {
		t.Fatalf("AdmissionResult is invalid: %+v", result.Decision())
	}
	if !result.Decision().IsAdmitted() {
		t.Fatal("AdmissionResult is not admitted")
	}
	if result.Decision().IsDenied() {
		t.Fatal("AdmissionResult is denied, want admitted")
	}
	if !result.Decision().HasSideEffect() {
		t.Fatal("AdmissionResult has no side effect")
	}
	if result.HasGrant() {
		t.Fatal("AdmissionResult has grant, want none")
	}
	if _, ok := result.Grant(); ok {
		t.Fatal("AdmissionResult GrantDecision() ok=true, want false")
	}
	if !result.HasMetadata() {
		t.Fatal("AdmissionResult has no metadata")
	}
	if got, want := result.Decision(), admission.CommitDecision(admission.ReasonAdmitted); got != want {
		t.Fatalf("decision = %+v, want %+v", got, want)
	}
	if metadata, ok := result.Metadata(); !ok || metadata != snap {
		t.Fatalf("metadata = (%+v, %t), want (%+v, true)", metadata, ok, snap)
	}
}

func TestDecisionAdmissionResultDeniedMapsToBudgetExhausted(t *testing.T) {
	snap := snapshot.Snapshot[Snapshot]{
		Revision: snapshot.ZeroRevision.Next(),
		Value:    exhaustedSnapshotValue(),
	}
	decision := Decision{
		Allowed:  false,
		Reason:   ReasonExhausted,
		Snapshot: snap,
	}

	result := decision.AdmissionResult()
	if !result.IsValid() {
		t.Fatalf("AdmissionResult is invalid: %+v", result.Decision())
	}
	if !result.Decision().IsDenied() {
		t.Fatal("AdmissionResult is not denied")
	}
	if result.Decision().IsAdmitted() {
		t.Fatal("AdmissionResult is admitted, want denied")
	}
	if result.Decision().HasSideEffect() {
		t.Fatal("AdmissionResult has side effect")
	}
	if result.HasGrant() {
		t.Fatal("AdmissionResult has grant, want none")
	}
	if _, ok := result.Grant(); ok {
		t.Fatal("AdmissionResult GrantDecision() ok=true, want false")
	}
	if !result.HasMetadata() {
		t.Fatal("AdmissionResult has no metadata")
	}
	if got, want := result.Decision(), admission.DenyDecision(admissionbuiltin.ReasonBudgetExhausted); got != want {
		t.Fatalf("decision = %+v, want %+v", got, want)
	}
	if metadata, ok := result.Metadata(); !ok || metadata != snap {
		t.Fatalf("metadata = (%+v, %t), want (%+v, true)", metadata, ok, snap)
	}
}

func TestDecisionAdmissionResultInvalidDecisionStaysInvalid(t *testing.T) {
	tests := []struct {
		name     string
		decision Decision
	}{
		{
			name:     "zero",
			decision: Decision{},
		},
		{
			name: "invalid allowed reason",
			decision: Decision{
				Allowed:  true,
				Reason:   ReasonExhausted,
				Snapshot: validGenericSnapshot(),
			},
		},
		{
			name: "invalid denied reason",
			decision: Decision{
				Allowed:  false,
				Reason:   ReasonUnknown,
				Snapshot: validGenericSnapshot(),
			},
		},
		{
			name: "denied with allowed reason",
			decision: Decision{
				Allowed:  false,
				Reason:   ReasonAllowed,
				Snapshot: validGenericSnapshot(),
			},
		},
		{
			name: "allowed with zero snapshot",
			decision: Decision{
				Allowed: true,
				Reason:  ReasonAllowed,
			},
		},
		{
			name: "denied with zero snapshot",
			decision: Decision{
				Allowed: false,
				Reason:  ReasonExhausted,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.decision.AdmissionResult(); result.IsValid() {
				t.Fatalf("AdmissionResult().IsValid() = true, want false: %+v", result.Decision())
			}
		})
	}
}

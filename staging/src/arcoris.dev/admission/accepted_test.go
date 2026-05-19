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

package admission

import "testing"

func TestAdmitDecision(t *testing.T) {
	t.Parallel()

	decision := Admit(ReasonAdmitted)
	if !decision.IsValid() {
		t.Fatalf("decision should be valid: %+v", decision)
	}
	if decision.Outcome != OutcomeAdmitted {
		t.Fatalf("outcome = %v, want %v", decision.Outcome, OutcomeAdmitted)
	}
	if decision.Effect != EffectNone {
		t.Fatalf("effect = %v, want %v", decision.Effect, EffectNone)
	}
}

func TestAdmittedDecision(t *testing.T) {
	t.Parallel()

	if got := Admitted(); got != Admit(ReasonAdmitted) {
		t.Fatalf("Admitted = %+v, want default admitted decision", got)
	}
}

func TestAcceptedResult(t *testing.T) {
	t.Parallel()

	result := Accepted(ReasonAdmitted, "snapshot")
	if !result.IsValid() {
		t.Fatalf("accepted result should be valid: %+v", result.Decision())
	}
	if !result.IsAdmitted() {
		t.Fatal("accepted result should be admitted")
	}
	if result.HasGrant() {
		t.Fatal("accepted result should not carry a grant")
	}
	if metadata, ok := result.Metadata(); !ok || metadata != "snapshot" {
		t.Fatalf("metadata = (%q, %v), want (snapshot, true)", metadata, ok)
	}
}

func TestAcceptedNoMetadataResult(t *testing.T) {
	t.Parallel()

	result := AcceptedNoMetadata(ReasonAdmitted)
	if !result.IsValid() {
		t.Fatalf("accepted result should be valid: %+v", result.Decision())
	}
	if result.HasMetadata() {
		t.Fatal("accepted no-metadata result should not carry metadata")
	}
}

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

	requireDecision(t, AdmitDecision(ReasonAdmitted), Decision{
		Outcome: OutcomeAdmitted,
		Reason:  ReasonAdmitted,
		Effect:  EffectNone,
	})
}

func TestAdmittedDecision(t *testing.T) {
	t.Parallel()

	requireDecision(t, AdmittedDecision(), Decision{
		Outcome: OutcomeAdmitted,
		Reason:  ReasonAdmitted,
		Effect:  EffectNone,
	})
}

func TestAcceptedResult(t *testing.T) {
	t.Parallel()

	result := AcceptedResult(ReasonAdmitted, "metadata")
	requireResultShape(t, result, AdmitDecision(ReasonAdmitted), false, true)
	if metadata, ok := result.Metadata(); !ok || metadata != "metadata" {
		t.Fatalf("Metadata() = (%q, %t), want metadata,true", metadata, ok)
	}
}

func TestAcceptedNoMetadataResult(t *testing.T) {
	t.Parallel()

	result := AcceptedNoMetadataResult(ReasonAdmitted)
	requireResultShape(t, result, AdmitDecision(ReasonAdmitted), false, false)
}

func TestAcceptedConstructorsWithInvalidReasonReturnInvalidValues(t *testing.T) {
	t.Parallel()

	invalid := Reason("bad-reason")
	if AdmitDecision(invalid).IsValid() {
		t.Fatal("AdmitDecision with invalid reason is valid")
	}
	if AcceptedResult(invalid, "metadata").IsValid() {
		t.Fatal("AcceptedResult with invalid reason is valid")
	}
	if AcceptedNoMetadataResult(invalid).IsValid() {
		t.Fatal("AcceptedNoMetadataResult with invalid reason is valid")
	}
}

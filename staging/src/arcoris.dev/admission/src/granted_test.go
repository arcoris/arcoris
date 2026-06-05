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

func TestGrantDecision(t *testing.T) {
	t.Parallel()

	requireDecision(t, GrantDecision(ReasonAdmitted), Decision{
		Outcome: OutcomeAdmitted,
		Reason:  ReasonAdmitted,
		Effect:  EffectOwned,
	})
}

func TestGrantedResult(t *testing.T) {
	t.Parallel()

	result := GrantedResult(ReasonAdmitted, "grant", "metadata")
	requireResultShape(t, result, GrantDecision(ReasonAdmitted), true, true)
	if grant, ok := result.Grant(); !ok || grant != "grant" {
		t.Fatalf("Grant() = (%q, %t), want grant,true", grant, ok)
	}
}

func TestGrantedNoMetadataResult(t *testing.T) {
	t.Parallel()

	result := GrantedNoMetadataResult(ReasonAdmitted, "grant")
	requireResultShape(t, result, GrantDecision(ReasonAdmitted), true, false)
}

func TestGrantedConstructorsWithInvalidReasonReturnInvalidValues(t *testing.T) {
	t.Parallel()

	invalid := Reason("bad-reason")
	if GrantDecision(invalid).IsValid() {
		t.Fatal("GrantDecision with invalid reason is valid")
	}
	if GrantedResult(invalid, "grant", "metadata").IsValid() {
		t.Fatal("GrantedResult with invalid reason is valid")
	}
	if GrantedNoMetadataResult(invalid, "grant").IsValid() {
		t.Fatal("GrantedNoMetadataResult with invalid reason is valid")
	}
}

func TestGrantedNoMetadataResultDoesNotRetainMetadataReferences(t *testing.T) {
	t.Parallel()

	type metadata struct{ value string }
	result := Result[string, *metadata]{
		decision: GrantDecision(ReasonAdmitted),
		grant:    "grant",
		hasGrant: true,
		metadata: &metadata{
			value: "should not be visible",
		},
		hasMetadata: false,
	}

	if got, ok := result.Metadata(); ok || got != nil {
		t.Fatalf("Metadata() = (%v, %t), want nil,false", got, ok)
	}
}

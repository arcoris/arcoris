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

	decision := Grant(ReasonAdmitted)
	if !decision.IsValid() {
		t.Fatalf("decision should be valid: %+v", decision)
	}
	if !decision.RequiresGrant() {
		t.Fatal("owned decision should require a grant")
	}
	if !decision.AllowsGrant() {
		t.Fatal("owned decision should allow a grant")
	}
}

func TestGrantedResult(t *testing.T) {
	t.Parallel()

	result := Granted(
		ReasonAdmitted,
		"lease",
		"snapshot",
	)
	if !result.IsValid() {
		t.Fatalf("granted result should be valid: %+v", result.Decision())
	}
	if grant, ok := result.Grant(); !ok || grant != "lease" {
		t.Fatalf("grant = (%q, %v), want (lease, true)", grant, ok)
	}
	if metadata, ok := result.Metadata(); !ok || metadata != "snapshot" {
		t.Fatalf("metadata = (%q, %v), want (snapshot, true)", metadata, ok)
	}
}

func TestGrantedNoMetadataResult(t *testing.T) {
	t.Parallel()

	result := GrantedNoMetadata(ReasonAdmitted, "lease")
	if !result.IsValid() {
		t.Fatalf("granted result should be valid: %+v", result.Decision())
	}
	if result.HasMetadata() {
		t.Fatal("granted no-metadata result should not carry metadata")
	}
}

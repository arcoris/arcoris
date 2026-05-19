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

func TestCommitDecision(t *testing.T) {
	t.Parallel()

	decision := Commit(ReasonAdmitted)
	if !decision.IsValid() {
		t.Fatalf("decision should be valid: %+v", decision)
	}
	if !decision.IsAdmitted() {
		t.Fatal("committed decision should admit work")
	}
	if !decision.HasSideEffect() {
		t.Fatal("committed decision should record a side effect")
	}
	if decision.AllowsGrant() {
		t.Fatal("committed decision should not allow caller-owned grants")
	}
}

func TestCommittedResult(t *testing.T) {
	t.Parallel()

	result := Committed(ReasonAdmitted, "budget-snapshot")
	if !result.IsValid() {
		t.Fatalf("committed result should be valid: %+v", result.Decision())
	}
	if got := result.Decision(); got != Commit(ReasonAdmitted) {
		t.Fatalf("decision = %+v, want admitted committed decision", got)
	}
	if !result.IsAdmitted() {
		t.Fatal("committed result should admit work")
	}
	if !result.HasSideEffect() {
		t.Fatal("committed result should record a side effect")
	}
	if result.HasGrant() {
		t.Fatal("committed result should not carry a grant")
	}
	if !result.HasMetadata() {
		t.Fatal("committed result should carry metadata")
	}
	if metadata, ok := result.Metadata(); !ok || metadata != "budget-snapshot" {
		t.Fatalf("metadata = (%q, %v), want (budget-snapshot, true)", metadata, ok)
	}
}

func TestCommittedNoMetadataResult(t *testing.T) {
	t.Parallel()

	result := CommittedNoMetadata(ReasonAdmitted)
	if !result.IsValid() {
		t.Fatalf("committed result should be valid: %+v", result.Decision())
	}
	if result.HasGrant() {
		t.Fatal("committed no-metadata result should not carry a grant")
	}
	if result.HasMetadata() {
		t.Fatal("committed no-metadata result should not carry metadata")
	}
}

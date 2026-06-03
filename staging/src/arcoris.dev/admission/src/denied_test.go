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

func TestDenyDecision(t *testing.T) {
	t.Parallel()

	decision := Deny(Reason("capacity_exhausted"))
	if !decision.IsValid() {
		t.Fatalf("decision should be valid: %+v", decision)
	}
	if !decision.IsDenied() {
		t.Fatal("denied decision should reject the attempt")
	}
	if decision.HasSideEffect() {
		t.Fatal("denied decision should not record side effects")
	}
}

func TestDeniedResult(t *testing.T) {
	t.Parallel()

	result := Denied(Reason("capacity_exhausted"), "snapshot")
	if !result.IsValid() {
		t.Fatalf("denied result should be valid: %+v", result.Decision())
	}
	if got := result.Decision(); got != Deny(Reason("capacity_exhausted")) {
		t.Fatalf("decision = %+v, want denied no-effect decision", got)
	}
	if result.HasGrant() {
		t.Fatal("denied result should not carry a grant")
	}
	if !result.HasMetadata() {
		t.Fatal("denied result should carry metadata")
	}
	if metadata, ok := result.Metadata(); !ok || metadata != "snapshot" {
		t.Fatalf("metadata = (%q, %v), want (snapshot, true)", metadata, ok)
	}
}

func TestDeniedForResult(t *testing.T) {
	t.Parallel()

	result := DeniedFor[string](Reason("capacity_exhausted"), "snapshot")
	if !result.IsValid() {
		t.Fatalf("denied result should be valid: %+v", result.Decision())
	}
	if !result.HasMetadata() {
		t.Fatal("denied typed result should carry metadata")
	}
	if grant, ok := result.Grant(); ok || grant != "" {
		t.Fatalf("grant = (%q, %v), want zero value and false", grant, ok)
	}
}

func TestDeniedNoMetadataResult(t *testing.T) {
	t.Parallel()

	result := DeniedNoMetadata(Reason("capacity_exhausted"))
	if !result.IsValid() {
		t.Fatalf("denied result should be valid: %+v", result.Decision())
	}
	if result.HasGrant() {
		t.Fatal("denied no-metadata result should not carry a grant")
	}
	if result.HasMetadata() {
		t.Fatal("denied no-metadata result should not carry metadata")
	}
}

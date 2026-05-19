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

func TestDeferDecision(t *testing.T) {
	t.Parallel()

	decision := Defer(ReasonDeferred)
	if !decision.IsValid() {
		t.Fatalf("decision should be valid: %+v", decision)
	}
	if !decision.IsDeferred() {
		t.Fatal("deferred decision should leave retry ownership with the caller")
	}
	if decision.HasSideEffect() {
		t.Fatal("deferred decision should not record side effects")
	}
}

func TestDeferredResult(t *testing.T) {
	t.Parallel()

	result := Deferred(ReasonDeferred, "snapshot")
	if !result.IsValid() {
		t.Fatalf("deferred result should be valid: %+v", result.Decision())
	}
	if got := result.Decision(); got != Defer(ReasonDeferred) {
		t.Fatalf("decision = %+v, want deferred no-effect decision", got)
	}
	if result.HasGrant() {
		t.Fatal("deferred result should not carry a grant")
	}
	if !result.HasMetadata() {
		t.Fatal("deferred result should carry metadata")
	}
	if metadata, ok := result.Metadata(); !ok || metadata != "snapshot" {
		t.Fatalf("metadata = (%q, %v), want (snapshot, true)", metadata, ok)
	}
}

func TestDeferredForResult(t *testing.T) {
	t.Parallel()

	result := DeferredFor[string](ReasonDeferred, "snapshot")
	if !result.IsValid() {
		t.Fatalf("deferred result should be valid: %+v", result.Decision())
	}
	if !result.HasMetadata() {
		t.Fatal("deferred typed result should carry metadata")
	}
	if grant, ok := result.Grant(); ok || grant != "" {
		t.Fatalf("grant = (%q, %v), want zero value and false", grant, ok)
	}
}

func TestDeferredNoMetadataResult(t *testing.T) {
	t.Parallel()

	result := DeferredNoMetadata(ReasonDeferred)
	if !result.IsValid() {
		t.Fatalf("deferred result should be valid: %+v", result.Decision())
	}
	if result.HasGrant() {
		t.Fatal("deferred no-metadata result should not carry a grant")
	}
	if result.HasMetadata() {
		t.Fatal("deferred no-metadata result should not carry metadata")
	}
}

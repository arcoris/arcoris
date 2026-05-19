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

func TestQueueDecision(t *testing.T) {
	t.Parallel()

	decision := Queue(ReasonQueued)
	if !decision.IsValid() {
		t.Fatalf("decision should be valid: %+v", decision)
	}
	if !decision.IsQueued() {
		t.Fatal("queued decision should accept system-owned waiting")
	}
	if decision.IsTerminal() {
		t.Fatal("queued decision should not be terminal")
	}
	if !decision.AllowsGrant() {
		t.Fatal("queued decision should allow an optional queue handle")
	}
}

func TestQueuedResult(t *testing.T) {
	t.Parallel()

	result := Queued(
		ReasonQueued,
		"ticket",
		"snapshot",
	)
	if !result.IsValid() {
		t.Fatalf("queued result should be valid: %+v", result.Decision())
	}
	if grant, ok := result.Grant(); !ok || grant != "ticket" {
		t.Fatalf("grant = (%q, %v), want (ticket, true)", grant, ok)
	}
	if metadata, ok := result.Metadata(); !ok || metadata != "snapshot" {
		t.Fatalf("metadata = (%q, %v), want (snapshot, true)", metadata, ok)
	}
}

func TestQueuedNoGrantResult(t *testing.T) {
	t.Parallel()

	result := QueuedNoGrant(ReasonQueued, "snapshot")
	if !result.IsValid() {
		t.Fatalf("queued result should be valid: %+v", result.Decision())
	}
	if result.HasGrant() {
		t.Fatal("queued no-grant result should not carry a grant")
	}
	if !result.HasMetadata() {
		t.Fatal("queued no-grant result should carry metadata")
	}
}

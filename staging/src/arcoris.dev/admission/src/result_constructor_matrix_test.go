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

func TestResultConstructorsReturnValidResults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		result   Result[string, string]
		grant    bool
		metadata bool
	}{
		{name: "accepted", result: acceptedString(ReasonAdmitted), metadata: true},
		{name: "accepted no metadata", result: acceptedStringNoMetadata(ReasonAdmitted)},
		{name: "committed", result: committedString(ReasonAdmitted), metadata: true},
		{name: "committed no metadata", result: committedStringNoMetadata(ReasonAdmitted)},
		{name: "granted", result: grantedString(ReasonAdmitted), grant: true, metadata: true},
		{name: "granted no metadata", result: grantedStringNoMetadata(ReasonAdmitted), grant: true},
		{name: "denied", result: deniedString(ReasonDenied), metadata: true},
		{name: "denied for", result: deniedForString(ReasonDenied), metadata: true},
		{name: "denied no metadata", result: deniedStringNoMetadata(ReasonDenied)},
		{name: "deferred", result: deferredString(ReasonDeferred), metadata: true},
		{name: "deferred for", result: deferredForString(ReasonDeferred), metadata: true},
		{name: "deferred no metadata", result: deferredStringNoMetadata(ReasonDeferred)},
		{name: "queued", result: queuedString(ReasonQueued), grant: true, metadata: true},
		{name: "queued no grant", result: queuedStringNoGrant(ReasonQueued), metadata: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if !tt.result.IsValid() {
				t.Fatalf("constructor result is invalid: %+v", tt.result.Decision())
			}
			if got := tt.result.HasGrant(); got != tt.grant {
				t.Fatalf("HasGrant() = %t, want %t", got, tt.grant)
			}
			if got := tt.result.HasMetadata(); got != tt.metadata {
				t.Fatalf("HasMetadata() = %t, want %t", got, tt.metadata)
			}
		})
	}
}

func TestResultConstructorsWithInvalidReasonReturnInvalidResults(t *testing.T) {
	t.Parallel()

	invalid := Reason("bad-reason")
	tests := []struct {
		name   string
		result Result[string, string]
	}{
		{name: "accepted", result: acceptedString(invalid)},
		{name: "accepted no metadata", result: acceptedStringNoMetadata(invalid)},
		{name: "committed", result: committedString(invalid)},
		{name: "committed no metadata", result: committedStringNoMetadata(invalid)},
		{name: "granted", result: grantedString(invalid)},
		{name: "granted no metadata", result: grantedStringNoMetadata(invalid)},
		{name: "denied", result: deniedString(invalid)},
		{name: "denied for", result: deniedForString(invalid)},
		{name: "denied no metadata", result: deniedStringNoMetadata(invalid)},
		{name: "deferred", result: deferredString(invalid)},
		{name: "deferred for", result: deferredForString(invalid)},
		{name: "deferred no metadata", result: deferredStringNoMetadata(invalid)},
		{name: "queued", result: queuedString(invalid)},
		{name: "queued no grant", result: queuedStringNoGrant(invalid)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.result.IsValid() {
				t.Fatalf("constructor result with invalid reason is valid: %+v", tt.result.Decision())
			}
		})
	}
}

func TestDeniedForAndDeferredForDoNotManufactureGrant(t *testing.T) {
	t.Parallel()

	denied := DeniedFor[string](ReasonDenied, "metadata")
	if denied.HasGrant() {
		t.Fatal("DeniedFor manufactured a grant")
	}
	if grant, ok := denied.Grant(); ok || grant != "" {
		t.Fatalf("DeniedFor Grant() = (%q, %t), want zero,false", grant, ok)
	}

	deferred := DeferredFor[string](ReasonDeferred, "metadata")
	if deferred.HasGrant() {
		t.Fatal("DeferredFor manufactured a grant")
	}
	if grant, ok := deferred.Grant(); ok || grant != "" {
		t.Fatalf("DeferredFor Grant() = (%q, %t), want zero,false", grant, ok)
	}
}

func TestQueuedAllowsGrantButDoesNotRequireGrant(t *testing.T) {
	t.Parallel()

	withGrant := Queued(ReasonQueued, "ticket", "metadata")
	if !withGrant.IsValid() || !withGrant.HasGrant() {
		t.Fatalf("Queued with grant = valid %t grant %t, want true,true", withGrant.IsValid(), withGrant.HasGrant())
	}

	withoutGrant := QueuedNoGrant(ReasonQueued, "metadata")
	if !withoutGrant.IsValid() || withoutGrant.HasGrant() {
		t.Fatalf("QueuedNoGrant = valid %t grant %t, want true,false", withoutGrant.IsValid(), withoutGrant.HasGrant())
	}
}

func TestGrantedRequiresGrantAtConstructionShape(t *testing.T) {
	t.Parallel()

	result := resultWith[string, string](
		Grant(ReasonAdmitted),
		noneString(),
		someString("metadata"),
	)
	if result.IsValid() {
		t.Fatal("owned result without grant is valid")
	}
}

func acceptedString(reason Reason) Result[string, string] {
	result := Accepted(reason, "metadata")
	return resultWith[string, string](result.Decision(), noneString(), someString("metadata"))
}

func acceptedStringNoMetadata(reason Reason) Result[string, string] {
	result := AcceptedNoMetadata(reason)
	return resultWith[string, string](result.Decision(), noneString(), noneString())
}

func committedString(reason Reason) Result[string, string] {
	result := Committed(reason, "metadata")
	return resultWith[string, string](result.Decision(), noneString(), someString("metadata"))
}

func committedStringNoMetadata(reason Reason) Result[string, string] {
	result := CommittedNoMetadata(reason)
	return resultWith[string, string](result.Decision(), noneString(), noneString())
}

func grantedString(reason Reason) Result[string, string] {
	return Granted(reason, "grant", "metadata")
}

func grantedStringNoMetadata(reason Reason) Result[string, string] {
	result := GrantedNoMetadata(reason, "grant")
	return resultWith[string, string](result.Decision(), someString("grant"), noneString())
}

func deniedString(reason Reason) Result[string, string] {
	result := Denied(reason, "metadata")
	return resultWith[string, string](result.Decision(), noneString(), someString("metadata"))
}

func deniedForString(reason Reason) Result[string, string] {
	return DeniedFor[string](reason, "metadata")
}

func deniedStringNoMetadata(reason Reason) Result[string, string] {
	result := DeniedNoMetadata(reason)
	return resultWith[string, string](result.Decision(), noneString(), noneString())
}

func deferredString(reason Reason) Result[string, string] {
	result := Deferred(reason, "metadata")
	return resultWith[string, string](result.Decision(), noneString(), someString("metadata"))
}

func deferredForString(reason Reason) Result[string, string] {
	return DeferredFor[string](reason, "metadata")
}

func deferredStringNoMetadata(reason Reason) Result[string, string] {
	result := DeferredNoMetadata(reason)
	return resultWith[string, string](result.Decision(), noneString(), noneString())
}

func queuedString(reason Reason) Result[string, string] {
	return Queued(reason, "ticket", "metadata")
}

func queuedStringNoGrant(reason Reason) Result[string, string] {
	result := QueuedNoGrant(reason, "metadata")
	return resultWith[string, string](result.Decision(), noneString(), someString("metadata"))
}

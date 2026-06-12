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

package objectstore

import (
	"testing"

	"arcoris.dev/apimachinery/api/meta/stamp"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/value"
)

func TestStateCloneDetachesObjectMetadata(t *testing.T) {
	state := validCommittedState()
	state.Object.ObjectMeta.Deletion = &stamp.Deletion{}

	cloned := state.Clone()
	state.Object.ObjectMeta.Deletion = nil

	if cloned.Object.ObjectMeta.Deletion == nil {
		t.Fatalf("cloned deletion pointer was mutated")
	}
}

func TestStateCloneDetachesDesiredValue(t *testing.T) {
	state := validCommittedState()
	cloned := state.Clone()
	state.Object.Desired = value.StringValue("mutated")

	got, ok := cloned.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("cloned desired = %q, %v; want desired, true", got, ok)
	}
}

func TestStateCloneDetachesObservedValue(t *testing.T) {
	state := validCommittedState()
	observed := value.StringValue("observed")
	state.Object = state.Object.WithObserved(observed)

	cloned := state.Clone()
	*state.Object.Observed = value.StringValue("mutated")

	got, ok := cloned.Object.Observed.AsString()
	if !ok || got != "observed" {
		t.Fatalf("cloned observed = %q, %v; want observed, true", got, ok)
	}
}

func TestStateCloneDetachesOwnershipDocument(t *testing.T) {
	state := validCommittedState()
	state.Ownership = ownershipWithEntry()

	cloned := state.Clone()
	state.Ownership.Desired.Entries[0].Fields[0] = "$.other"

	if got, want := cloned.Ownership.Desired.Entries[0].Fields[0], objectownership.Path("$.spec"); got != want {
		t.Fatalf("cloned ownership field = %q; want %q", got, want)
	}
}

func TestStateClonePreservesOwnershipNilVsEmptyShape(t *testing.T) {
	state := validCommittedState()
	state.Ownership = objectownership.Document{
		Version: objectownership.DocumentVersionV1,
		Desired: objectownership.Surface{
			Entries: []objectownership.Entry{},
		},
	}

	cloned := state.Clone()
	if cloned.Ownership.Desired.Entries == nil {
		t.Fatalf("cloned entries = nil; want non-nil empty")
	}

	state.Ownership.Desired.Entries = nil
	cloned = state.Clone()
	if cloned.Ownership.Desired.Entries != nil {
		t.Fatalf("cloned entries = %#v; want nil", cloned.Ownership.Desired.Entries)
	}
}

func TestStateClonePreservesRevision(t *testing.T) {
	state := validCommittedState()
	state.Revision = 42

	cloned := state.Clone()

	if cloned.Revision != 42 {
		t.Fatalf("revision = %v; want 42", cloned.Revision)
	}
}

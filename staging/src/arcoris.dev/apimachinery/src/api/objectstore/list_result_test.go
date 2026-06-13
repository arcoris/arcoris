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

	"arcoris.dev/apimachinery/api/value"
)

func TestListResultHelpers(t *testing.T) {
	if !(ListResult{}).IsZero() {
		t.Fatalf("zero ListResult is not zero")
	}

	result := ListResult{
		Items:    []ListItem{{Key: validKey(), State: validCommittedState()}},
		Revision: 2,
	}
	if result.IsZero() {
		t.Fatalf("non-empty ListResult is zero")
	}
	if result.Len() != 1 {
		t.Fatalf("Len() = %d; want 1", result.Len())
	}
}

func TestListResultCloneDetachesItemsAndStates(t *testing.T) {
	result := ListResult{
		Items:    []ListItem{{Key: validKey(), State: validCommittedState()}},
		Revision: 2,
	}

	cloned := result.Clone()
	cloned.Items[0].State.Object.Desired = value.StringValue("mutated")
	cloned.Items = append(cloned.Items, ListItem{})

	if len(result.Items) != 1 {
		t.Fatalf("source items length = %d; want 1", len(result.Items))
	}
	got, ok := result.Items[0].State.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("source desired = %q, %v; want desired, true", got, ok)
	}
	if cloned.Revision != result.Revision {
		t.Fatalf("cloned revision = %v; want %v", cloned.Revision, result.Revision)
	}
}

func TestListResultClonePreservesNilItems(t *testing.T) {
	result := ListResult{Revision: 1}

	cloned := result.Clone()

	if cloned.Items != nil {
		t.Fatalf("cloned Items = %#v; want nil", cloned.Items)
	}
	if cloned.Revision != result.Revision {
		t.Fatalf("cloned revision = %v; want %v", cloned.Revision, result.Revision)
	}
}

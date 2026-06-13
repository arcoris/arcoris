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

package objectlifecycle

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func TestListResultHelpers(t *testing.T) {
	zero := ListResult{}
	if !zero.IsZero() {
		t.Fatalf("zero IsZero() = false; want true")
	}
	if zero.Len() != 0 {
		t.Fatalf("zero Len() = %d; want 0", zero.Len())
	}

	result := ListResult{
		Items:    []objectstore.ListItem{{Key: objectstore.MustKey(testGVR(), testName(1))}},
		Revision: 1,
	}
	if result.IsZero() {
		t.Fatalf("non-zero IsZero() = true; want false")
	}
	if result.Len() != 1 {
		t.Fatalf("Len() = %d; want 1", result.Len())
	}
}

func TestListResultCloneDetachesItemsAndStates(t *testing.T) {
	original := ListResult{
		Items: []objectstore.ListItem{
			{
				Key: objectstore.MustKey(testGVR(), testName(1)),
				State: objectstore.State{
					Object:   testObjectWithDesired(1, value.StringValue("original")),
					Revision: 1,
				},
			},
		},
		Revision: 2,
	}

	cloned := original.Clone()
	cloned.Items[0].State.Object.Desired = value.StringValue("clone-mutated")
	cloned.Items[0] = objectstore.ListItem{}

	if original.Len() != 1 {
		t.Fatalf("original len = %d; want 1", original.Len())
	}
	got, ok := original.Items[0].State.Object.Desired.AsString()
	if !ok || got != "original" {
		t.Fatalf("original desired = %q, %v; want original, true", got, ok)
	}
	if cloned.Revision != original.Revision {
		t.Fatalf("clone revision = %v; want %v", cloned.Revision, original.Revision)
	}
}

func TestListResultClonePreservesNilItems(t *testing.T) {
	clone := ListResult{Revision: 1}.Clone()

	if clone.Items != nil {
		t.Fatalf("clone items = %#v; want nil", clone.Items)
	}
	if clone.Revision != 1 {
		t.Fatalf("clone revision = %v; want 1", clone.Revision)
	}
}

/*
   Copyright 2026 The ARCORIS Authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package schema

import "testing"

// ObjectKind tests keep the interface schema-only and verify the lightweight
// utility implementations do not require any runtime object dependency.

// TestObjectKindHolderStoresAndReturnsGVK verifies mutable schema identity storage.
func TestObjectKindHolderStoresAndReturnsGVK(t *testing.T) {
	initial := GroupVersionKind{Version: "v1", Kind: "Pod"}
	holder := NewObjectKindHolder(initial)
	if got := holder.GroupVersionKind(); got != initial {
		t.Fatalf("initial GroupVersionKind = %+v, want %+v", got, initial)
	}

	updated := GroupVersionKind{Group: "control.arcoris.dev", Version: "v1alpha1", Kind: "WorkloadClass"}
	holder.SetGroupVersionKind(updated)
	if got := holder.GroupVersionKind(); got != updated {
		t.Fatalf("updated GroupVersionKind = %+v, want %+v", got, updated)
	}
}

// TestEmptyObjectKindNoOpBehavior verifies the no-op implementation never stores a kind.
func TestEmptyObjectKindNoOpBehavior(t *testing.T) {
	var objectKind ObjectKind = EmptyObjectKind{}
	objectKind.SetGroupVersionKind(GroupVersionKind{Version: "v1", Kind: "Pod"})
	if got := objectKind.GroupVersionKind(); got != (GroupVersionKind{}) {
		t.Fatalf("GroupVersionKind = %+v, want zero", got)
	}
}

// TestNilObjectKindHolderNoOpBehavior verifies nil holders behave like empty holders.
func TestNilObjectKindHolderNoOpBehavior(t *testing.T) {
	var holder *ObjectKindHolder
	holder.SetGroupVersionKind(GroupVersionKind{Version: "v1", Kind: "Pod"})
	if got := holder.GroupVersionKind(); got != (GroupVersionKind{}) {
		t.Fatalf("nil holder GroupVersionKind = %+v, want zero", got)
	}
}

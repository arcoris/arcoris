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

package admissioncatalog

import (
	"testing"

	"arcoris.dev/admission"
)

func TestDescriptorStoreDeclareLookupListAndClone(t *testing.T) {
	store := newReasonStore()
	if !store.declare(reasonDescriptor(testOtherReason)) {
		t.Fatal("declare returned false for first descriptor")
	}
	if !store.declare(reasonDescriptor(testReason)) {
		t.Fatal("declare returned false for second descriptor")
	}
	if store.declare(reasonDescriptor(testReason)) {
		t.Fatal("declare returned true for duplicate descriptor")
	}

	if !store.has(testReason) {
		t.Fatal("store does not contain reason")
	}
	if descriptor, ok := store.get(testReason); !ok || descriptor.Reason != testReason {
		t.Fatalf("get = %+v, %v", descriptor, ok)
	}
	list := store.list()
	if got, want := list[0].Reason, testReason; got != want {
		t.Fatalf("first listed reason = %s, want %s", got, want)
	}

	clone := store.clone()
	list[0] = reasonDescriptor(admission.Reason("mutated_reason"))
	if clone.has(admission.Reason("mutated_reason")) {
		t.Fatal("mutating returned list changed cloned store")
	}
	store.byKey[testReason] = reasonDescriptor(admission.Reason("replacement_reason"))
	if descriptor, _ := clone.get(testReason); descriptor.Reason != testReason {
		t.Fatal("clone changed after source mutation")
	}
}

func TestDescriptorStoreZeroValueReadsEmpty(t *testing.T) {
	var store descriptorStore[admission.Reason, ReasonDescriptor]
	if store.len() != 0 {
		t.Fatalf("len = %d, want 0", store.len())
	}
	if store.has(testReason) {
		t.Fatal("zero store reports descriptor")
	}
	if list := store.list(); len(list) != 0 {
		t.Fatalf("list length = %d, want 0", len(list))
	}
}

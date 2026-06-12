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

package object

import "testing"

func TestNew(t *testing.T) {
	desired := testDesired{Replicas: 3}
	typeMeta := validTypeMeta()
	objectMeta := validObjectMeta()

	obj := New[testDesired, testObserved](
		typeMeta,
		objectMeta,
		desired,
	)

	if obj.TypeMeta != typeMeta {
		t.Fatalf("TypeMeta = %#v", obj.TypeMeta)
	}
	if obj.ObjectMeta.Name != objectMeta.Name || obj.ObjectMeta.Namespace != objectMeta.Namespace {
		t.Fatalf("ObjectMeta = %#v", obj.ObjectMeta)
	}
	if obj.Desired != desired {
		t.Fatalf("Desired = %#v", obj.Desired)
	}
	if obj.HasObserved() {
		t.Fatal("New() HasObserved() = true")
	}
}

func TestNewCopiesMetadata(t *testing.T) {
	typeMeta := validTypeMeta()
	objectMeta := validObjectMeta()
	obj := New[testDesired, testObserved](
		typeMeta,
		objectMeta,
		testDesired{Replicas: 3},
	)

	typeMeta.Kind = "Other"
	objectMeta.Labels["role"] = "mutated"

	if obj.TypeMeta.Kind != "Worker" {
		t.Fatalf("stored kind = %q, want detached type metadata", obj.TypeMeta.Kind)
	}
	if got := obj.ObjectMeta.Labels["role"]; got != "worker" {
		t.Fatalf("stored label = %q, want detached metadata", got)
	}
}

func TestNewObserved(t *testing.T) {
	observed := testObserved{ReadyReplicas: 2}
	obj := NewObserved(
		validTypeMeta(),
		validObjectMeta(),
		testDesired{Replicas: 3},
		observed,
	)

	if !obj.HasObserved() {
		t.Fatal("NewObserved() HasObserved() = false")
	}
	if obj.Observed == nil || *obj.Observed != observed {
		t.Fatalf("Observed = %#v", obj.Observed)
	}

	observed.ReadyReplicas = 9
	if obj.Observed.ReadyReplicas != 2 {
		t.Fatalf("Observed changed after caller mutation: %#v", obj.Observed)
	}
	if obj.Observed == &observed {
		t.Fatal("NewObserved() reused caller variable address")
	}
}

func TestObservedValue(t *testing.T) {
	obj := New[testDesired, testObserved](
		validTypeMeta(),
		validObjectMeta(),
		testDesired{Replicas: 3},
	)

	observed, ok := obj.ObservedValue()
	if ok {
		t.Fatalf("ObservedValue() ok = true, value = %#v", observed)
	}

	withObserved := obj.WithObserved(testObserved{ReadyReplicas: 2})
	observed, ok = withObserved.ObservedValue()
	if !ok {
		t.Fatal("ObservedValue() ok = false")
	}
	if observed.ReadyReplicas != 2 {
		t.Fatalf("ObservedValue() = %#v", observed)
	}

	observed.ReadyReplicas = 9
	if withObserved.Observed.ReadyReplicas != 2 {
		t.Fatal("ObservedValue() exposed internal observed pointer")
	}
}

func TestObservedValueIsShallow(t *testing.T) {
	type observedPayload struct {
		Values []int
	}

	obj := NewObserved(
		validTypeMeta(),
		validObjectMeta(),
		testDesired{Replicas: 3},
		observedPayload{Values: []int{1}},
	)

	observed, ok := obj.ObservedValue()
	if !ok {
		t.Fatal("ObservedValue() ok = false")
	}
	observed.Values[0] = 9
	if obj.Observed.Values[0] != 9 {
		t.Fatal("ObservedValue() unexpectedly deep-copied reference-bearing payload")
	}
}

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

import (
	"testing"

	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/meta/labels"
)

func TestObjectUpdateMethods(t *testing.T) {
	original := New[testDesired, testObserved](
		validTypeMeta(),
		validObjectMeta(),
		testDesired{Replicas: 3},
	)

	otherType := meta.TypeMeta{Kind: "Other"}
	otherMeta := meta.ObjectMeta{
		Name:   "other",
		Labels: labels.Set{"role": "other"},
	}

	withType := original.WithTypeMeta(otherType)
	if withType.TypeMeta != otherType {
		t.Fatalf("WithTypeMeta() TypeMeta = %#v", withType.TypeMeta)
	}
	if original.TypeMeta != validTypeMeta() {
		t.Fatal("WithTypeMeta() mutated original")
	}

	withMeta := original.WithObjectMeta(otherMeta)
	if withMeta.ObjectMeta.Name != "other" {
		t.Fatalf("WithObjectMeta() ObjectMeta = %#v", withMeta.ObjectMeta)
	}
	if original.ObjectMeta.Name != "main" {
		t.Fatal("WithObjectMeta() mutated original")
	}

	otherMeta.Labels["role"] = "mutated"
	if withMeta.ObjectMeta.Labels["role"] != "other" {
		t.Fatal("WithObjectMeta() did not detach replacement metadata")
	}

	withDesired := original.WithDesired(testDesired{Replicas: 5})
	if withDesired.Desired.Replicas != 5 {
		t.Fatalf("WithDesired() Desired = %#v", withDesired.Desired)
	}
	if original.Desired.Replicas != 3 {
		t.Fatal("WithDesired() mutated original")
	}

	withObserved := original.WithObserved(testObserved{ReadyReplicas: 2})
	if !withObserved.HasObserved() || withObserved.Observed.ReadyReplicas != 2 {
		t.Fatalf("WithObserved() Observed = %#v", withObserved.Observed)
	}
	observedValue, ok := withObserved.ObservedValue()
	if !ok || observedValue.ReadyReplicas != 2 {
		t.Fatalf("WithObserved() ObservedValue() = %#v, %v", observedValue, ok)
	}
	if original.HasObserved() {
		t.Fatal("WithObserved() mutated original")
	}

	observed := testObserved{ReadyReplicas: 7}
	withObserved = original.WithObserved(observed)
	observed.ReadyReplicas = 9
	if withObserved.Observed.ReadyReplicas != 7 {
		t.Fatal("WithObserved() reused caller variable address")
	}

	withoutObserved := withObserved.WithoutObserved()
	if withoutObserved.HasObserved() {
		t.Fatal("WithoutObserved() HasObserved() = true")
	}
	if !withObserved.HasObserved() {
		t.Fatal("WithoutObserved() mutated original")
	}
}

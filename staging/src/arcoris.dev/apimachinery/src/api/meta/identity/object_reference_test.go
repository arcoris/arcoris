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

package identity

import (
	"encoding/json"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
)

func TestObjectReference(t *testing.T) {
	ref := ObjectReference{
		APIVersion: apiidentity.GroupVersion{Group: "control.arcoris.dev", Version: "v1"},
		Kind:       "Worker",
		Namespace:  "system",
		Name:       "worker",
		UID:        "uid-1",
	}

	if ref.ObjectName().String() != "system/worker" {
		t.Fatalf("ObjectName() = %q", ref.ObjectName())
	}
	if ref.ObjectIdentity().String() != "system/worker#uid-1" {
		t.Fatalf("ObjectIdentity() = %q", ref.ObjectIdentity())
	}
	if ref.GroupVersionKind().String() != "control.arcoris.dev/v1#Worker" {
		t.Fatalf("GroupVersionKind() = %q", ref.GroupVersionKind())
	}
	if ref.IsZero() {
		t.Fatal("non-zero ObjectReference IsZero() = true")
	}
	if !(ObjectReference{}).IsZero() {
		t.Fatal("zero ObjectReference IsZero() = false")
	}
}

func TestObjectReferenceJSONFields(t *testing.T) {
	data, err := json.Marshal(ObjectReference{
		APIVersion: apiidentity.GroupVersion{Group: "control.arcoris.dev", Version: "v1"},
		Kind:       "Worker",
		Namespace:  "system",
		Name:       "worker",
		UID:        "uid-1",
	})
	requireNoError(t, err)

	var got map[string]any
	requireNoError(t, json.Unmarshal(data, &got))

	if got["apiVersion"] != "control.arcoris.dev/v1" {
		t.Fatalf("apiVersion = %#v", got["apiVersion"])
	}
	if got["kind"] != "Worker" {
		t.Fatalf("kind = %#v", got["kind"])
	}
	if got["namespace"] != "system" || got["name"] != "worker" || got["uid"] != "uid-1" {
		t.Fatalf("object reference JSON = %#v", got)
	}
	if _, ok := got["APIVersion"]; ok {
		t.Fatalf("unexpected Go field name in JSON: %s", data)
	}
}

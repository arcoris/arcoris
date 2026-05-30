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

package meta

import (
	"encoding/json"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
)

func TestTypeMeta(t *testing.T) {
	gvk := apiidentity.GroupVersionKind{Group: "control.arcoris.dev", Version: "v1", Kind: "Worker"}
	meta := FromGroupVersionKind(gvk)

	if got := meta.GroupVersionKind(); got != gvk {
		t.Fatalf("GroupVersionKind() = %#v, want %#v", got, gvk)
	}
	if got := meta.String(); got != gvk.String() {
		t.Fatalf("String() = %q, want %q", got, gvk.String())
	}
	if !(TypeMeta{}).IsZero() {
		t.Fatal("zero TypeMeta IsZero() = false")
	}
	if got := (TypeMeta{}).String(); got != "" {
		t.Fatalf("zero TypeMeta String() = %q, want empty", got)
	}
	if meta.IsZero() {
		t.Fatal("non-zero TypeMeta IsZero() = true")
	}
}

func TestTypeMetaJSONFields(t *testing.T) {
	type testObject struct {
		TypeMeta `json:",inline"`
	}

	data, err := json.Marshal(testObject{
		TypeMeta: FromGroupVersionKind(apiidentity.GroupVersionKind{
			Group:   "control.arcoris.dev",
			Version: "v1",
			Kind:    "Worker",
		}),
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
	if _, ok := got["APIVersion"]; ok {
		t.Fatalf("unexpected Go field name in JSON: %s", data)
	}
	if _, ok := got["Kind"]; ok {
		t.Fatalf("unexpected Go field name in JSON: %s", data)
	}
}

func TestTypeMetaJSONOmitsZeroAPIVersion(t *testing.T) {
	data, err := json.Marshal(TypeMeta{})
	requireNoError(t, err)

	var got map[string]any
	requireNoError(t, json.Unmarshal(data, &got))
	if _, ok := got["apiVersion"]; ok {
		t.Fatalf("zero APIVersion encoded in JSON: %s", data)
	}
}

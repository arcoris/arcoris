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
)

func TestObjectName(t *testing.T) {
	global := ObjectName{Name: "worker"}
	if global.String() != "worker" {
		t.Fatalf("String() = %q", global.String())
	}
	if global.IsZero() {
		t.Fatal("global ObjectName IsZero() = true")
	}

	namespaced := ObjectName{Namespace: "system", Name: "worker"}
	if namespaced.String() != "system/worker" {
		t.Fatalf("String() = %q", namespaced.String())
	}
	if namespaced.IsZero() {
		t.Fatal("namespaced ObjectName IsZero() = true")
	}
	if !(ObjectName{}).IsZero() {
		t.Fatal("zero ObjectName IsZero() = false")
	}
}

func TestObjectNameCanonicalText(t *testing.T) {
	text, err := (ObjectName{Namespace: "system", Name: "worker"}).CanonicalText()
	requireNoError(t, err)
	if text != "system/worker" {
		t.Fatalf("CanonicalText() = %q, want %q", text, "system/worker")
	}

	text, err = (ObjectName{Name: "worker"}).CanonicalText()
	requireNoError(t, err)
	if text != "worker" {
		t.Fatalf("CanonicalText() = %q, want %q", text, "worker")
	}

	_, err = (ObjectName{Namespace: "system"}).CanonicalText()
	requireErrorIs(t, err, ErrInvalidObjectName)
}

func TestObjectNameJSONFields(t *testing.T) {
	data, err := json.Marshal(ObjectName{Namespace: "system", Name: "worker"})
	requireNoError(t, err)

	var got map[string]any
	requireNoError(t, json.Unmarshal(data, &got))

	if got["namespace"] != "system" || got["name"] != "worker" {
		t.Fatalf("object name JSON = %#v", got)
	}
	if _, ok := got["Namespace"]; ok {
		t.Fatalf("unexpected Go field name in JSON: %s", data)
	}
}

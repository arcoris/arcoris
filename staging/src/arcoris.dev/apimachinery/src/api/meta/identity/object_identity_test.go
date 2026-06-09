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

func TestObjectIdentity(t *testing.T) {
	ident := ObjectIdentity{Namespace: "system", Name: "worker", UID: "uid-1"}
	if ident.ObjectName().String() != "system/worker" {
		t.Fatalf("ObjectName() = %q", ident.ObjectName())
	}
	if ident.String() != "system/worker#uid-1" {
		t.Fatalf("String() = %q", ident.String())
	}
	if ident.IsZero() {
		t.Fatal("non-zero ObjectIdentity IsZero() = true")
	}
	if !(ObjectIdentity{}).IsZero() {
		t.Fatal("zero ObjectIdentity IsZero() = false")
	}
}

func TestObjectIdentityCanonicalText(t *testing.T) {
	text, err := (ObjectIdentity{Namespace: "system", Name: "worker", UID: "uid-1"}).CanonicalText()
	requireNoError(t, err)
	if text != "system/worker#uid-1" {
		t.Fatalf("CanonicalText() = %q, want %q", text, "system/worker#uid-1")
	}

	_, err = (ObjectIdentity{Namespace: "system", Name: "worker"}).CanonicalText()
	requireErrorIs(t, err, ErrInvalidObjectIdentity)
}

func TestObjectIdentityJSONFields(t *testing.T) {
	data, err := json.Marshal(ObjectIdentity{Namespace: "system", Name: "worker", UID: "uid-1"})
	requireNoError(t, err)

	var got map[string]any
	requireNoError(t, json.Unmarshal(data, &got))

	if got["namespace"] != "system" || got["name"] != "worker" || got["uid"] != "uid-1" {
		t.Fatalf("object identity JSON = %#v", got)
	}
	if _, ok := got["Namespace"]; ok {
		t.Fatalf("unexpected Go field name in JSON: %s", data)
	}
}

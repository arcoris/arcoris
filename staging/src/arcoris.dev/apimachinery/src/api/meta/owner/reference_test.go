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

package owner

import (
	"encoding/json"
	"testing"
)

func TestReference(t *testing.T) {
	ref := validReference(true)
	if ref.IsZero() {
		t.Fatal("non-zero Reference IsZero() = true")
	}
	if !(Reference{}).IsZero() {
		t.Fatal("zero Reference IsZero() = false")
	}
}

func TestReferenceJSONFields(t *testing.T) {
	data, err := json.Marshal(validReference(true))
	requireNoError(t, err)

	var got map[string]any
	requireNoError(t, json.Unmarshal(data, &got))

	if _, ok := got["ref"].(map[string]any); !ok {
		t.Fatalf("ref = %#v", got["ref"])
	}
	if got["controller"] != true {
		t.Fatalf("controller = %#v", got["controller"])
	}
	if _, ok := got["Ref"]; ok {
		t.Fatalf("unexpected Go field name in JSON: %s", data)
	}
}

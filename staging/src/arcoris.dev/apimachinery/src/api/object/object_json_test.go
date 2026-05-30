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
	"encoding/json"
	"testing"
)

func TestObjectJSONShape(t *testing.T) {
	obj := NewObserved(
		validTypeMeta(),
		wireObjectMeta(),
		testDesired{Replicas: 3},
		testObserved{ReadyReplicas: 2},
	)

	data, err := json.Marshal(obj)
	requireNoError(t, err)

	got := decodeObject(t, data)
	if got["apiVersion"] != "control.arcoris.dev/v1" {
		t.Fatalf("apiVersion = %#v", got["apiVersion"])
	}
	if got["kind"] != "Worker" {
		t.Fatalf("kind = %#v", got["kind"])
	}

	metadata, ok := got["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata = %#v", got["metadata"])
	}
	if metadata["name"] != "main" || metadata["namespace"] != "system" {
		t.Fatalf("metadata = %#v", metadata)
	}

	desired, ok := got["desired"].(map[string]any)
	if !ok || desired["replicas"] != float64(3) {
		t.Fatalf("desired = %#v", got["desired"])
	}

	observed, ok := got["observed"].(map[string]any)
	if !ok || observed["readyReplicas"] != float64(2) {
		t.Fatalf("observed = %#v", got["observed"])
	}
}

func TestObjectJSONOmitsNilObserved(t *testing.T) {
	obj := New[testDesired, testObserved](
		validTypeMeta(),
		wireObjectMeta(),
		testDesired{Replicas: 3},
	)

	data, err := json.Marshal(obj)
	requireNoError(t, err)

	got := decodeObject(t, data)
	if _, ok := got["observed"]; ok {
		t.Fatalf("observed encoded for nil observed: %s", data)
	}
}

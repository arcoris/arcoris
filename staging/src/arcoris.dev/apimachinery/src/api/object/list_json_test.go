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

func TestListJSONShape(t *testing.T) {
	item := New[testDesired, testObserved](
		validTypeMeta(),
		wireObjectMeta(),
		testDesired{Replicas: 3},
	)
	list := NewList(
		validListTypeMeta(),
		validListMeta(),
		[]Object[testDesired, testObserved]{item},
	)

	data, err := json.Marshal(list)
	requireNoError(t, err)

	got := decodeObject(t, data)
	if got["apiVersion"] != "control.arcoris.dev/v1" {
		t.Fatalf("apiVersion = %#v", got["apiVersion"])
	}
	if got["kind"] != "WorkerList" {
		t.Fatalf("kind = %#v", got["kind"])
	}

	metadata, ok := got["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata = %#v", got["metadata"])
	}
	if metadata["resourceVersion"] != "rv-1" {
		t.Fatalf("resourceVersion = %#v", metadata["resourceVersion"])
	}
	if metadata["continue"] != "token-1" {
		t.Fatalf("continue = %#v", metadata["continue"])
	}
	if metadata["remainingItemCount"] != float64(1) {
		t.Fatalf("remainingItemCount = %#v", metadata["remainingItemCount"])
	}

	items, ok := got["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %#v", got["items"])
	}
}

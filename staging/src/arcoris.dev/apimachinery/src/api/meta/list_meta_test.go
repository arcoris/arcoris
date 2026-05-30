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

func TestListMeta(t *testing.T) {
	if !(ListMeta{}).IsZero() {
		t.Fatal("zero ListMeta IsZero() = false")
	}

	count := uint64(3)
	meta := ListMeta{
		ResourceVersion:    "rv-1",
		ContinueToken:      "page-1",
		RemainingItemCount: &count,
	}
	if meta.IsZero() {
		t.Fatal("non-zero ListMeta IsZero() = true")
	}
}

func TestListMetaJSONFields(t *testing.T) {
	type testWorker struct {
		TypeMeta `json:",inline"`
	}
	type testWorkerList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		Items []testWorker `json:"items"`
	}

	count := uint64(10)
	data, err := json.Marshal(testWorkerList{
		TypeMeta: FromGroupVersionKind(apiidentity.GroupVersionKind{
			Group:   "control.arcoris.dev",
			Version: "v1",
			Kind:    "WorkerList",
		}),
		ListMeta: ListMeta{
			ResourceVersion:    "rv-1",
			ContinueToken:      "token-1",
			RemainingItemCount: &count,
		},
		Items: []testWorker{
			{
				TypeMeta: FromGroupVersionKind(apiidentity.GroupVersionKind{
					Group:   "control.arcoris.dev",
					Version: "v1",
					Kind:    "Worker",
				}),
			},
		},
	})
	requireNoError(t, err)

	var got map[string]any
	requireNoError(t, json.Unmarshal(data, &got))

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
	if metadata["remainingItemCount"] != float64(10) {
		t.Fatalf("remainingItemCount = %#v", metadata["remainingItemCount"])
	}
	if _, ok := metadata["ContinueToken"]; ok {
		t.Fatalf("unexpected Go field name in JSON: %s", data)
	}
}

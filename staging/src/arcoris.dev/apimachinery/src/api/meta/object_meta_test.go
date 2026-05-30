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
	"time"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/stamp"
)

func testTime() time.Time {
	return time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
}

func TestObjectMeta(t *testing.T) {
	if !(ObjectMeta{}).IsZero() {
		t.Fatal("zero ObjectMeta IsZero() = false")
	}

	meta := validObjectMeta()
	if meta.IsZero() {
		t.Fatal("non-zero ObjectMeta IsZero() = true")
	}

	if meta.ObjectName() != (metaidentity.ObjectName{Namespace: "system", Name: "worker"}) {
		t.Fatalf("ObjectName() = %#v", meta.ObjectName())
	}
	if meta.ObjectIdentity() != (metaidentity.ObjectIdentity{Namespace: "system", Name: "worker", UID: "uid-1"}) {
		t.Fatalf("ObjectIdentity() = %#v", meta.ObjectIdentity())
	}
}

func TestObjectMetaIsZeroNonNilDeletion(t *testing.T) {
	meta := ObjectMeta{Deletion: &stamp.Deletion{}}
	if meta.IsZero() {
		t.Fatal("ObjectMeta with non-nil deletion IsZero() = true")
	}
}

func TestObjectMetaJSONFields(t *testing.T) {
	data, err := json.Marshal(ObjectMeta{
		Name:      "worker",
		Namespace: "system",
		UID:       "uid-1",
	})
	requireNoError(t, err)

	var got map[string]any
	requireNoError(t, json.Unmarshal(data, &got))

	if got["name"] != "worker" {
		t.Fatalf("name = %#v", got["name"])
	}
	if got["namespace"] != "system" {
		t.Fatalf("namespace = %#v", got["namespace"])
	}
	if got["uid"] != "uid-1" {
		t.Fatalf("uid = %#v", got["uid"])
	}
	if _, ok := got["Name"]; ok {
		t.Fatalf("unexpected Go field name in JSON: %s", data)
	}
	if _, ok := got["CreatedAt"]; ok {
		t.Fatalf("unexpected Go field name in JSON: %s", data)
	}
}

func TestObjectMetaJSONOmitsZeroCreatedAt(t *testing.T) {
	data, err := json.Marshal(ObjectMeta{})
	requireNoError(t, err)

	var got map[string]any
	requireNoError(t, json.Unmarshal(data, &got))
	if _, ok := got["createdAt"]; ok {
		t.Fatalf("zero CreatedAt encoded in JSON: %s", data)
	}
}

func TestObjectMetaEmbeddedResourceJSON(t *testing.T) {
	type testWorkerDesired struct {
		Replicas int32 `json:"replicas"`
	}
	type testWorkerObserved struct {
		ReadyReplicas int32 `json:"readyReplicas"`
	}
	type testWorker struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`

		Desired  testWorkerDesired  `json:"desired,omitempty"`
		Observed testWorkerObserved `json:"observed,omitempty"`
	}

	data, err := json.Marshal(testWorker{
		TypeMeta: FromGroupVersionKind(apiidentity.GroupVersionKind{
			Group:   "control.arcoris.dev",
			Version: "v1",
			Kind:    "Worker",
		}),
		ObjectMeta: ObjectMeta{
			Name:      "main",
			Namespace: "system",
		},
		Desired:  testWorkerDesired{Replicas: 3},
		Observed: testWorkerObserved{ReadyReplicas: 2},
	})
	requireNoError(t, err)

	var got map[string]any
	requireNoError(t, json.Unmarshal(data, &got))

	if got["apiVersion"] != "control.arcoris.dev/v1" || got["kind"] != "Worker" {
		t.Fatalf("type metadata JSON = %s", data)
	}

	metadata, ok := got["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata = %#v", got["metadata"])
	}
	if metadata["name"] != "main" || metadata["namespace"] != "system" {
		t.Fatalf("metadata JSON = %#v", metadata)
	}

	desired, ok := got["desired"].(map[string]any)
	if !ok || desired["replicas"] != float64(3) {
		t.Fatalf("desired JSON = %#v", got["desired"])
	}

	observed, ok := got["observed"].(map[string]any)
	if !ok || observed["readyReplicas"] != float64(2) {
		t.Fatalf("observed JSON = %#v", got["observed"])
	}
}

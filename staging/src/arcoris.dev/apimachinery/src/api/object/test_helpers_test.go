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
	"errors"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/meta/labels"
)

type testDesired struct {
	Replicas int32 `json:"replicas"`
}

type testObserved struct {
	ReadyReplicas int32 `json:"readyReplicas"`
}

type uninspectedPayload struct {
	Value string `json:"value"`
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func decodeObject(t *testing.T, data []byte) map[string]any {
	t.Helper()

	var out map[string]any
	requireNoError(t, json.Unmarshal(data, &out))

	return out
}

func validTypeMeta() meta.TypeMeta {
	return meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
		Group:   "control.arcoris.dev",
		Version: "v1",
		Kind:    "Worker",
	})
}

func validListTypeMeta() meta.TypeMeta {
	return meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
		Group:   "control.arcoris.dev",
		Version: "v1",
		Kind:    "WorkerList",
	})
}

func validObjectMeta() meta.ObjectMeta {
	return meta.ObjectMeta{
		Name:      "main",
		Namespace: "system",
		UID:       "uid-1",
		Labels:    labels.Set{"role": "worker"},
	}
}

func wireObjectMeta() meta.ObjectMeta {
	return meta.ObjectMeta{
		Name:      "main",
		Namespace: "system",
	}
}

func validListMeta() meta.ListMeta {
	count := uint64(1)

	return meta.ListMeta{
		ResourceVersion:    "rv-1",
		ContinueToken:      "token-1",
		RemainingItemCount: &count,
	}
}

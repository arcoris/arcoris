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

package objectlifecycle

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func TestGetExistingReturnsFound(t *testing.T) {
	executor := testExecutor(t)
	created := createObject(t, executor, 1, "api:v1", owner("creator"))

	result, err := executor.Get(
		context.Background(),
		GetRequest{Resource: testGVR(), Object: testName(1)},
	)
	requireNoError(t, err)

	requireEffect(t, result, OperationGet, EffectFound)
	if result.State.Revision != created.State.Revision {
		t.Fatalf("revision = %v; want %v", result.State.Revision, created.State.Revision)
	}
	requireImage(t, result.State, "api:v1")
}

func TestGetMissingReturnsNotFound(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.Get(
		context.Background(),
		GetRequest{Resource: testGVR(), Object: testName(1)},
	)

	requireLifecycleError(t, err, ErrNotFound, ErrorReasonNotFound)
}

func TestGetMissingResourceReturnsResourceNotFound(t *testing.T) {
	executor := testExecutor(t)
	req := GetRequest{Resource: testGVR(), Object: testName(1)}
	req.Resource.Resource = "unknowns"

	_, err := executor.Get(context.Background(), req)

	requireLifecycleError(t, err, ErrResourceNotFound, ErrorReasonResourceNotFound)
}

func TestGetDoesNotValidateDesiredPayload(t *testing.T) {
	store := testStore(t)
	key := objectstore.MustKey(testGVR(), testName(1))
	_, err := store.Create(
		context.Background(),
		key,
		objectstore.State{
			Object:    testObjectWithDesired(1, value.StringValue("descriptor-invalid")),
			Ownership: objectownership.Document{Version: objectownership.DocumentVersionV1},
		},
	)
	requireNoError(t, err)
	executor := testExecutor(t, WithStore(store))

	result, err := executor.Get(context.Background(), GetRequest{Resource: testGVR(), Object: testName(1)})
	requireNoError(t, err)

	requireEffect(t, result, OperationGet, EffectFound)
	if got, ok := result.State.Object.Desired.AsString(); !ok || got != "descriptor-invalid" {
		t.Fatalf("Desired = %q, %v; want descriptor-invalid, true", got, ok)
	}
}

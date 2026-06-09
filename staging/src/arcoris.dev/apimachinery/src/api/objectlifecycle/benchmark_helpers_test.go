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

	"arcoris.dev/apimachinery/api/objectmemorystore"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/resourcecatalog"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func benchmarkExecutor(b *testing.B) *Executor {
	b.Helper()

	store, err := objectmemorystore.New()
	if err != nil {
		b.Fatalf("objectmemorystore.New() error: %v", err)
	}
	catalog := resourcecatalog.New(nil)
	if err := catalog.Register(testDefinition()); err != nil {
		b.Fatalf("Register() error: %v", err)
	}
	executor, err := NewExecutor(
		WithStore(store),
		WithResourceResolver(catalog),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
		WithObservedValidator(valuevalidation.SurfaceValidator{}),
	)
	if err != nil {
		b.Fatalf("NewExecutor() error: %v", err)
	}

	return executor
}

func benchmarkCreateRequests(b *testing.B, n int) []CreateRequest {
	b.Helper()

	requests := make([]CreateRequest, n)
	for i := range requests {
		requests[i] = CreateRequest{Object: testObject(i+1, "api:v1"), Owner: owner("creator")}
	}

	return requests
}

func benchmarkApplyRequests(b *testing.B, n int) []ApplyRequest {
	b.Helper()

	requests := make([]ApplyRequest, n)
	for i := range requests {
		requests[i] = ApplyRequest{Object: testObject(i+1, "api:v1"), Owner: owner("creator")}
	}

	return requests
}

func benchmarkCreateObject(b *testing.B, executor *Executor, index int, image string) Result {
	b.Helper()

	result, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(index, image), Owner: owner("creator")},
	)
	if err != nil {
		b.Fatalf("Create() error: %v", err)
	}

	return result
}

func benchmarkPreparedDeletes(b *testing.B, executor *Executor, n int) []DeleteRequest {
	b.Helper()

	requests := make([]DeleteRequest, n)
	for i := range requests {
		created := benchmarkCreateObject(b, executor, i+1, "api:v1")
		requests[i] = DeleteRequest{
			Resource: testGVR(),
			Object:   testName(i + 1),
			Expected: created.State.Revision,
		}
	}

	return requests
}

func benchmarkPreparedStoreCreates(b *testing.B, executor *Executor, n int) ([]objectstore.Key, []objectstore.State) {
	b.Helper()

	keys := make([]objectstore.Key, n)
	states := make([]objectstore.State, n)
	for i := range keys {
		obj := testObject(i+1, "api:v1")
		prepared, err := executor.prepareObjectRequest(OperationCreate, obj)
		if err != nil {
			b.Fatalf("prepareObjectRequest() error: %v", err)
		}
		ownership, err := executor.initialOwnership(OperationCreate, prepared.key, owner("creator"), obj.Desired, prepared.resolved)
		if err != nil {
			b.Fatalf("initialOwnership() error: %v", err)
		}
		keys[i] = prepared.key
		states[i] = inputState(obj, ownership)
	}

	return keys, states
}

func benchmarkPreparedStoreUpdates(b *testing.B, executor *Executor, n int) ([]objectstore.Key, []objectstore.Revision, []objectstore.State) {
	b.Helper()

	keys, states := benchmarkPreparedStoreCreates(b, executor, n)
	revisions := make([]objectstore.Revision, n)
	for i := range keys {
		committed, err := executor.store.Create(context.Background(), keys[i], states[i])
		if err != nil {
			b.Fatalf("store.Create() error: %v", err)
		}
		revisions[i] = committed.Revision
		states[i] = objectstore.State{Object: testObject(i+1, "api:v2"), Ownership: states[i].Ownership}
	}

	return keys, revisions, states
}

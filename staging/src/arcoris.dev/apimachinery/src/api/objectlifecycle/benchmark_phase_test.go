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

	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/objectstore"
)

func BenchmarkInitialOwnership(b *testing.B) {
	executor := benchmarkExecutor(b)
	obj := testObject(1, "api:v1")
	prepared, err := executor.prepareObjectRequest(OperationCreate, obj)
	if err != nil {
		b.Fatalf("prepareObjectRequest() error: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := executor.initialOwnership(OperationCreate, prepared.key, owner("creator"), obj.Desired, prepared.resolved); err != nil {
			b.Fatalf("initialOwnership() error: %v", err)
		}
	}
}

func BenchmarkValidateObject(b *testing.B) {
	executor := benchmarkExecutor(b)
	obj := testObject(1, "api:v1")
	prepared, err := executor.prepareObjectRequest(OperationCreate, obj)
	if err != nil {
		b.Fatalf("prepareObjectRequest() error: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := executor.validateObject(OperationCreate, prepared.key, obj, prepared.resolved); err != nil {
			b.Fatalf("validateObject() error: %v", err)
		}
	}
}

func BenchmarkApplyExistingComputeOnly(b *testing.B) {
	executor := benchmarkExecutor(b)
	live := benchmarkCreateObject(b, executor, 1, "api:v1").State
	ownership, err := stateOwnership(OperationApply, objectstoreKeyForBenchmark(b), live.Ownership)
	if err != nil {
		b.Fatalf("stateOwnership() error: %v", err)
	}
	obj := testObject(1, "api:v2")
	prepared, err := executor.prepareObjectRequest(OperationApply, obj)
	if err != nil {
		b.Fatalf("prepareObjectRequest() error: %v", err)
	}
	request := objectapply.Request{
		Owner:     owner("creator"),
		Live:      live.Object,
		Applied:   obj,
		Resource:  prepared.resolved.definition,
		Ownership: ownership,
	}
	options := executor.optionsForApply(false)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := objectapply.Apply(request, options); err != nil {
			b.Fatalf("objectapply.Apply() error: %v", err)
		}
	}
}

func BenchmarkStoreCreateFromLifecycleState(b *testing.B) {
	executor := benchmarkExecutor(b)
	keys, states := benchmarkPreparedStoreCreates(b, executor, b.N)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := executor.store.Create(context.Background(), keys[i], states[i]); err != nil {
			b.Fatalf("store.Create() error: %v", err)
		}
	}
}

func BenchmarkStoreUpdateFromLifecycleState(b *testing.B) {
	executor := benchmarkExecutor(b)
	keys, revisions, states := benchmarkPreparedStoreUpdates(b, executor, b.N)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := executor.store.Update(context.Background(), keys[i], revisions[i], states[i]); err != nil {
			b.Fatalf("store.Update() error: %v", err)
		}
	}
}

func objectstoreKeyForBenchmark(b *testing.B) objectstore.Key {
	b.Helper()

	key, err := objectstore.NewKey(testGVR(), testName(1))
	if err != nil {
		b.Fatalf("objectstore.NewKey() error: %v", err)
	}

	return key
}

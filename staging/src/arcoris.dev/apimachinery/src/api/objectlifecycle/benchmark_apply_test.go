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
	"fmt"
	"sync/atomic"
	"testing"

	"arcoris.dev/apimachinery/api/objectmemorystore"
	"arcoris.dev/apimachinery/api/resourcecatalog"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func BenchmarkExecutorApplyCreate(b *testing.B) {
	executor := benchmarkExecutor(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := executor.Apply(
			context.Background(),
			ApplyRequest{Object: testObject(i+1, "api:v1"), Owner: owner("creator")},
		)
		if err != nil {
			b.Fatalf("Apply() error: %v", err)
		}
	}
}

func BenchmarkExecutorApplyExisting(b *testing.B) {
	executor := benchmarkExecutor(b)
	createObjectForBenchmark(b, executor, 1, "api:v1")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := executor.Apply(
			context.Background(),
			ApplyRequest{Object: testObject(1, fmt.Sprintf("api:%d", i)), Owner: owner("creator")},
		)
		if err != nil {
			b.Fatalf("Apply() error: %v", err)
		}
	}
}

func BenchmarkExecutorCreate(b *testing.B) {
	executor := benchmarkExecutor(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := executor.Create(
			context.Background(),
			CreateRequest{Object: testObject(i+1, "api:v1"), Owner: owner("creator")},
		)
		if err != nil {
			b.Fatalf("Create() error: %v", err)
		}
	}
}

func BenchmarkExecutorDelete(b *testing.B) {
	executor := benchmarkExecutor(b)
	current := createObjectForBenchmark(b, executor, 1, "api:v1")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := executor.Delete(
			context.Background(),
			DeleteRequest{Resource: testGVR(), Object: testName(1), Expected: current.State.Revision},
		)
		if err != nil {
			b.Fatalf("Delete() error: %v", err)
		}
		current = createObjectForBenchmark(b, executor, 1, "api:v1")
	}
}

func BenchmarkExecutorGet(b *testing.B) {
	executor := benchmarkExecutor(b)
	createObjectForBenchmark(b, executor, 1, "api:v1")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := executor.Get(context.Background(), GetRequest{Resource: testGVR(), Object: testName(1)})
		if err != nil {
			b.Fatalf("Get() error: %v", err)
		}
	}
}

func BenchmarkExecutorParallelApplyDistinctObjects(b *testing.B) {
	executor := benchmarkExecutor(b)
	var next atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := int(next.Add(1))
			_, err := executor.Apply(
				context.Background(),
				ApplyRequest{Object: testObject(index, "api:v1"), Owner: owner("creator")},
			)
			if err != nil {
				b.Fatalf("Apply() error: %v", err)
			}
		}
	})
}

func BenchmarkExecutorParallelApplySameObject(b *testing.B) {
	executor := benchmarkExecutor(b)
	createObjectForBenchmark(b, executor, 1, "api:v1")
	var next atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := next.Add(1)
			_, _ = executor.Apply(
				context.Background(),
				ApplyRequest{
					Object: testObject(1, fmt.Sprintf("api:%d", index)),
					Owner:  owner("creator"),
					Force:  true,
				},
			)
		}
	})
}

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

func createObjectForBenchmark(b *testing.B, executor *Executor, index int, image string) Result {
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

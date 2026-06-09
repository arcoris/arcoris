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
	"sync/atomic"
	"testing"
)

func BenchmarkExecutorParallelApplyDistinctObjects(b *testing.B) {
	executor := benchmarkExecutor(b)
	requests := benchmarkApplyRequests(b, b.N)
	var next atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := int(next.Add(1)) - 1
			if _, err := executor.Apply(context.Background(), requests[index]); err != nil {
				b.Fatalf("Apply() error: %v", err)
			}
		}
	})
}

func BenchmarkExecutorParallelApplySameObject(b *testing.B) {
	executor := benchmarkExecutor(b)
	benchmarkCreateObject(b, executor, 1, "api:v1")
	request := ApplyRequest{Object: testObject(1, "api:v2"), Owner: owner("creator"), Force: true}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = executor.Apply(context.Background(), request)
		}
	})
}

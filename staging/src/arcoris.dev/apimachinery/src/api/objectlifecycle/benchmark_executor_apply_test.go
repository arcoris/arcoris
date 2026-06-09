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
)

func BenchmarkExecutorApplyCreatePrepared(b *testing.B) {
	executor := benchmarkExecutor(b)
	requests := benchmarkApplyRequests(b, b.N)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := executor.Apply(context.Background(), requests[i]); err != nil {
			b.Fatalf("Apply() error: %v", err)
		}
	}
}

func BenchmarkExecutorApplyExistingPrepared(b *testing.B) {
	executor := benchmarkExecutor(b)
	benchmarkCreateObject(b, executor, 1, "api:v1")
	request := ApplyRequest{Object: testObject(1, "api:v2"), Owner: owner("creator")}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := executor.Apply(context.Background(), request); err != nil {
			b.Fatalf("Apply() error: %v", err)
		}
	}
}

func BenchmarkExecutorEndToEndApplyCreate(b *testing.B) {
	executor := benchmarkExecutor(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := ApplyRequest{Object: testObject(i+1, "api:v1"), Owner: owner("creator")}
		if _, err := executor.Apply(context.Background(), req); err != nil {
			b.Fatalf("Apply() error: %v", err)
		}
	}
}

func BenchmarkExecutorEndToEndApplyExisting(b *testing.B) {
	executor := benchmarkExecutor(b)
	benchmarkCreateObject(b, executor, 1, "api:v1")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := ApplyRequest{Object: testObject(1, "api:v2"), Owner: owner("creator")}
		if _, err := executor.Apply(context.Background(), req); err != nil {
			b.Fatalf("Apply() error: %v", err)
		}
	}
}

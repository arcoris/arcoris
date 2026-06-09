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

func BenchmarkExecutorCreatePrepared(b *testing.B) {
	executor := benchmarkExecutor(b)
	requests := benchmarkCreateRequests(b, b.N)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := executor.Create(context.Background(), requests[i]); err != nil {
			b.Fatalf("Create() error: %v", err)
		}
	}
}

func BenchmarkExecutorEndToEndCreate(b *testing.B) {
	executor := benchmarkExecutor(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := CreateRequest{Object: testObject(i+1, "api:v1"), Owner: owner("creator")}
		if _, err := executor.Create(context.Background(), req); err != nil {
			b.Fatalf("Create() error: %v", err)
		}
	}
}

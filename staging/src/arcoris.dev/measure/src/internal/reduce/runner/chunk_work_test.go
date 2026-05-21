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


package runner

import (
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func TestChunkWorkersUsesChunkCount(t *testing.T) {
	work := chunkWorkers(10, core.Options{Workers: 8, ChunkSize: 100})
	if work.size != 100 {
		t.Fatalf("chunk size = %d, want 100", work.size)
	}
	if work.count != 1 {
		t.Fatalf("chunk count = %d, want 1", work.count)
	}
	if work.workers != 1 {
		t.Fatalf("workers = %d, want 1", work.workers)
	}
}

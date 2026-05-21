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


package planner

import (
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func TestFixedChunksCoversInputWithFixedChunks(t *testing.T) {
	got := FixedChunks(10, core.Options{ChunkSize: 4}, nil)
	want := []core.Range{{Start: 0, End: 4}, {Start: 4, End: 8}, {Start: 8, End: 10}}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("range[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
	assertPlanCovers(t, got, 10)
}

func TestFixedChunksIgnoresWorkersAndMinItemsPerWorker(t *testing.T) {
	got := FixedChunks(10, core.Options{Workers: 100, MinItemsPerWorker: 1000, ChunkSize: 4}, nil)
	want := []core.Range{{Start: 0, End: 4}, {Start: 4, End: 8}, {Start: 8, End: 10}}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("range[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestFixedChunksHandlesEmptyInput(t *testing.T) {
	got := FixedChunks(0, core.Options{ChunkSize: 4}, []core.Range{{Start: 1, End: 2}})
	if len(got) != 0 {
		t.Fatalf("FixedChunks empty = %#v, want no ranges", got)
	}
}

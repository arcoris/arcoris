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
	"math"
	"testing"
)

func TestChunkRangeClampsLastChunk(t *testing.T) {
	got := chunkRange(10, 4, 2)
	if got.Start != 8 || got.End != 10 {
		t.Fatalf("chunkRange() = [%d,%d), want [8,10)", got.Start, got.End)
	}
}

func TestChunkCountHandlesInvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		chunk int
		want  int
	}{
		{name: "zero items", n: 0, chunk: 4, want: 0},
		{name: "negative items", n: -1, chunk: 4, want: 0},
		{name: "zero chunk", n: 10, chunk: 0, want: 0},
		{name: "negative chunk", n: 10, chunk: -1, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := chunkCount(tt.n, tt.chunk); got != tt.want {
				t.Fatalf(
					"chunkCount(%d, %d) = %d, want %d",
					tt.n,
					tt.chunk,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestChunkCountAvoidsOverflowFormula(t *testing.T) {
	got := chunkCount(math.MaxInt, math.MaxInt)
	if got != 1 {
		t.Fatalf("chunkCount(MaxInt, MaxInt) = %d, want 1", got)
	}
}

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

import "testing"

func TestFixedChunksChunkBlockAssignsContiguousChunks(t *testing.T) {
	want := [][2]int{{0, 4}, {4, 8}, {8, 12}, {12, 16}}
	for worker := range want {
		start, end := fixedChunkBlock(16, 4, worker)
		if start != want[worker][0] || end != want[worker][1] {
			t.Fatalf(
				"worker %d block = [%d,%d), want [%d,%d)",
				worker,
				start,
				end,
				want[worker][0],
				want[worker][1],
			)
		}
	}
}

func TestFixedChunksChunkBlockBalancesRemainder(t *testing.T) {
	want := [][2]int{{0, 4}, {4, 7}, {7, 10}}
	for worker := range want {
		start, end := fixedChunkBlock(10, 3, worker)
		if start != want[worker][0] || end != want[worker][1] {
			t.Fatalf(
				"worker %d block = [%d,%d), want [%d,%d)",
				worker,
				start,
				end,
				want[worker][0],
				want[worker][1],
			)
		}
	}
}

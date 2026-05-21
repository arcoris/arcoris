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

	"arcoris.dev/measure/internal/reduce/core"
)

func TestActiveWorkersHandlesZeroJobs(t *testing.T) {
	if got := activeWorkers(8, 0); got != 0 {
		t.Fatalf("activeWorkers() = %d, want 0", got)
	}
}

func TestActiveWorkersCapsToJobs(t *testing.T) {
	if got := activeWorkers(8, 3); got != 3 {
		t.Fatalf("activeWorkers() = %d, want 3", got)
	}
}

func TestActiveWorkersUsesRequestedWorkersWhenValid(t *testing.T) {
	if got := activeWorkers(4, 10); got != 4 {
		t.Fatalf("activeWorkers() = %d, want 4", got)
	}
}

func TestShouldReduceSequentiallyIgnoresExplicitStrategySequential(t *testing.T) {
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1, Strategy: core.StrategySequential}
	if shouldReduceSequentially(100, opts) {
		t.Fatal(
			"shouldReduceSequentially() = true, want false; dispatch handles StrategySequential explicitly",
		)
	}
}

func TestShouldReduceSequentiallyRespectsWorkerCount(t *testing.T) {
	opts := core.Options{Workers: 1, MinItemsPerWorker: 1, Strategy: core.StrategyBalanced}
	if !shouldReduceSequentially(100, opts) {
		t.Fatal("shouldReduceSequentially() = false, want true")
	}
}

func TestShouldReduceSequentiallyRespectsMinimumItemsPerWorker(t *testing.T) {
	opts := core.Options{Workers: 8, MinItemsPerWorker: 100, Strategy: core.StrategyBalanced}
	if !shouldReduceSequentially(100, opts) {
		t.Fatal("shouldReduceSequentially() = false, want true below threshold")
	}
	opts.MinItemsPerWorker = 50
	if shouldReduceSequentially(100, opts) {
		t.Fatal("shouldReduceSequentially() = true, want false at threshold")
	}
}

func TestEnoughItemsForParallelBoundaries(t *testing.T) {
	tests := []struct {
		name              string
		n                 int
		minItemsPerWorker int
		want              bool
	}{
		{name: "empty", n: 0, minItemsPerWorker: 1, want: false},
		{name: "disabled threshold", n: 1, minItemsPerWorker: 0, want: true},
		{name: "below two workers", n: 99, minItemsPerWorker: 50, want: false},
		{name: "exactly two workers", n: 100, minItemsPerWorker: 50, want: true},
		{
			name:              "large threshold avoids overflow",
			n:                 math.MaxInt,
			minItemsPerWorker: math.MaxInt,
			want:              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := enoughItemsForParallel(tt.n, tt.minItemsPerWorker); got != tt.want {
				t.Fatalf(
					"enoughItemsForParallel(%d, %d) = %v, want %v",
					tt.n,
					tt.minItemsPerWorker,
					got,
					tt.want,
				)
			}
		})
	}
}

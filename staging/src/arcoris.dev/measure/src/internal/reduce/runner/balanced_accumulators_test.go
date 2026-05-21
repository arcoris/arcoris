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

func TestAccumulateBalancedWorkerPartialsSupportsLazyInitialization(t *testing.T) {
	type partial struct {
		Buckets []int
	}

	got, ok := accumulateBalancedWorkerPartials(
		32,
		core.Options{
			Workers:           4,
			MinItemsPerWorker: 1,
			Strategy:          core.StrategyBalanced,
		},
		nil,
		func(_ int, r core.Range, dst *partial) {
			if dst.Buckets == nil {
				dst.Buckets = make([]int, 1)
			}
			dst.Buckets[0] += r.Len()
		},
		func(dst *partial, src partial) {
			if dst.Buckets == nil {
				dst.Buckets = make([]int, len(src.Buckets))
			}
			for i := range src.Buckets {
				dst.Buckets[i] += src.Buckets[i]
			}
		},
	)
	if !ok {
		t.Fatal("accumulateBalancedWorkerPartials returned false for non-empty input")
	}
	if len(got.Buckets) != 1 || got.Buckets[0] != 32 {
		t.Fatalf("got buckets = %#v, want [32]", got.Buckets)
	}
}

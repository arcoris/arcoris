/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package planner

import (
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func TestPlanDispatchesByStrategy(t *testing.T) {
	seq := Plan(10, core.Options{Strategy: core.StrategySequential}, nil)
	if len(seq) != 1 {
		t.Fatalf("sequential len = %d, want 1", len(seq))
	}
	fixed := Plan(10, core.Options{Strategy: core.StrategyFixedChunks, ChunkSize: 3}, nil)
	if len(fixed) != 4 {
		t.Fatalf("fixed len = %d, want 4", len(fixed))
	}
	balanced := Plan(
		1000,
		core.Options{
			Strategy:          core.StrategyBalanced,
			Workers:           2,
			MinItemsPerWorker: 100,
		},
		nil,
	)
	if len(balanced) != 2 {
		t.Fatalf("balanced len = %d, want 2", len(balanced))
	}
	dyn := Plan(10, core.Options{Strategy: core.StrategyDynamicChunks, ChunkSize: 3}, nil)
	if len(dyn) != len(fixed) {
		t.Fatalf("dynamic inspection len = %d, want fixed chunk len %d", len(dyn), len(fixed))
	}
}

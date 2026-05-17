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

	"arcoris.dev/measure/internal/reduce"
)

func TestPlanDispatchesByStrategy(t *testing.T) {
	seq := Plan(10, reduce.Options{Strategy: reduce.StrategySequential}, nil)
	if len(seq) != 1 {
		t.Fatalf("sequential len = %d, want 1", len(seq))
	}
	fixed := Plan(10, reduce.Options{Strategy: reduce.StrategyFixed, ChunkSize: 3}, nil)
	if len(fixed) != 4 {
		t.Fatalf("fixed len = %d, want 4", len(fixed))
	}
	stat := Plan(1000, reduce.Options{Strategy: reduce.StrategyStatic, Workers: 2, MinItemsPerWorker: 100}, nil)
	if len(stat) != 2 {
		t.Fatalf("static len = %d, want 2", len(stat))
	}
}

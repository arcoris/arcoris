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

package runner

import (
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func TestReduceWrapsMapper(t *testing.T) {
	got, ok := Reduce[int](
		10,
		core.Options{Workers: 2, MinItemsPerWorker: 1, Strategy: core.StrategyBalanced},
		func(r core.Range) int { return r.Len() },
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("Reduce returned false for non-empty input")
	}
	if got != 10 {
		t.Fatalf("Reduce() = %d, want 10", got)
	}
}

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

func BenchmarkStatic(b *testing.B) {
	var ranges []core.Range
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1024, Strategy: core.StrategyStatic}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ranges = Static(1_000_000, opts, ranges)
	}
	_ = ranges
}

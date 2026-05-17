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

	"arcoris.dev/measure/internal/reduce"
)

func TestRunnerReusesScratch(t *testing.T) {
	r := New[int](reduce.Options{Workers: 4, MinItemsPerWorker: 10, Strategy: reduce.StrategyStatic})
	for i := 0; i < 3; i++ {
		got, ok := r.DoInto(100, func(rng reduce.Range, dst *int) {
			for x := rng.Start; x < rng.End; x++ {
				*dst += x
			}
		}, func(dst *int, src int) { *dst += src })
		if !ok {
			t.Fatal("expected ok")
		}
		if got != 4950 {
			t.Fatalf("got %d, want 4950", got)
		}
	}
}

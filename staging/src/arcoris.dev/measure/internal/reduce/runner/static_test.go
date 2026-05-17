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

func TestDoIntoProcessesEveryIndexOnce(t *testing.T) {
	result, ok := DoInto[int](
		1000,
		reduce.Options{Workers: 4, MinItemsPerWorker: 100, Strategy: reduce.StrategyStatic},
		nil,
		func(r reduce.Range, dst *int) {
			for i := r.Start; i < r.End; i++ {
				*dst += i
			}
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("expected non-empty reduction")
	}
	want := 999 * 1000 / 2
	if result != want {
		t.Fatalf("result = %d, want %d", result, want)
	}
}

func TestDoIntoSequentialFastPath(t *testing.T) {
	var calls int
	_, ok := DoInto[int](
		10,
		reduce.Options{Workers: 8, MinItemsPerWorker: 100, Strategy: reduce.StrategyStatic},
		nil,
		func(r reduce.Range, dst *int) {
			calls++
			*dst = r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("expected ok")
	}
	if calls != 1 {
		t.Fatalf("mapper calls = %d, want 1", calls)
	}
}

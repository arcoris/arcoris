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
	"sync/atomic"
	"testing"

	"arcoris.dev/measure/internal/reduce"
)

func TestDoDynamicIntoProcessesEveryIndexOnce(t *testing.T) {
	result, ok := DoDynamicInto[int](
		1000,
		reduce.Options{Workers: 4, MinItemsPerWorker: 100, ChunkSize: 17, Strategy: reduce.StrategyDynamic},
		nil,
		func(_ int, r reduce.Range, dst *int) {
			for i := r.Start; i < r.End; i++ {
				*dst += i
			}
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("expected ok")
	}
	want := 999 * 1000 / 2
	if result != want {
		t.Fatalf("result = %d, want %d", result, want)
	}
}

func TestDoDynamicIntoSkipsIdlePartials(t *testing.T) {
	var mergeCalls atomic.Int64
	result, ok := DoDynamicInto[int](
		10,
		reduce.Options{Workers: 8, MinItemsPerWorker: 1, ChunkSize: 100, Strategy: reduce.StrategyDynamic},
		nil,
		func(_ int, r reduce.Range, dst *int) {
			*dst += r.Len()
		},
		func(dst *int, src int) {
			mergeCalls.Add(1)
			*dst += src
		},
	)
	if !ok {
		t.Fatal("expected ok")
	}
	if result != 10 {
		t.Fatalf("result = %d, want 10", result)
	}
	if got := mergeCalls.Load(); got != 0 {
		t.Fatalf("merge calls = %d, want 0 for one active chunk", got)
	}
}

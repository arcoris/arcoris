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

	"arcoris.dev/measure/internal/reduce/core"
)

func TestReduceSequentiallyMapsFullRangeOnce(t *testing.T) {
	var calls atomic.Int64
	got, ok := reduceSequentially[int](23, func(r core.Range, dst *int) {
		calls.Add(1)
		if r.Start != 0 || r.End != 23 {
			t.Errorf("range = [%d,%d), want [0,23)", r.Start, r.End)
		}
		*dst = r.Len()
	})
	if !ok {
		t.Fatal("reduceSequentially returned false for non-empty input")
	}
	if got != 23 {
		t.Fatalf("reduceSequentially() = %d, want 23", got)
	}
	if calls.Load() != 1 {
		t.Fatalf("mapper calls = %d, want 1", calls.Load())
	}
}

func TestReduceSequentiallyIndexedUsesWorkerSlotZero(t *testing.T) {
	got, ok := reduceSequentiallyIndexed[int](5, func(worker int, r core.Range, dst *int) {
		if worker != 0 {
			t.Errorf("worker = %d, want 0", worker)
		}
		*dst = r.Len()
	})
	if !ok {
		t.Fatal("reduceSequentiallyIndexed returned false for non-empty input")
	}
	if got != 5 {
		t.Fatalf("reduceSequentiallyIndexed() = %d, want 5", got)
	}
}

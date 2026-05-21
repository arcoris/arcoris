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

func FuzzFixedChunksCoverage(f *testing.F) {
	f.Add(1000, 64)
	f.Add(1, 64)
	f.Add(0, 64)
	f.Fuzz(func(t *testing.T, n int, chunk int) {
		if n < 0 || n > 1_000_000 {
			return
		}
		plan := FixedChunks(n, core.Options{ChunkSize: chunk}, nil)
		pos := 0
		for i, r := range plan {
			if r.Empty() {
				t.Fatalf("range[%d] empty: %#v", i, r)
			}
			if r.Start != pos {
				t.Fatalf("range[%d].Start=%d want %d plan=%#v", i, r.Start, pos, plan)
			}
			pos = r.End
		}
		if pos != n {
			t.Fatalf("covered %d, want %d plan=%#v", pos, n, plan)
		}
	})
}

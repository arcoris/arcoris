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

import "arcoris.dev/measure/internal/reduce/core"

// FixedChunks returns stable fixed-size chunks covering [0:n).
//
// FixedChunks ignores Workers and MinItemsPerWorker. It is a planning primitive for
// fixed-grain execution and tuning; runners decide how many workers consume the
// resulting chunks.
func FixedChunks(n int, opts core.Options, dst []core.Range) []core.Range {
	dst = dst[:0]
	if n <= 0 {
		return dst
	}

	opts = core.NormalizeOptions(opts)
	chunk := opts.ChunkSize
	if chunk <= 0 {
		chunk = core.DefaultChunkSize
	}
	chunks := 1 + (n-1)/chunk
	if cap(dst) < chunks {
		dst = make([]core.Range, 0, chunks)
	}

	for start := 0; start < n; {
		end := start + chunk
		if end > n {
			end = n
		}
		dst = append(dst, core.Range{Start: start, End: end})
		start = end
	}
	return dst
}

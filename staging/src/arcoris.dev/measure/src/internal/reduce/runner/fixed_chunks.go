// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package runner

// fixedChunkBlock returns the half-open chunk-index block assigned to worker.
//
// FixedChunks uses contiguous chunk blocks to keep ownership deterministic and
// locality-friendly. The distribution is balanced by at most one chunk.
func fixedChunkBlock(
	chunks int,
	workers int,
	worker int,
) (int, int) {
	if chunks <= 0 || workers <= 0 || worker < 0 || worker >= workers {
		return 0, 0
	}
	base := chunks / workers
	rem := chunks % workers
	start := worker*base + min(worker, rem)
	end := start + base
	if worker < rem {
		end++
	}
	return start, end
}

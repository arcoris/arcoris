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

import "arcoris.dev/measure/internal/reduce/core"

// chunkWork groups the normalized scheduling inputs for chunk-based execution.
type chunkWork struct {
	// size is the validated chunk size that runners should use.
	size int

	// count is the number of chunks needed to cover the input.
	count int

	// workers is the bounded worker count derived from chunk count.
	workers int
}

// chunkWorkers normalizes chunk execution settings for fixed and dynamic paths.
//
// Chunk-based execution decides worker count from chunk count, not from raw item
// count, because only chunks are scheduled to workers.
func chunkWorkers(n int, opts core.Options) chunkWork {
	size := opts.ChunkSize
	if size <= 0 {
		size = core.DefaultChunkSize
	}
	count := chunkCount(n, size)
	return chunkWork{
		size:    size,
		count:   count,
		workers: activeWorkers(opts.Workers, count),
	}
}

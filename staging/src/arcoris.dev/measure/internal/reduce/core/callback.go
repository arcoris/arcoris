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

package core

// Mapper computes a complete partial result for one range or chunk.
//
// Mapper is part of the Reduce family: each call owns one fresh mapping
// operation and returns the full partial for that operation. The runner may
// merge the produced partial later. Keep hot per-element loops inside the
// mapper so runners do not add callback overhead. Runners may invoke Mapper
// concurrently, so captured shared state must be immutable or synchronized by
// the caller.
type Mapper[T any] func(Range) T

// IntoMapper writes a complete partial for one range or chunk into dst.
//
// IntoMapper is part of the Reduce family. The destination is fresh for the
// current mapping operation, so implementations may assign the result directly
// and are not required to preserve previous dst contents. The runner may merge
// the produced partial later. This callback is generic and safe, but it is not
// the maximum-performance worker-local accumulation contract.
type IntoMapper[T any] func(Range, *T)

// IndexedIntoMapper writes a complete partial while exposing the execution slot.
//
// IndexedIntoMapper is the indexed Reduce-family form. The destination is fresh
// for the current range or chunk, so implementations may assign it. Balanced
// range-local runners pass the worker currently processing the planned range.
// FixedChunks and DynamicChunks runners pass the worker that owns or claimed the
// current chunk. The worker value is stable only within one reduction call.
type IndexedIntoMapper[T any] func(worker int, r Range, dst *T)

// Accumulator adds one range or chunk directly into a worker-local partial.
//
// Accumulator is part of the Accumulate family. The destination belongs to a
// worker-local partial, and a runner may call the accumulator multiple times
// with the same dst. Implementations must accumulate into dst and must not
// blindly overwrite previous state unless replacement is the intended
// accumulation behavior. Captured shared state must be immutable or synchronized
// by the caller.
type Accumulator[T any] func(r Range, dst *T)

// IndexedAccumulator adds one range or chunk into a worker-local partial while
// exposing the execution slot.
//
// IndexedAccumulator has the same ownership and accumulation contract as
// Accumulator. The worker value identifies the worker-local partial being
// updated during the current reduction call.
type IndexedAccumulator[T any] func(worker int, r Range, dst *T)

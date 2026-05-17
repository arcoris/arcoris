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

// Mapper computes a complete partial result for one planned range.
//
// Mapper is called once per range, not once per element. Keep hot per-element
// loops inside the mapper so runners do not add callback overhead. Runners may
// invoke Mapper concurrently, so any captured state must be immutable or
// synchronized by the caller.
type Mapper[T any] func(Range) T

// IntoMapper maps one planned range into dst.
//
// Runners pass a zero-value destination for each planned range. Implementations
// may assign fields or accumulate into dst, whichever is faster for the partial
// type. The destination is owned by the current mapper call; captured state is
// still the caller's responsibility to protect when reductions run in parallel.
type IntoMapper[T any] func(Range, *T)

// IndexedIntoMapper maps a range while exposing the execution slot.
//
// Static range-local runners pass the worker slot currently processing the
// planned range. Fixed and dynamic worker-local runners pass the worker slot
// that claimed the current chunk; each chunk receives a fresh destination and
// the runner folds chunk partials into that worker's accumulator with Merger.
// The worker value is stable only within one reduction call.
type IndexedIntoMapper[T any] func(worker int, r Range, dst *T)

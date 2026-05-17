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

import "arcoris.dev/measure/internal/reduce/core"

// Runner owns normalized options and scratch storage for repeated reductions.
//
// Runner is not safe for concurrent use because each call mutates its reusable
// scratch buffers. Use one Runner per concurrent caller.
type Runner[T any] struct {
	// opts is normalized once during construction to avoid repeated defaulting.
	opts core.Options

	// scratch stores reusable ranges, partial-result slots, and active-worker
	// flags between calls.
	scratch core.Scratch[T]
}

// New returns a Runner with normalized options and empty scratch storage.
func New[T any](opts core.Options) Runner[T] {
	return Runner[T]{opts: core.NormalizeOptions(opts)}
}

// ReduceInto executes a reduction using runner-owned scratch buffers.
func (r *Runner[T]) ReduceInto(
	n int,
	mapRange core.IntoMapper[T],
	mergeFn core.Merger[T],
) (T, bool) {
	return ReduceInto(
		n,
		r.opts,
		&r.scratch,
		mapRange,
		mergeFn,
	)
}

// ReduceIndexedInto executes an indexed reduction using runner-owned scratch
// buffers.
func (r *Runner[T]) ReduceIndexedInto(
	n int,
	mapRange core.IndexedIntoMapper[T],
	mergeFn core.Merger[T],
) (T, bool) {
	return ReduceIndexedInto(
		n,
		r.opts,
		&r.scratch,
		mapRange,
		mergeFn,
	)
}

// AccumulateInto executes a worker-local accumulation reduction using
// runner-owned scratch buffers.
func (r *Runner[T]) AccumulateInto(
	n int,
	accumulate core.Accumulator[T],
	mergeFn core.Merger[T],
) (T, bool) {
	return AccumulateInto(
		n,
		r.opts,
		&r.scratch,
		accumulate,
		mergeFn,
	)
}

// AccumulateIndexedInto executes an indexed worker-local accumulation reduction
// using runner-owned scratch buffers.
func (r *Runner[T]) AccumulateIndexedInto(
	n int,
	accumulate core.IndexedAccumulator[T],
	mergeFn core.Merger[T],
) (T, bool) {
	return AccumulateIndexedInto(
		n,
		r.opts,
		&r.scratch,
		accumulate,
		mergeFn,
	)
}

// Reset clears runner-owned scratch contents while retaining backing storage.
func (r *Runner[T]) Reset() { r.scratch.Reset() }

// Clear zeroes runner-owned scratch storage while retaining capacity.
func (r *Runner[T]) Clear() { r.scratch.Clear() }

// Release drops runner-owned scratch backing storage.
func (r *Runner[T]) Release() { r.scratch.Release() }

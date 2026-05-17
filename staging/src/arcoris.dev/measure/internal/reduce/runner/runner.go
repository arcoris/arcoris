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

import "arcoris.dev/measure/internal/reduce"

// Runner owns normalized options and scratch storage for repeated reductions.
//
// Runner is not safe for concurrent use because each call mutates its reusable
// scratch buffers. Use one Runner per concurrent caller.
type Runner[T any] struct {
	// opts is normalized once during construction to avoid repeated defaulting.
	opts reduce.Options

	// scratch stores reusable ranges and partial-result slots between calls.
	scratch reduce.Scratch[T]
}

// New returns a Runner with normalized options and empty scratch storage.
func New[T any](opts reduce.Options) Runner[T] {
	return Runner[T]{opts: reduce.NormalizeOptions(opts)}
}

// DoInto executes a reduction using runner-owned scratch buffers.
func (r *Runner[T]) DoInto(n int, mapRange reduce.IntoMapper[T], mergeFn reduce.Merger[T]) (T, bool) {
	return DoInto(n, r.opts, &r.scratch, mapRange, mergeFn)
}

// DoIndexedInto executes an indexed reduction using runner-owned scratch
// buffers.
func (r *Runner[T]) DoIndexedInto(n int, mapRange reduce.IndexedIntoMapper[T], mergeFn reduce.Merger[T]) (T, bool) {
	return DoIndexedInto(n, r.opts, &r.scratch, mapRange, mergeFn)
}

// Reset clears runner-owned scratch contents while retaining backing storage.
func (r *Runner[T]) Reset() { r.scratch.Reset() }

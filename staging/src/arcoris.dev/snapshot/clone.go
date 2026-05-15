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

package snapshot

// CloneFunc copies a value across a snapshot ownership boundary.
//
// CloneFunc is explicit because generic Go code cannot know how to safely copy
// arbitrary values. Structs containing slices, maps, pointers, or mutable nested
// state usually need a domain-specific clone function. Immutable or copy-safe
// values may use Identity.
type CloneFunc[T any] func(T) T

// Identity returns value unchanged.
//
// Identity is appropriate only for immutable or copy-safe values. It must not be
// used for values whose maps, slices, pointers, buffers, or other mutable fields
// could be mutated by readers after a snapshot is returned.
func Identity[T any](val T) T {
	return val
}

// requireClone returns clone or panics when clone is nil.
//
// Store calls this during construction so nil clone functions fail early instead
// of silently exposing mutable internal state.
func requireClone[T any](clone CloneFunc[T]) CloneFunc[T] {
	if clone == nil {
		panic("snapshot: nil clone function")
	}

	return clone
}

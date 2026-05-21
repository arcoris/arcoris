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

package liveconfig

// CloneFunc copies a configuration value across the Holder ownership boundary.
//
// CloneFunc is required when T contains slices, maps, pointers, buffers, or any
// mutable nested state that callers or readers could mutate after Apply or
// Snapshot. Immutable and copy-safe values may use the package default identity
// clone. The clone runs before normalization, validation, equality checks, and
// publication.
//
// A clone function should return a value that is fully owned by the holder. For
// nested mutable state, that usually means copying every map, slice, and pointer
// target that remains reachable from the published configuration.
type CloneFunc[T any] func(T) T

// identityClone returns val unchanged.
//
// The identity clone is correct only for immutable or copy-safe configuration
// values. Callers that publish mutable aggregate values must provide WithClone
// or guarantee by other means that nobody mutates the value after publication.
func identityClone[T any](val T) T {
	return val
}

// cloneValue applies cfg's configured ownership-boundary clone.
//
// The helper exists so the candidate pipeline names each phase explicitly and
// tests can exercise clone behavior without duplicating Holder setup.
func cloneValue[T any](cfg config[T], val T) T {
	return cfg.clone(val)
}

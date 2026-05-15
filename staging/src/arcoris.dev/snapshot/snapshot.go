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

// Snapshot is a lightweight point-in-time value read from a source.
//
// Snapshot carries the source-local Revision at which Value was observed. The
// type itself does not perform cloning, copying, or immutability enforcement.
// Those guarantees come from the source that produced the snapshot.
//
// Store returns cloned values and is safe for mutable payloads when the supplied
// CloneFunc is correct. Publisher returns immutable published values and relies
// on callers not mutating values after Publish.
type Snapshot[T any] struct {
	// Revision is the source-local revision of Value.
	Revision Revision

	// Value is the typed value observed at Revision.
	Value T
}

// IsZeroRevision reports whether the snapshot has ZeroRevision.
func (s Snapshot[T]) IsZeroRevision() bool {
	return s.Revision.IsZero()
}

// ChangedSince reports whether the snapshot revision differs from revision.
func (s Snapshot[T]) ChangedSince(rev Revision) bool {
	return s.Revision.ChangedSince(rev)
}

// WithValue returns a snapshot with the same Revision and a different Value.
//
// WithValue is useful for adapting a read model while preserving the revision
// that the derived value came from. The method does not clone either value.
func (s Snapshot[T]) WithValue(val T) Snapshot[T] {
	return Snapshot[T]{
		Revision: s.Revision,
		Value:    val,
	}
}

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

package snapshot

// Snapshot is a lightweight point-in-time value read from one source.
//
// R is the authoritative revision type for that source. It may be LocalRevision
// for Store and Publisher, or a domain-owned revision such as an object-store
// revision. Snapshot only requires comparability so callers can detect whether a
// source has changed since a previously observed revision.
//
// Snapshot does not clone, freeze, or otherwise protect Value. Ownership and
// immutability are source responsibilities.
type Snapshot[R comparable, T any] struct {
	// Revision is the source or domain revision at which Value was observed.
	Revision R

	// Value is the typed read-model value observed at Revision.
	Value T
}

// ChangedSince reports whether the snapshot revision differs from prev.
//
// The check is equality-based. Snapshot does not assume R is ordered, dense, or
// globally comparable across unrelated sources.
func (s Snapshot[R, T]) ChangedSince(prev R) bool {
	return s.Revision != prev
}

// WithValue returns a snapshot with the same Revision and a different Value.
//
// WithValue is useful for deriving a smaller or transformed read model while
// preserving the revision that made the value visible. The method does not clone
// either value.
func (s Snapshot[R, T]) WithValue(value T) Snapshot[R, T] {
	return Snapshot[R, T]{
		Revision: s.Revision,
		Value:    value,
	}
}

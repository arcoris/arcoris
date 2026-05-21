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

// Source is a read-only provider of lightweight snapshots.
//
// Source is a consumer-facing interface. Functions that only need to read the
// current state should accept Source[T] rather than depending on a concrete Store
// or Publisher implementation.
type Source[T any] interface {
	// Snapshot returns the source's current lightweight snapshot.
	Snapshot() Snapshot[T]
}

// StampedSource is a read-only provider of stamped snapshots.
//
// StampedSource should be used when consumers need local update time metadata in
// addition to the value and revision.
type StampedSource[T any] interface {
	// Stamped returns the source's current stamped snapshot.
	Stamped() Stamped[T]
}

// RevisionSource exposes the current source-local revision without requiring the
// caller to read the value.
//
// RevisionSource is useful for cheap change checks. It does not imply a global
// ordering across independent sources.
type RevisionSource interface {
	// Revision returns the latest committed or published revision known to the
	// source.
	Revision() Revision
}

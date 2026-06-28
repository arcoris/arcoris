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

// Source is an always-available provider of lightweight snapshots.
//
// Source should be used for components such as Store and Publisher that can
// return a snapshot immediately after construction. Sources that may not be
// ready should implement SnapshotReader instead.
type Source[R comparable, T any] interface {
	// Snapshot returns the source's current lightweight snapshot.
	Snapshot() Snapshot[R, T]
}

// StampedSource is an always-available provider of stamped snapshots.
//
// StampedSource should be used when consumers need local update time metadata in
// addition to the value and revision. Sources that may not be ready should
// implement StampedReader instead.
type StampedSource[R comparable, T any] interface {
	// Stamped returns the source's current stamped snapshot.
	Stamped() Stamped[R, T]
}

// RevisionSource exposes the current revision without requiring the caller to
// read a value.
//
// RevisionSource is useful for cheap change checks. It does not imply a global
// ordering across independent sources.
type RevisionSource[R comparable] interface {
	// Revision returns the latest committed or published revision known to the
	// source.
	Revision() R
}

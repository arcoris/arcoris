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

// Reader is a fallible provider of lightweight snapshots.
//
// Reader should be used for sources that are not always ready, for example a
// read model that cannot produce a meaningful snapshot until an initial
// replacement has been applied.
type Reader[R comparable, T any] interface {
	// ReadSnapshot returns the source's current lightweight snapshot or an error
	// explaining why no snapshot is currently available.
	ReadSnapshot() (Snapshot[R, T], error)
}

// StampedReader is a fallible provider of stamped snapshots.
type StampedReader[R comparable, T any] interface {
	// ReadStamped returns the source's current stamped snapshot or an error
	// explaining why no stamped snapshot is currently available.
	ReadStamped() (Stamped[R, T], error)
}

// RevisionReader is a fallible provider of current revisions.
type RevisionReader[R comparable] interface {
	// ReadRevision returns the source's current revision or an error explaining
	// why no revision is currently available.
	ReadRevision() (R, error)
}

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

package objectcache

import "arcoris.dev/apimachinery/api/objectstore"

// objectVersion records what this cache observed for one key at a collection
// boundary. Revision is the cache observation revision, not necessarily the
// committed State revision retained in State.
type objectVersion struct {
	// Revision is the collection/cache observation revision at which this
	// version became visible to the cache.
	Revision objectstore.Revision
	// State is set only for live versions. Tombstones intentionally retain no
	// object state because deletion state is represented by Live=false.
	State objectstore.State
	// Live distinguishes retained object versions from retained delete
	// tombstones.
	Live bool
}

// liveVersion returns a detached live history entry.
//
// The observation revision is supplied by Replace or ApplyChange. It may differ
// from state.Revision when a replacement observes an older object state at a
// newer collection boundary.
func liveVersion(revision objectstore.Revision, state objectstore.State) objectVersion {
	return objectVersion{Revision: revision, State: state.Clone(), Live: true}
}

// tombstoneVersion returns a deletion marker at a cache observation revision.
//
// Tombstones let GetAt prove absence after a delete while keeping deleted
// objects absent from the latest collection.
func tombstoneVersion(revision objectstore.Revision) objectVersion {
	return objectVersion{Revision: revision}
}

// clone detaches the retained live state before a version crosses a cache
// boundary. Tombstones intentionally carry no State.
func (v objectVersion) clone() objectVersion {
	if !v.Live {
		return objectVersion{Revision: v.Revision}
	}

	return objectVersion{Revision: v.Revision, State: v.State.Clone(), Live: true}
}

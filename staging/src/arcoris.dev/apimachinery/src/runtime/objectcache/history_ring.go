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

// versionRing is a fixed-capacity per-object history buffer.
//
// The ring is not internally synchronized. Cache owns synchronization and
// mutates records only while holding Cache.mu.
type versionRing struct {
	// entries stores ring slots. len(entries) may be lower than cap(entries)
	// until the ring first fills.
	entries []objectVersion
	// start is the index of the oldest retained entry when len > 0.
	start int
	// len is the number of valid retained entries in entries.
	len int
}

// newVersionRing creates a fixed-capacity ring. Non-positive capacity produces
// a disabled ring that silently ignores appends.
func newVersionRing(capacity int) versionRing {
	if capacity <= 0 {
		return versionRing{}
	}

	return versionRing{entries: make([]objectVersion, 0, capacity)}
}

// append records version as the newest entry and evicts only this object's
// oldest version when the per-object capacity is full. The overwritten slot is
// cleared before replacement so evicted object graphs do not stay reachable
// through the backing array.
func (r *versionRing) append(version objectVersion) {
	capacity := cap(r.entries)
	if capacity == 0 {
		return
	}

	version = version.clone()
	if r.len < capacity {
		r.entries = append(r.entries, version)
		r.len++
		return
	}

	r.entries[r.start] = objectVersion{}
	r.entries[r.start] = version
	r.start = (r.start + 1) % capacity
}

// newestToOldest visits retained versions from newest to oldest until fn
// returns false.
func (r versionRing) newestToOldest(fn func(objectVersion) bool) {
	for offset := r.len - 1; offset >= 0; offset-- {
		if !fn(r.entry(offset)) {
			return
		}
	}
}

// oldestToNewest visits retained versions from oldest to newest until fn
// returns false. Tests use this to prove eviction order.
func (r versionRing) oldestToNewest(fn func(objectVersion) bool) {
	for offset := 0; offset < r.len; offset++ {
		if !fn(r.entry(offset)) {
			return
		}
	}
}

func (r versionRing) entry(offset int) objectVersion {
	if r.len == 0 {
		return objectVersion{}
	}
	return r.entries[(r.start+offset)%cap(r.entries)].clone()
}

// count reports retained entries. It exists for invariant tests, not as public
// cache API.
func (r versionRing) count() int {
	return r.len
}

// capacity reports the per-object retention bound carried by this ring.
func (r versionRing) capacity() int {
	return cap(r.entries)
}

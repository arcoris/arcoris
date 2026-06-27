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

// objectRecord owns retained history for one object key.
type objectRecord struct {
	// key is retained for diagnostics and invariant tests. The records map also
	// keys by this value; storing it here keeps the record self-describing.
	key      objectstore.Key
	versions versionRing
}

// newObjectRecord creates an empty bounded history record for key.
func newObjectRecord(key objectstore.Key, retained int) *objectRecord {
	return &objectRecord{key: key, versions: newVersionRing(retained)}
}

// append records one observed version. A nil receiver is a no-op so optional
// history paths can stay simple when history is disabled.
func (r *objectRecord) append(version objectVersion) {
	if r == nil {
		return
	}
	r.versions.append(version)
}

// newestToOldest scans retained history from the most recent observation back.
// A nil receiver has no retained proof and therefore visits nothing.
func (r *objectRecord) newestToOldest(fn func(objectVersion) bool) {
	if r == nil {
		return
	}
	r.versions.newestToOldest(fn)
}

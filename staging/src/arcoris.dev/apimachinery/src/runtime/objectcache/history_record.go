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
	key      objectstore.Key
	versions versionRing
}

func newObjectRecord(key objectstore.Key, retained int) *objectRecord {
	return &objectRecord{key: key, versions: newVersionRing(retained)}
}

func (r *objectRecord) append(version objectVersion) {
	if r == nil {
		return
	}
	r.versions.append(version)
}

func (r *objectRecord) newestToOldest(fn func(objectVersion) bool) {
	if r == nil {
		return
	}
	r.versions.newestToOldest(fn)
}

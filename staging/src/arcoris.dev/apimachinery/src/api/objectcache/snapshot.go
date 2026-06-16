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

// Snapshot is an immutable materialized object collection at one observed store
// revision watermark.
//
// Snapshot owns detached internal list items. It does not observe later store
// changes and does not validate resource descriptors, scopes, or query/resource
// consistency.
type Snapshot struct {
	// col is immutable after construction and stores only detached items.
	col collection
}

// IsZero reports whether s has no items and no observed revision.
func (s Snapshot) IsZero() bool {
	return s.col.isZero()
}

// Len returns the number of cached items.
func (s Snapshot) Len() int {
	return s.col.len()
}

// Revision returns the observed store revision watermark captured by s.
func (s Snapshot) Revision() objectstore.Revision {
	return s.col.revision
}

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

import (
	"sync"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Cache is a concurrency-safe materialized read model for one object collection.
//
// Cache is inactive until Replace installs the first collection read. Revision
// zero is a valid ready boundary for an empty initial collection, so readiness
// is tracked separately from the revision watermark.
//
// The latest collection is the source of truth for current live objects.
// Optional object records retain bounded per-key history for GetAt and
// PreviousLive. Cache contains a mutex and must not be copied after
// construction.
type Cache struct {
	// mu protects readiness, latest state, and per-object history records.
	mu sync.RWMutex
	// collection is the structural collection this cache owns.
	collection objectstore.ListRequest
	// ready records whether an initial Replace succeeded.
	ready bool
	// latest is the current materialized collection. It is meaningful only when
	// ready is true.
	latest collection
	// historyEnabled records whether historical reads are supported.
	historyEnabled bool
	// retainedVersionsPerObject bounds every object record independently.
	retainedVersionsPerObject int
	// records stores optional per-object version rings. Deleted keys may remain
	// here as tombstones even when absent from latest.
	records map[objectstore.Key]*objectRecord
}

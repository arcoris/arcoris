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
	"arcoris.dev/apimachinery/runtime/objectreflector"
)

var _ objectreflector.Sink = (*Cache)(nil)

// Cache is a concurrency-safe materialized read model for one object collection.
//
// Cache is inactive until Replace installs the first collection read. Revision
// zero is a valid ready boundary for an empty initial collection, so readiness
// is tracked separately from the revision watermark. Cache contains a mutex and
// must not be copied after construction.
type Cache struct {
	// mu protects collection, ready, and col.
	mu sync.RWMutex
	// collection is the structural collection this cache owns.
	collection objectstore.ListRequest
	// ready records whether an initial Replace succeeded.
	ready bool
	// col is the current materialized collection. It is meaningful only when
	// ready is true.
	col collection
}

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

// View is a detached latest collection read model at one objectstore revision.
//
// View is created by Cache.ReadSnapshot. It owns a detached copy of the cache's
// latest live collection, does not take Cache locks, and never observes future
// Cache mutations. View is safe for concurrent read use after creation.
type View struct {
	// collection is the structural collection this view represents.
	collection objectstore.ListRequest

	// latest is the detached live collection visible at collection revision.
	latest collection
}

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

// GetResult is the result of a latest or historical object read.
type GetResult struct {
	// Key is the requested object key.
	Key objectstore.Key
	// State is populated only when Found is true. It is always detached from
	// cache-owned state.
	State objectstore.State
	// Found reports whether the object was live at the served cache revision.
	Found bool
	// Revision is the cache observation revision served by the read.
	//
	// Get returns the current cache collection revision. GetAt returns the
	// requested cache observation revision. PreviousLive returns the retained
	// live version's observation revision.
	Revision objectstore.Revision
}

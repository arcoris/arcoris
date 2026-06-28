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
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/snapshot"
)

// ReadSnapshot returns a detached latest collection view at one objectstore revision.
//
// Cache is not always ready: before the first successful Replace there is no
// collection boundary to publish. For that reason Cache implements
// snapshot.SnapshotReader instead of the always-available snapshot.Source
// interface.
func (c *Cache) ReadSnapshot() (snapshot.Snapshot[objectstore.Revision, View], error) {
	if c == nil {
		return snapshot.Snapshot[objectstore.Revision, View]{}, ErrInvalidCache
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.ready {
		return snapshot.Snapshot[objectstore.Revision, View]{}, ErrNotReady
	}

	view := View{
		collection: c.collection,
		latest:     c.latest.clone(),
	}

	return snapshot.Snapshot[objectstore.Revision, View]{
		Revision: c.latest.revision,
		Value:    view,
	}, nil
}

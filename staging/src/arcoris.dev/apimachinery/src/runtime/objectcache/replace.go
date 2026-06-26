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
	"context"
	"errors"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

// Replace atomically installs a complete collection read.
//
// Replace validates and builds detached replacement state before taking the
// write lock. The lock is held only for the stale-read check and pointer swap,
// so readers never observe a partially rebuilt collection.
func (c *Cache) Replace(ctx context.Context, read storewatchapi.CollectionRead) error {
	if c == nil {
		return ErrInvalidCache
	}
	if ctx == nil {
		// The reflector always supplies a context, but accepting nil keeps Sink
		// usage from tests and small tools forgiving.
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := read.Validate(); err != nil {
		return errorWith(ErrInvalidRead, err)
	}
	if read.Collection() != c.collection {
		return errorWith(ErrInvalidRead, ErrOutsideCollection)
	}

	result := read.Result()
	next, err := buildCollection(result)
	if err != nil {
		return err
	}
	if err := objectstore.ValidateListResult(c.collection, next.listResult()); err != nil {
		return errorWith(ErrInvalidRead, err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ready && read.Revision().Before(c.latest.revision) {
		return errorWith(ErrStaleRead, errors.New("collection read revision is before cache revision"))
	}

	c.latest = next
	c.resetHistoryLocked(next)
	c.ready = true

	return nil
}

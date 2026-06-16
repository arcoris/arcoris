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
	"fmt"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Replace atomically replaces the entire cache collection.
//
// If result is invalid or older than the current cache revision, the existing
// cache state is left unchanged. Equal revisions are accepted as idempotent
// refreshes. Replacement item order becomes the new cache order.
func (c *Cache) Replace(result objectstore.ListResult) error {
	if c == nil {
		return fmt.Errorf("%w: nil cache", ErrInvalidCache)
	}

	next, err := buildCollection(result, ErrInvalidCache)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.col.isZero() && next.revision.Before(c.col.revision) {
		return staleSnapshotError(c.col.revision, next.revision)
	}

	c.col = next
	return nil
}

func staleSnapshotError(current objectstore.Revision, incoming objectstore.Revision) error {
	return fmt.Errorf(
		"%w: current revision %s, incoming revision %s",
		ErrStaleSnapshot,
		current.String(),
		incoming.String(),
	)
}

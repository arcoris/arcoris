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
)

// ApplyChange atomically applies one committed change after the current cache boundary.
//
// Shape and collection validation happen before taking the write lock. The
// lock protects readiness, monotonic revision checks, per-key preconditions,
// and the final mutation, so failed changes leave the cache untouched.
func (c *Cache) ApplyChange(ctx context.Context, change objectstore.Change) error {
	if c == nil {
		return ErrInvalidCache
	}
	if ctx == nil {
		// Reflector calls always provide a context. Treating nil as background
		// keeps direct Sink use deterministic instead of panicking mid-pipeline.
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := objectstore.ValidateChange(change); err != nil {
		return errorWith(ErrInvalidChange, err)
	}
	if !objectstore.ChangeMatchesListRequest(change, c.collection) {
		return errorWith(ErrInvalidChange, ErrOutsideCollection)
	}

	change = change.Clone()

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.ready {
		return ErrNotReady
	}
	if !c.latest.revision.Before(change.Revision) {
		return errorWith(ErrStaleChange, errors.New("change revision does not advance cache revision"))
	}
	if err := c.latest.validateApply(change); err != nil {
		return err
	}

	c.latest.applyValidated(change)
	c.appendHistoryLocked(change)

	return nil
}

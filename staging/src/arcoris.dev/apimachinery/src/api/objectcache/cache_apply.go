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
	"errors"
	"fmt"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Apply atomically applies a committed object store change to c.
//
// Changes must validate through objectstore, be newer than the current cache
// revision, and be consistent with the current cached key state. Apply assumes
// callers provide a complete ordered change stream for this cached collection.
// It detects stale revisions and per-key before-revision mismatches, but missed
// unrelated-key changes require a future watch/reflector recovery layer.
func (c *Cache) Apply(change objectstore.Change) error {
	if c == nil {
		return fmt.Errorf("%w: nil cache", ErrInvalidCache)
	}
	if err := change.Validate(); err != nil {
		return invalidChangeError(err)
	}
	change = change.Clone()

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.col.revision.Before(change.Revision) {
		return staleChangeError(c.col.revision, change.Revision)
	}
	if err := c.col.validateApply(change); err != nil {
		return err
	}

	c.col.applyValidated(change)
	return nil
}

// invalidChangeError keeps ErrInvalidChange visible while preserving the
// objectstore.Change validation cause for errors.Is callers.
func invalidChangeError(cause error) error {
	if cause == nil {
		return ErrInvalidChange
	}

	return errors.Join(ErrInvalidChange, cause)
}

// staleChangeError reports a monotonic revision violation without converting it
// into an invalid-change error. A stale change can be structurally valid while
// still being unusable for this cache revision.
func staleChangeError(current objectstore.Revision, incoming objectstore.Revision) error {
	return fmt.Errorf(
		"%w: current revision %s, incoming revision %s",
		ErrStaleChange,
		current.String(),
		incoming.String(),
	)
}

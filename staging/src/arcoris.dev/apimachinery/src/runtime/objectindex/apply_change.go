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

package objectindex

import (
	"context"

	"arcoris.dev/apimachinery/api/objectstore"
)

// ApplyChange incrementally updates memberships for one committed change.
//
// Extraction is completed before state mutation. If extraction fails, the
// previous index state is preserved.
func (i *Index) ApplyChange(ctx context.Context, change objectstore.Change) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := i.validateStatic(); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := change.Validate(); err != nil {
		return err
	}

	afterValues, err := i.extractChangeValues(ctx, change)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	if err := i.validateStateLocked(); err != nil {
		return err
	}

	if change.Kind == objectstore.ChangeCreated || change.Kind == objectstore.ChangeUpdated {
		i.state.updateObject(change.Key, afterValues)
	} else {
		i.state.removeObject(change.Key)
	}

	return nil
}

// extractChangeValues validates all extractor-visible states for the change
// before any mutation. Updated and deleted changes extract Before as well as
// After so extractor failures remain all-or-nothing even for old membership
// removal.
func (i *Index) extractChangeValues(ctx context.Context, change objectstore.Change) (map[Name]valueSet, error) {
	switch change.Kind {
	case objectstore.ChangeCreated:
		return i.extractStateValues(ctx, change.Key, change.After)
	case objectstore.ChangeUpdated:
		if _, err := i.extractStateValues(ctx, change.Key, change.Before); err != nil {
			return nil, err
		}

		return i.extractStateValues(ctx, change.Key, change.After)
	case objectstore.ChangeDeleted:
		if _, err := i.extractStateValues(ctx, change.Key, change.Before); err != nil {
			return nil, err
		}

		return nil, nil
	default:
		return nil, ErrInvalidIndex
	}
}

// extractStateValues adapts a committed object state into the ListItem shape
// used by extractors. The state is cloned before validation/extraction so
// extractor code cannot retain mutable index-owned data.
func (i *Index) extractStateValues(ctx context.Context, key objectstore.Key, state objectstore.State) (map[Name]valueSet, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	item := objectstore.ListItem{Key: key, State: state.Clone()}
	if err := objectstore.ValidateListItem(item); err != nil {
		return nil, err
	}

	return extractValues(i.names, i.definitions, item)
}

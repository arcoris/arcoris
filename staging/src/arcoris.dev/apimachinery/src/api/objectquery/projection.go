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

package objectquery

import "arcoris.dev/apimachinery/api/objectstore"

// ChangeProjectionKind classifies a committed change against a predicate.
type ChangeProjectionKind uint8

// Change projection outcomes.
const (
	// ChangeProjectionIgnored means neither the previous nor final state
	// belongs to the predicate result set.
	ChangeProjectionIgnored ChangeProjectionKind = iota
	// ChangeProjectionEntered means the change causes membership to appear.
	ChangeProjectionEntered
	// ChangeProjectionUpdated means the object matched before and after.
	ChangeProjectionUpdated
	// ChangeProjectionLeft means the change causes membership to disappear.
	ChangeProjectionLeft
)

// ChangeProjection reports whether a committed change crosses predicate membership.
type ChangeProjection struct {
	// Kind classifies the membership transition.
	Kind ChangeProjectionKind
	// Before records whether the committed Before state matched.
	Before bool
	// After records whether the committed After state matched.
	After bool
}

// ProjectChange evaluates a committed objectstore.Change against p.
//
// Projection is pure query semantics for future filtered watch/cache layers; it
// does not publish events, mutate caches, or interpret storage history.
func (p Predicate) ProjectChange(change objectstore.Change) (ChangeProjection, error) {
	if err := change.Validate(); err != nil {
		return ChangeProjection{}, invalidChangeError(err)
	}

	switch change.Kind {
	case objectstore.ChangeCreated:
		after := p.Match(objectstore.ListItem{Key: change.Key, State: change.After})
		if after {
			return ChangeProjection{Kind: ChangeProjectionEntered, After: true}, nil
		}
		return ChangeProjection{Kind: ChangeProjectionIgnored}, nil
	case objectstore.ChangeUpdated:
		before := p.Match(objectstore.ListItem{Key: change.Key, State: change.Before})
		after := p.Match(objectstore.ListItem{Key: change.Key, State: change.After})
		switch {
		case !before && after:
			return ChangeProjection{Kind: ChangeProjectionEntered, Before: before, After: after}, nil
		case before && after:
			return ChangeProjection{Kind: ChangeProjectionUpdated, Before: before, After: after}, nil
		case before && !after:
			return ChangeProjection{Kind: ChangeProjectionLeft, Before: before, After: after}, nil
		default:
			return ChangeProjection{Kind: ChangeProjectionIgnored}, nil
		}
	case objectstore.ChangeDeleted:
		before := p.Match(objectstore.ListItem{Key: change.Key, State: change.Before})
		if before {
			return ChangeProjection{Kind: ChangeProjectionLeft, Before: true}, nil
		}
		return ChangeProjection{Kind: ChangeProjectionIgnored}, nil
	default:
		return ChangeProjection{}, invalidChangeError(ErrInvalidChange)
	}
}

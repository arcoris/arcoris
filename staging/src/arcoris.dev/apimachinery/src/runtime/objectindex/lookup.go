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

import "arcoris.dev/apimachinery/api/objectstore"

// Lookup returns the keys currently indexed under name/value.
//
// The returned slice is detached. Results preserve membership order: Replace
// follows collection read order, and ApplyChange appends new memberships after
// existing memberships while preserving unchanged memberships in place.
func (i *Index) Lookup(name Name, value Value) ([]objectstore.Key, error) {
	if name == "" || value == "" {
		return nil, ErrInvalidIndex
	}
	if err := i.validateStatic(); err != nil {
		return nil, err
	}

	i.mu.RLock()
	defer i.mu.RUnlock()
	if err := i.validateStateLocked(); err != nil {
		return nil, err
	}

	values, ok := i.state.byName[name]
	if !ok {
		return nil, ErrUnknownIndex
	}
	keys := values[value]
	if keys == nil {
		return nil, nil
	}

	return keys.keys(), nil
}

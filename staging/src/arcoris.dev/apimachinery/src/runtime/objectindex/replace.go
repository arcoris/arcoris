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

	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

// Replace rebuilds all index memberships from a reflected collection read.
//
// The rebuild is all-or-nothing. Extractor failures or context cancellation
// leave the previous index state untouched.
func (i *Index) Replace(ctx context.Context, read storewatchapi.CollectionRead) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := i.validateStatic(); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := read.Validate(); err != nil {
		return err
	}

	next := newIndexState(i.names)
	for _, item := range read.Items() {
		if err := ctx.Err(); err != nil {
			return err
		}
		valuesByName, err := extractValues(i.names, i.definitions, item)
		if err != nil {
			return err
		}
		next.addObject(item.Key, valuesByName)
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	i.mu.Lock()
	i.state = next
	i.mu.Unlock()

	return nil
}

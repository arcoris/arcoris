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

package objectworkqueue

import (
	"errors"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Item is one object-keyed unit of reconciliation work.
type Item struct {
	// Key identifies the object whose reconciliation should be attempted. Add,
	// TryAdd, and Done reject items whose Key is not a valid objectstore key.
	Key objectstore.Key
}

// validateItem checks the object identity boundary before queue state changes.
func validateItem(item Item) error {
	if err := objectstore.ValidateKey(item.Key); err != nil {
		return errors.Join(ErrInvalidItem, err)
	}
	return nil
}

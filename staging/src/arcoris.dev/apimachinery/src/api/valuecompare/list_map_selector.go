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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/listmapkey"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// extractListMapSelector delegates selector construction to the shared helper.
func (c *comparer) extractListMapSelector(
	indexPath fieldpath.Path,
	item value.Value,
	element types.Type,
	keys []types.FieldName,
) (fieldpath.Selector, error) {
	selector, err := listmapkey.ExtractSelector(
		indexPath,
		item,
		element,
		keys,
		listmapkey.Options{Resolver: c.resolver, MaxDepth: c.maxDepth},
	)
	if err != nil {
		return fieldpath.Selector{}, compareListMapKeyError(err)
	}

	return selector, nil
}

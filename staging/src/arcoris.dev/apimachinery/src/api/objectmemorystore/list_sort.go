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

package objectmemorystore

import (
	"cmp"
	"slices"

	"arcoris.dev/apimachinery/api/objectstore"
)

// sortListItems puts collection reads into stable storage-key order.
func sortListItems(items []objectstore.ListItem) {
	slices.SortFunc(items, func(a, b objectstore.ListItem) int {
		return compareListKeys(a.Key, b.Key)
	})
}

// compareListKeys orders keys by full storage identity.
func compareListKeys(a, b objectstore.Key) int {
	if result := cmp.Compare(a.Resource.Group.String(), b.Resource.Group.String()); result != 0 {
		return result
	}
	if result := cmp.Compare(a.Resource.Version.String(), b.Resource.Version.String()); result != 0 {
		return result
	}
	if result := cmp.Compare(a.Resource.Resource.String(), b.Resource.Resource.String()); result != 0 {
		return result
	}
	if result := cmp.Compare(a.Object.Namespace.String(), b.Object.Namespace.String()); result != 0 {
		return result
	}

	return cmp.Compare(a.Object.Name.String(), b.Object.Name.String())
}

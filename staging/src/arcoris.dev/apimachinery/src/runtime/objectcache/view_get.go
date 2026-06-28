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

import "arcoris.dev/apimachinery/api/objectstore"

// Get returns the latest live state for key from v.
//
// A missing key is a known absence answer when key belongs to v's collection.
// Get returns detached state and does not consult Cache history records.
func (v View) Get(key objectstore.Key) (GetResult, error) {
	if !objectstore.KeyMatchesListRequest(key, v.collection) {
		return GetResult{}, ErrOutsideCollection
	}

	result := GetResult{Key: key, Revision: v.latest.revision}
	item, ok := v.latest.item(key)
	if !ok {
		return result, nil
	}

	result.State = item.State
	result.Found = true

	return result, nil
}

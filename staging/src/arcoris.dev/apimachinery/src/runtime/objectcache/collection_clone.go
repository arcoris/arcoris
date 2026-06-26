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

// clone returns a detached collection copy for validate-then-mutate paths and
// private tests. Public reads use narrower cloning helpers to avoid unnecessary
// map allocation.
func (col collection) clone() collection {
	out := collection{revision: col.revision}
	if len(col.order) == 0 {
		return out
	}

	out.order = append([]objectstore.Key(nil), col.order...)
	out.items = make(map[objectstore.Key]objectstore.ListItem, len(col.items))
	for key, item := range col.items {
		out.items[key] = item.Clone()
	}

	return out
}

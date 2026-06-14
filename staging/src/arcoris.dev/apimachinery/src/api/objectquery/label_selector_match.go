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

import (
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectstore"
)

// match checks item labels against the selector.
func (s LabelSelector) match(item objectstore.ListItem) bool {
	return matchLabelRequirements(s.requirements, func(key string) (string, bool) {
		value, ok := item.State.Object.ObjectMeta.Labels.Get(labels.Key(key))
		return value.String(), ok
	})
}

// matchLabelRequirements applies AND semantics without allocating per item.
func matchLabelRequirements(requirements []LabelRequirement, lookup func(string) (string, bool)) bool {
	for _, req := range requirements {
		if !req.req.match(lookup) {
			return false
		}
	}

	return true
}

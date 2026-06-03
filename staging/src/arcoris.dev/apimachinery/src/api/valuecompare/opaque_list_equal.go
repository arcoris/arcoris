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
	"arcoris.dev/apimachinery/api/value"
)

// equalOpaqueList compares unknown list payloads by physical order.
func (c *comparer) equalOpaqueList(path fieldpath.Path, oldValue value.Value, newValue value.Value) (bool, error) {
	oldList, _ := oldValue.List()
	newList, _ := newValue.List()
	n := oldList.Len()
	if n != newList.Len() {
		return false, nil
	}

	for i := 0; i < n; i++ {
		oldItem, _ := oldList.At(i)
		newItem, _ := newList.At(i)
		equal, err := c.equalOpaqueValue(path.Index(i), oldItem, newItem)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}

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

package valuemerge

import (
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/value"
)

// listItems returns detached list items for present list operands.
func listItems(o operand) []value.Value {
	if o.Absent() || o.Value().IsNull() {
		return nil
	}

	view, _ := o.Value().List()
	return view.Items()
}

// itemAt returns a presence-aware list item.
func itemAt(items []value.Value, index int) operand {
	if index < 0 || index >= len(items) {
		return valuepresence.Absent()
	}

	return valuepresence.Present(items[index])
}

// appendItem appends a cloned item when present.
func appendItem(items []value.Value, item operand) []value.Value {
	if item.Absent() {
		return items
	}

	return append(items, item.Value().Clone())
}

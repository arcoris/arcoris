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
	"strings"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectsurface"
	"arcoris.dev/apimachinery/api/value"
)

// fieldLookup is the internal result of resolving a selectable field path.
// present distinguishes missing from present null.
type fieldLookup struct {
	// value is meaningful only when present is true.
	value value.Value
	// present reports that every path segment existed.
	present bool
}

// lookupFieldValue resolves a registered surface-relative field path against a
// list item. It intentionally supports only explicitly queryable surfaces.
func lookupFieldValue(item objectstore.ListItem, ref FieldRef) fieldLookup {
	var root value.Value
	kinds := objectsurface.Kinds()
	switch ref.Surface {
	case kinds.Desired():
		root = item.State.Object.Desired
	case kinds.Observed():
		if item.State.Object.Observed == nil {
			return fieldLookup{}
		}
		root = *item.State.Object.Observed
	default:
		return fieldLookup{}
	}
	if root.IsZero() {
		return fieldLookup{}
	}

	current := root
	for _, segment := range strings.Split(ref.Path.String(), ".") {
		if segment == "" {
			return fieldLookup{}
		}
		record, ok := current.AsRecord()
		if !ok {
			return fieldLookup{}
		}
		next, ok := record.Get(value.MemberName(segment))
		if !ok {
			return fieldLookup{}
		}
		current = next
	}

	return fieldLookup{value: current, present: true}
}

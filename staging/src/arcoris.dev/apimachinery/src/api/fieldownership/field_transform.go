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

package fieldownership

import "arcoris.dev/apimachinery/api/fieldpath"

// unionTransform returns the canonical union of current and fields.
func unionTransform(current fieldpath.Set, fields fieldpath.Set) fieldpath.Set {
	switch {
	case current.IsEmpty():
		return fields
	case fields.IsEmpty():
		return current
	default:
		return current.Union(fields)
	}
}

// removeExactTransform removes only paths that exactly match fields.
func removeExactTransform(current fieldpath.Set, fields fieldpath.Set) fieldpath.Set {
	return current.Difference(fields)
}

// removeOverlapTransform removes current paths that structurally overlap fields.
func removeOverlapTransform(current fieldpath.Set, fields fieldpath.Set) fieldpath.Set {
	result := fieldpath.EmptySet()
	current.ForEach(func(_ int, p fieldpath.Path) bool {
		if overlapsAny(p, fields) {
			return true
		}

		result = result.Insert(p)
		return true
	})

	return result
}

// overlapsAny reports whether path overlaps any path in fields.
func overlapsAny(path fieldpath.Path, fields fieldpath.Set) bool {
	overlaps := false
	fields.ForEach(func(_ int, candidate fieldpath.Path) bool {
		if path.Overlaps(candidate) {
			overlaps = true
			return false
		}

		return true
	})

	return overlaps
}

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

package valueapply

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/valuecompare"
)

// changedAppliedFields returns the applied subset that would change Live.
//
// This set drives conflict detection. Unchanged applied fields can become
// shared ownership without conflicting with existing owners. Matching is
// structural rather than exact because lower layers may report a changed parent
// while Applied mentions a child, or the other way around.
func changedAppliedFields(applied fieldpath.Set, changes valuecompare.Result) fieldpath.Set {
	changed := changes.Changed()
	result := fieldpath.EmptySet()

	applied.ForEach(func(_ int, appliedPath fieldpath.Path) bool {
		if overlapsAny(changed, appliedPath) {
			result = result.Insert(appliedPath)
		}
		return true
	})

	return result
}

// droppedFields returns fields previously owned by Owner but omitted by Applied.
//
// Dropped fields are release/deletion candidates, not conflict candidates. A
// previous field remains owned when Applied contains that exact path or an
// ancestor path that covers it.
func droppedFields(oldOwnerFields fieldpath.Set, appliedFields fieldpath.Set) fieldpath.Set {
	dropped := fieldpath.EmptySet()

	oldOwnerFields.ForEach(func(_ int, oldPath fieldpath.Path) bool {
		if coveredByApplied(oldPath, appliedFields) {
			return true
		}

		dropped = dropped.Insert(oldPath)
		return true
	})

	return dropped
}

// mergeFields combines fields copied from Applied with dropped fields selected
// for deletion from Live.
func mergeFields(appliedFields fieldpath.Set, deletedFields fieldpath.Set) fieldpath.Set {
	return appliedFields.Union(deletedFields)
}

// overlapsAny reports whether path has structural overlap with any path in set.
func overlapsAny(set fieldpath.Set, path fieldpath.Path) bool {
	return set.Overlaps(path)
}

// coveredByApplied reports whether Applied keeps ownership of oldPath.
func coveredByApplied(oldPath fieldpath.Path, appliedFields fieldpath.Set) bool {
	covered := false
	appliedFields.ForEach(func(_ int, appliedPath fieldpath.Path) bool {
		if oldPath.HasPrefix(appliedPath) {
			covered = true
			return false
		}
		return true
	})

	return covered
}

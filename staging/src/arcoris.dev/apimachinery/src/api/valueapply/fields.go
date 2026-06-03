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
// shared ownership without conflicting with existing owners.
func changedAppliedFields(applied fieldpath.Set, changes valuecompare.Result) fieldpath.Set {
	return applied.Intersection(changes.Changed())
}

// droppedFields returns fields previously owned by Owner but omitted by Applied.
//
// Dropped fields are release/deletion candidates, not conflict candidates.
func droppedFields(oldOwnerFields fieldpath.Set, appliedFields fieldpath.Set) fieldpath.Set {
	return oldOwnerFields.Difference(appliedFields)
}

// mergeFields combines fields copied from Applied with dropped fields selected
// for deletion from Live.
func mergeFields(appliedFields fieldpath.Set, deletedFields fieldpath.Set) fieldpath.Set {
	return appliedFields.Union(deletedFields)
}

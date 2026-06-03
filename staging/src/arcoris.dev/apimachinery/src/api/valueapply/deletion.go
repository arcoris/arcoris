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
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
)

// deletableDroppedFields filters dropped fields to those not protected by any
// other owner.
//
// Structural overlap is conservative: exact, ancestor, and descendant ownership
// by another owner all preserve the live value and only release Owner.
func deletableDroppedFields(
	ownership fieldownership.State,
	owner fieldownership.Owner,
	dropped fieldpath.Set,
) fieldpath.Set {
	deletable := fieldpath.EmptySet()
	for _, path := range dropped.Paths() {
		if hasOtherOverlappingOwner(ownership, owner, path) {
			continue
		}

		deletable = deletable.Insert(path)
	}

	return deletable
}

// hasOtherOverlappingOwner reports whether path is protected by an owner other
// than the applying owner.
func hasOtherOverlappingOwner(
	ownership fieldownership.State,
	owner fieldownership.Owner,
	path fieldpath.Path,
) bool {
	for _, record := range ownership.OverlappingOwners(path) {
		if record.Owner != owner {
			return true
		}
	}

	return false
}

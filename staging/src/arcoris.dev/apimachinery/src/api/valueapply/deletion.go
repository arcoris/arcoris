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

// planMergeFields computes deletion and merge fields after conflict and force
// policy allow the apply to proceed.
func (a *applier) planMergeFields(req Request, prepared preparedApply) (preparedApply, error) {
	deleted, err := deletableDroppedFields(req.Ownership, req.Owner, prepared.DroppedFields)
	if err != nil {
		return prepared, err
	}

	prepared.DeletedFields = deleted
	prepared.MergeFields = mergeFields(prepared.AppliedFields, prepared.DeletedFields)

	return prepared, nil
}

// deletableDroppedFields filters dropped fields to those not protected by any
// other owner.
//
// Structural overlap is conservative: exact, ancestor, and descendant ownership
// by another owner all preserve the live value and only release Owner.
func deletableDroppedFields(
	ownership fieldownership.State,
	owner fieldownership.Owner,
	dropped fieldpath.Set,
) (fieldpath.Set, error) {
	deletable := fieldpath.EmptySet()
	var overlapErr error

	dropped.ForEach(func(_ int, path fieldpath.Path) bool {
		protected, err := hasOtherOverlappingOwner(ownership, owner, path)
		if err != nil {
			overlapErr = wrapAt(
				path,
				ErrInvalidRequest,
				ErrorReasonInvalidRequest,
				"failed to check dropped field ownership overlap",
				err,
			)
			return false
		}
		if protected {
			return true
		}

		deletable = deletable.Insert(path)
		return true
	})
	if overlapErr != nil {
		return fieldpath.Set{}, overlapErr
	}

	return deletable, nil
}

// hasOtherOverlappingOwner reports whether path is protected by an owner other
// than the applying owner.
func hasOtherOverlappingOwner(
	ownership fieldownership.State,
	owner fieldownership.Owner,
	path fieldpath.Path,
) (bool, error) {
	records, err := ownership.OverlappingPaths(path)
	if err != nil {
		return false, err
	}

	hasOther := false
	records.ForEach(func(_ int, record fieldownership.OwnedPath) bool {
		if record.Owner != owner {
			hasOther = true
			return false
		}
		return true
	})

	return hasOther, nil
}

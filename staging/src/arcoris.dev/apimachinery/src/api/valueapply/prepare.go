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

// prepare computes the deterministic field metadata needed before conflict
// handling and merge.
func (a *applier) prepare(req Request) (Result, error) {
	result := Result{}

	appliedFields, err := a.extractAppliedFields(req)
	if err != nil {
		return result, err
	}
	result.AppliedFields = appliedFields

	oldOwnerFields := req.Ownership.FieldsFor(req.Owner)
	result.DroppedFields = droppedFields(oldOwnerFields, result.AppliedFields)

	changes, err := a.compare(req)
	if err != nil {
		return result, err
	}
	result.Changes = changes

	result.ChangedAppliedFields = changedAppliedFields(result.AppliedFields, result.Changes)

	conflicts, err := ownershipConflicts(req, result.ChangedAppliedFields)
	if err != nil {
		return result, err
	}
	result.Conflicts = conflicts

	result.DeletedFields = deletableDroppedFields(
		req.Ownership,
		req.Owner,
		result.DroppedFields,
	)
	result.MergeFields = mergeFields(result.AppliedFields, result.DeletedFields)

	return result, nil
}

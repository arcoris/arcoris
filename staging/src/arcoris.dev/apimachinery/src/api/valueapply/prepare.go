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

// prepare computes the deterministic field metadata needed before conflict,
// merge, and ownership-update policy.
func (a *applier) prepare(req Request) (preparedApply, error) {
	prepared := preparedApply{}

	appliedFields, err := a.extractAppliedFields(req)
	if err != nil {
		return prepared, err
	}
	prepared.AppliedFields = appliedFields

	oldOwnerFields := req.Ownership.FieldsFor(req.Owner)
	prepared.DroppedFields = droppedFields(oldOwnerFields, prepared.AppliedFields)

	changes, err := a.compare(req)
	if err != nil {
		return prepared, err
	}
	prepared.Changes = changes

	prepared.ChangedAppliedFields = changedAppliedFields(prepared.AppliedFields, prepared.Changes)

	conflicts, err := ownershipConflicts(req, prepared.ChangedAppliedFields)
	if err != nil {
		return prepared, err
	}
	prepared.Conflicts = conflicts

	return prepared, nil
}

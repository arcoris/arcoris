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

// ownershipConflicts asks fieldownership for structural overlaps against only
// changed applied fields.
func ownershipConflicts(
	req Request,
	changed fieldpath.Set,
) (fieldownership.ConflictSet, error) {
	conflicts, err := req.Ownership.Conflicts(req.Owner, changed)
	if err != nil {
		return fieldownership.ConflictSet{}, wrapAt(
			req.Path,
			ErrConflict,
			ErrorReasonConflict,
			"field ownership conflict check failed",
			err,
		)
	}

	return conflicts, nil
}

// rejectConflicts converts non-empty conflicts into a valueapply error unless
// Force has explicitly selected takeover behavior.
func (a *applier) rejectConflicts(req Request, prepared preparedApply) error {
	if prepared.Conflicts.IsEmpty() || a.opts.Force {
		return nil
	}

	return conflictError(req.Path, prepared.Conflicts)
}

// conflictError wraps a non-empty fieldownership conflict set with the
// valueapply sentinel while preserving fieldownership.ErrConflict as a cause.
func conflictError(path fieldpath.Path, conflicts fieldownership.ConflictSet) error {
	return wrapAt(
		path,
		ErrConflict,
		ErrorReasonConflict,
		"changed applied fields conflict with existing ownership",
		fieldownership.NewConflictError(conflicts),
	)
}

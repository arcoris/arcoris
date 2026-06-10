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

// rejectUnsupportedForceTakeover rejects forced conflicts that would require
// subtracting a child path from another owner's ancestor ownership.
func (a *applier) rejectUnsupportedForceTakeover(req Request, result Result) error {
	if !a.opts.Force || result.Conflicts.IsEmpty() {
		return nil
	}

	unsupported := unsupportedTakeoverConflicts(result.Conflicts)
	if unsupported.IsEmpty() {
		return nil
	}

	return unsupportedTakeoverError(req.Path, unsupported)
}

// unsupportedTakeoverConflicts returns conflicts whose owned path strictly
// contains the attempted path.
func unsupportedTakeoverConflicts(
	conflicts fieldownership.ConflictSet,
) fieldownership.ConflictSet {
	unsupported := make([]fieldownership.Conflict, 0, conflicts.Len())

	for _, conflict := range conflicts.Conflicts() {
		if ownedPathStrictlyContainsAttempted(conflict) {
			unsupported = append(unsupported, conflict)
		}
	}

	return fieldownership.NewConflictSet(unsupported...)
}

// ownedPathStrictlyContainsAttempted reports whether force would over-remove
// another owner's ancestor field.
func ownedPathStrictlyContainsAttempted(conflict fieldownership.Conflict) bool {
	return conflict.AttemptedPath.HasPrefix(conflict.OwnedPath) &&
		!conflict.AttemptedPath.Equal(conflict.OwnedPath)
}

// unsupportedTakeoverError preserves the fieldownership conflict cause while
// classifying the valueapply policy failure separately.
func unsupportedTakeoverError(
	path fieldpath.Path,
	conflicts fieldownership.ConflictSet,
) error {
	return wrapAt(
		path,
		ErrUnsupportedTakeover,
		ErrorReasonUnsupportedTakeover,
		"forced ownership takeover cannot precisely subtract ancestor ownership",
		fieldownership.NewConflictError(conflicts),
	)
}

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

// updateOwnership returns the post-apply ownership state after merge succeeds.
//
// Force removes only conflicting attempted paths from other owners. The current
// owner is then set exactly to AppliedFields, releasing omitted fields.
func (a *applier) updateOwnership(req Request, merged mergedApply) (completedApply, error) {
	ownership := req.Ownership
	var err error

	if a.opts.Force && !merged.Conflicts.IsEmpty() {
		ownership, err = ownership.RemoveOverlappingFieldsFromOthers(
			req.Owner,
			merged.Conflicts.AttemptedPaths(),
		)
		if err != nil {
			return completedApply{}, wrapAt(
				req.Path,
				ErrInvalidRequest,
				ErrorReasonInvalidRequest,
				"forced ownership takeover failed",
				err,
			)
		}
	}

	ownership, err = ownership.SetFields(req.Owner, merged.AppliedFields)
	if err != nil {
		return completedApply{}, wrapAt(
			req.Path,
			ErrInvalidRequest,
			ErrorReasonInvalidRequest,
			"owner field update failed",
			err,
		)
	}

	return completedApply{
		mergedApply: merged,
		Ownership:   ownership,
	}, nil
}

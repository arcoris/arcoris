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

// validateConflict checks one conflict record.
func validateConflict(conflict Conflict) error {
	if err := conflict.Owner.ValidateLexical(); err != nil {
		return wrapAt(
			"conflict.owner",
			ErrInvalidConflict,
			ErrorReasonInvalidConflictOwner,
			"conflict owner is invalid",
			err,
		)
	}
	if err := conflict.OwnedPath.ValidateStructure(); err != nil {
		return wrapAt(
			"conflict.ownedPath",
			ErrInvalidConflict,
			ErrorReasonInvalidConflictOwnedPath,
			"conflict owned path is invalid",
			err,
		)
	}
	if err := conflict.AttemptedPath.ValidateStructure(); err != nil {
		return wrapAt(
			"conflict.attemptedPath",
			ErrInvalidConflict,
			ErrorReasonInvalidConflictAttemptedPath,
			"conflict attempted path is invalid",
			err,
		)
	}

	return nil
}

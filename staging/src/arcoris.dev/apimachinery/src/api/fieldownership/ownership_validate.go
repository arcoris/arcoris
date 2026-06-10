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

// validateOwnedPath checks one owner/path query record.
func validateOwnedPath(path OwnedPath) error {
	if err := path.Owner.ValidateLexical(); err != nil {
		return wrapAt(
			"ownedPath.owner",
			ErrInvalidOwnedPath,
			ErrorReasonInvalidOwnedPathOwner,
			"owned path owner is invalid",
			err,
		)
	}
	if err := path.Path.ValidateStructure(); err != nil {
		return wrapAt(
			"ownedPath.path",
			ErrInvalidOwnedPath,
			ErrorReasonInvalidOwnedPathPath,
			"owned path is invalid",
			err,
		)
	}

	return nil
}

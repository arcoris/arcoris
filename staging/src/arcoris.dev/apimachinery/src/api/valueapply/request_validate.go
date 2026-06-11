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

// validateRequestShape checks caller-owned request metadata before invoking
// lower-level value packages.
func validateRequestShape(req Request) error {
	if err := req.Path.ValidateStructure(); err != nil {
		return wrapAt(
			req.Path,
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"base field path is invalid",
			err,
		)
	}
	if err := req.Owner.ValidateLexical(); err != nil {
		return wrapAt(
			req.Path,
			ErrInvalidOwner,
			ErrorReasonInvalidOwner,
			"field owner is invalid",
			err,
		)
	}
	if err := req.Ownership.ValidateStructure(); err != nil {
		return wrapAt(
			req.Path,
			ErrInvalidRequest,
			ErrorReasonInvalidRequest,
			"ownership state is invalid",
			err,
		)
	}

	return validateOwnershipScope(req.Path, req.Ownership)
}

// validateOwnershipScope rejects ownership paths outside the apply base path.
func validateOwnershipScope(base fieldpath.Path, ownership fieldownership.State) error {
	var scopeErr error
	ownership.ForEach(func(_ int, entry fieldownership.Entry) bool {
		entry.Fields().ForEach(func(_ int, ownedPath fieldpath.Path) bool {
			if ownedPath.Equal(base) || ownedPath.IsDescendantOf(base) {
				return true
			}

			scopeErr = errorfAt(
				ownedPath,
				ErrInvalidRequest,
				ErrorReasonInvalidRequest,
				"ownership field path %s is outside apply base path %s",
				ownedPath.CanonicalText(),
				base.CanonicalText(),
			)
			return false
		})

		return scopeErr == nil
	})

	return scopeErr
}

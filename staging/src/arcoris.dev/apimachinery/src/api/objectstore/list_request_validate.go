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

package objectstore

import "errors"

// IsValid reports whether request is structurally usable by Store.List.
func (request ListRequest) IsValid() bool {
	return ValidateListRequest(request) == nil
}

// ValidateListRequest checks the structural shape required for Store.List.
func ValidateListRequest(request ListRequest) error {
	if err := request.Resource.Validate(); err != nil {
		return errorFor(ErrorReasonInvalidListRequest, Key{}, 0, 0, errors.Join(ErrInvalidListRequest, err))
	}
	if err := ValidateListScope(request.Scope); err != nil {
		return err
	}

	return nil
}

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

// IsValid reports whether i is shaped like one committed list item.
func (i ListItem) IsValid() bool {
	return ValidateListItem(i) == nil
}

// ValidateListItem checks one committed live item returned by Store.List.
//
// The key is the authoritative storage identity. The state's metadata
// namespace/name must agree with that key so callers can safely use list items
// for later watch projection, cache indexing, or lifecycle reads without
// reconciling two identities.
func ValidateListItem(item ListItem) error {
	if err := ValidateKey(item.Key); err != nil {
		return errorFor(
			ErrorReasonInvalidListItem,
			item.Key,
			0,
			0,
			errors.Join(ErrInvalidListResult, err),
		)
	}
	if err := ValidateCommittedState(item.State); err != nil {
		return errorFor(
			ErrorReasonInvalidListItem,
			item.Key,
			item.State.Revision,
			0,
			errors.Join(ErrInvalidListResult, err),
		)
	}
	if !stateIdentityMatchesKey(item.Key, item.State) {
		return errorFor(
			ErrorReasonInvalidListItem,
			item.Key,
			0,
			0,
			errors.Join(ErrInvalidListResult, ErrInvalidKey),
		)
	}

	return nil
}

// stateIdentityMatchesKey checks the metadata namespace/name carried by state
// against the authoritative storage identity.
func stateIdentityMatchesKey(key Key, state State) bool {
	return state.Object.ObjectMeta.Namespace == key.Object.Namespace &&
		state.Object.ObjectMeta.Name == key.Object.Name
}

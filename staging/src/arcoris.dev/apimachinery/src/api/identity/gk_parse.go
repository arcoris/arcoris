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

package identity

// ParseGroupKind parses a canonical group/kind identity.
//
// Core-group kinds use only "Kind". Named groups use "group#Kind". The parser
// rejects the old dotted "Kind.group" grammar so kind identity cannot be
// confused with a Go package-like name.
func ParseGroupKind(value string) (GroupKind, error) {
	groupPart, kindPart, hasGroup, err := splitOptionalPair(
		identityNameGroupKind,
		value,
		kindSeparator,
		detailExpectedGroupKind,
	)
	if err != nil {
		if value == "" {
			return GroupKind{}, invalid(
				identityNameGroupKind,
				value,
				ErrorReasonEmptyValue,
				detailKindRequired,
			)
		}

		return GroupKind{}, err
	}

	var group Group
	if hasGroup {
		parsedGroup, err := ParseGroup(groupPart)
		if err != nil {
			return GroupKind{}, invalidValue(identityNameGroupKind, value, err)
		}

		group = parsedGroup
	} else {
		kindPart = groupPart
	}

	kind, err := ParseKind(kindPart)
	if err != nil {
		return GroupKind{}, invalidValue(identityNameGroupKind, value, err)
	}

	return GroupKind{Group: group, Kind: kind}, nil
}

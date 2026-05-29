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

// ParseGroupVersion parses a canonical group/version identity.
//
// Core-group versions are written as "v1". Named-group versions are written as
// "group/v1". URL-like resource paths are not accepted here; resources use the
// colon separator in GroupVersionResource.
func ParseGroupVersion(value string) (GroupVersion, error) {
	groupPart, versionPart, hasGroup, err := splitOptionalPair(
		identityNameGroupVersion,
		value,
		groupVersionSeparator,
		detailExpectedGroupVersion,
	)
	if err != nil {
		if value == "" {
			return GroupVersion{}, invalid(
				identityNameGroupVersion,
				value,
				ErrorReasonEmptyValue,
				detailVersionRequired,
			)
		}

		return GroupVersion{}, err
	}

	var group Group
	if hasGroup {
		parsedGroup, err := ParseGroup(groupPart)
		if err != nil {
			return GroupVersion{}, invalidValue(identityNameGroupVersion, value, err)
		}

		group = parsedGroup
	} else {
		versionPart = groupPart
	}

	version, err := ParseVersion(versionPart)
	if err != nil {
		return GroupVersion{}, invalidValue(identityNameGroupVersion, value, err)
	}

	return GroupVersion{Group: group, Version: version}, nil
}

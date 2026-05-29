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

// ParseGroupResource parses a canonical group/resource identity.
//
// Core resources use only "resource". Named groups use "group:resource". The
// old dotted "resource.group" form is rejected so group and resource axes stay
// visually distinct.
func ParseGroupResource(value string) (GroupResource, error) {
	groupPart, resourcePart, hasGroup, err := splitOptionalPair(
		identityNameGroupResource,
		value,
		resourceSeparator,
		detailExpectedGroupResource,
	)
	if err != nil {
		if value == "" {
			return GroupResource{}, invalid(
				identityNameGroupResource,
				value,
				ErrorReasonEmptyValue,
				detailResourceRequired,
			)
		}

		return GroupResource{}, err
	}

	var group Group
	if hasGroup {
		parsedGroup, err := ParseGroup(groupPart)
		if err != nil {
			return GroupResource{}, invalidValue(identityNameGroupResource, value, err)
		}

		group = parsedGroup
	} else {
		resourcePart = groupPart
	}

	resource, err := ParseResource(resourcePart)
	if err != nil {
		return GroupResource{}, invalidValue(identityNameGroupResource, value, err)
	}

	return GroupResource{Group: group, Resource: resource}, nil
}

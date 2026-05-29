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

// ParseGroupVersionResource parses a canonical group/version/resource identity.
//
// Resource identity uses ":" after GroupVersion. That keeps it distinct from
// group/version parsing and from URL paths, which are not identity grammar.
func ParseGroupVersionResource(value string) (GroupVersionResource, error) {
	gvPart, resourcePart, err := splitRequiredPair(
		identityNameGroupVersionResource,
		value,
		resourceSeparator,
		detailExpectedGroupVersionResource,
	)
	if err != nil {
		if value == "" {
			return GroupVersionResource{}, invalid(
				identityNameGroupVersionResource,
				value,
				ErrorReasonEmptyValue,
				detailGroupVersionAndResourceRequired,
			)
		}

		return GroupVersionResource{}, err
	}

	gv, err := ParseGroupVersion(gvPart)
	if err != nil {
		return GroupVersionResource{}, invalidValue(identityNameGroupVersionResource, value, err)
	}

	resource, err := ParseResource(resourcePart)
	if err != nil {
		return GroupVersionResource{}, invalidValue(identityNameGroupVersionResource, value, err)
	}

	return GroupVersionResource{
		Group:    gv.Group,
		Version:  gv.Version,
		Resource: resource,
	}, nil
}

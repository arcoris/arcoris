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

// ParseGroupVersionKind parses a canonical group/version/kind identity.
//
// A kind is separated from GroupVersion with "#". Comma-based diagnostics and
// split object-field forms are intentionally not part of identity parsing.
func ParseGroupVersionKind(value string) (GroupVersionKind, error) {
	gvPart, kindPart, err := splitRequiredPair(
		identityNameGroupVersionKind,
		value,
		kindSeparator,
		detailExpectedGroupVersionKind,
	)
	if err != nil {
		if value == "" {
			return GroupVersionKind{}, invalid(
				identityNameGroupVersionKind,
				value,
				ErrorReasonEmptyValue,
				detailGroupVersionAndKindRequired,
			)
		}

		return GroupVersionKind{}, err
	}

	gv, err := ParseGroupVersion(gvPart)
	if err != nil {
		return GroupVersionKind{}, invalidValue(identityNameGroupVersionKind, value, err)
	}

	kind, err := ParseKind(kindPart)
	if err != nil {
		return GroupVersionKind{}, invalidValue(identityNameGroupVersionKind, value, err)
	}

	return GroupVersionKind{
		Group:   gv.Group,
		Version: gv.Version,
		Kind:    kind,
	}, nil
}

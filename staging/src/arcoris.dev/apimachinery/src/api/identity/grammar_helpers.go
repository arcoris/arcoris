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

import "strings"

// joinGroupVersion returns the canonical GroupVersion spelling.
//
// The core group is omitted instead of serialized as an empty prefix. That
// keeps "v1" as the only core-group group/version spelling.
func joinGroupVersion(group Group, version Version) string {
	if group.IsZero() {
		return version.String()
	}

	return group.String() + groupVersionSeparator + version.String()
}

// joinGroupKind returns the canonical GroupKind spelling.
//
// Kind identity uses its own separator so it cannot be confused with resource
// collection identity or URL-like paths.
func joinGroupKind(group Group, kind Kind) string {
	if group.IsZero() {
		return kind.String()
	}

	return group.String() + kindSeparator + kind.String()
}

// joinGroupResource returns the canonical GroupResource spelling.
//
// Resource identity uses a colon for named groups and omits the core group
// prefix for core resources.
func joinGroupResource(group Group, resource Resource) string {
	if group.IsZero() {
		return resource.String()
	}

	return group.String() + resourceSeparator + resource.String()
}

// joinGroupVersionKind returns the canonical GroupVersionKind spelling.
func joinGroupVersionKind(gv GroupVersion, kind Kind) string {
	return gv.String() + kindSeparator + kind.String()
}

// joinGroupVersionResource returns the canonical GroupVersionResource spelling.
func joinGroupVersionResource(gv GroupVersion, resource Resource) string {
	return gv.String() + resourceSeparator + resource.String()
}

// joinResourcePath returns the canonical ResourcePath spelling.
func joinResourcePath(resource Resource, subresource Subresource) string {
	if subresource.IsZero() {
		return resource.String()
	}

	return resource.String() + subresourceSeparator + subresource.String()
}

// joinGroupVersionResourcePath returns the canonical versioned ResourcePath spelling.
func joinGroupVersionResourcePath(gv GroupVersion, path ResourcePath) string {
	return gv.String() + resourceSeparator + path.String()
}

// splitOptionalPair splits either a one-part identity or one separator-delimited pair.
//
// Optional-pair identities have a compact core-group spelling and an expanded
// named-group spelling. The helper rejects repeated separators and empty
// segments before callers validate the concrete segment grammars.
func splitOptionalPair(
	name string,
	value string,
	separator string,
	want string,
) (left string, right string, hasSeparator bool, err error) {
	if strings.Count(value, separator) > 1 {
		return "", "", false, invalid(
			name,
			value,
			ErrorReasonInvalidForm,
			want,
		)
	}

	left, right, hasSeparator = strings.Cut(value, separator)
	if hasSeparator && (left == "" || right == "") {
		return "", "", false, invalid(
			name,
			value,
			ErrorReasonInvalidForm,
			want+" without empty segments",
		)
	}

	return left, right, hasSeparator, nil
}

// splitRequiredPair splits an identity that must contain exactly one separator.
//
// Required-pair identities always have two semantic sides, such as
// GroupVersion and Kind. The helper only validates the shape; callers still own
// parsing and validating the left and right identity values.
func splitRequiredPair(
	name string,
	value string,
	separator string,
	want string,
) (left string, right string, err error) {
	if strings.Count(value, separator) != 1 {
		return "", "", invalid(
			name,
			value,
			ErrorReasonInvalidForm,
			want,
		)
	}

	left, right, _ = strings.Cut(value, separator)
	if left == "" || right == "" {
		return "", "", invalid(
			name,
			value,
			ErrorReasonInvalidForm,
			want+" without empty segments",
		)
	}

	return left, right, nil
}

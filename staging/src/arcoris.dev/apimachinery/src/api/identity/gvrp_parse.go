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

// ParseGroupVersionResourcePath parses a canonical versioned resource path.
//
// The parser first splits GroupVersion from ResourcePath with ":". The right
// side is then parsed by ResourcePath, which owns the optional subresource
// slash. URL-like group/version/resource paths are therefore rejected.
func ParseGroupVersionResourcePath(value string) (GroupVersionResourcePath, error) {
	gvPart, pathPart, err := splitRequiredPair(
		identityNameGroupVersionResourcePath,
		value,
		resourceSeparator,
		detailExpectedGroupVersionResourcePath,
	)
	if err != nil {
		if value == "" {
			return GroupVersionResourcePath{}, invalid(
				identityNameGroupVersionResourcePath,
				value,
				ErrorReasonEmptyValue,
				detailGroupVersionAndResourceRequired,
			)
		}

		return GroupVersionResourcePath{}, err
	}

	gv, err := ParseGroupVersion(gvPart)
	if err != nil {
		return GroupVersionResourcePath{}, invalidValue(
			identityNameGroupVersionResourcePath,
			value,
			err,
		)
	}

	path, err := ParseResourcePath(pathPart)
	if err != nil {
		return GroupVersionResourcePath{}, invalidValue(
			identityNameGroupVersionResourcePath,
			value,
			err,
		)
	}

	return GroupVersionResourcePath{
		Group:       gv.Group,
		Version:     gv.Version,
		Resource:    path.Resource,
		Subresource: path.Subresource,
	}, nil
}

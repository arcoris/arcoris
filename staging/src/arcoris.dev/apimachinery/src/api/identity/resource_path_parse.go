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

// ParseResourcePath parses a canonical resource/subresource path.
//
// Resource paths have at most one slash. The slash belongs only between
// Resource and Subresource; versioned resource identities use ":" before this
// path is parsed.
func ParseResourcePath(value string) (ResourcePath, error) {
	resourcePart, subresourcePart, hasSubresource, err := splitOptionalPair(
		identityNameResourcePath,
		value,
		subresourceSeparator,
		detailExpectedResourcePath,
	)
	if err != nil {
		if value == "" {
			return ResourcePath{}, invalid(
				identityNameResourcePath,
				value,
				ErrorReasonEmptyValue,
				detailResourceRequired,
			)
		}

		return ResourcePath{}, err
	}

	var subresource Subresource
	if hasSubresource {
		parsedSubresource, err := ParseSubresource(subresourcePart)
		if err != nil {
			return ResourcePath{}, invalidValue(identityNameResourcePath, value, err)
		}

		subresource = parsedSubresource
	}

	resource, err := ParseResource(resourcePart)
	if err != nil {
		return ResourcePath{}, invalidValue(identityNameResourcePath, value, err)
	}

	return ResourcePath{
		Resource:    resource,
		Subresource: subresource,
	}, nil
}

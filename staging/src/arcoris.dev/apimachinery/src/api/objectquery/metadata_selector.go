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

package objectquery

import "slices"

// canonicalMetadataRequirements validates, sorts, and deduplicates requirements.
func canonicalMetadataRequirements(
	path string,
	requirements []metadataRequirement,
	validateKey validatorFunc,
	validateValue validatorFunc,
) ([]metadataRequirement, error) {
	if len(requirements) == 0 {
		return nil, nil
	}

	out := make([]metadataRequirement, len(requirements))
	for i, req := range requirements {
		cloned := req.clone()
		if err := cloned.validate(path+"["+itoa(i)+"]", validateKey, validateValue); err != nil {
			return nil, err
		}
		cloned.values = canonicalValues(cloned.values)
		out[i] = cloned
	}

	slices.SortFunc(out, compareMetadataRequirement)
	return compactMetadataRequirements(out), nil
}

// compactMetadataRequirements removes exact duplicate canonical requirements.
func compactMetadataRequirements(requirements []metadataRequirement) []metadataRequirement {
	if len(requirements) == 0 {
		return nil
	}

	out := requirements[:1]
	for _, req := range requirements[1:] {
		if sameMetadataRequirement(req, out[len(out)-1]) {
			continue
		}
		out = append(out, req)
	}

	return out
}

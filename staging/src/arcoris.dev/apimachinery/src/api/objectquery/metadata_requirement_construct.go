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

// metadataRequirementFrom constructs, validates, and canonicalizes one
// metadata requirement for a concrete metadata domain.
//
// values is copied before validation so callers never retain a mutable handle
// into the compiled requirement.
func metadataRequirementFrom(
	path string,
	key string,
	op Operator,
	values []string,
	validateKey validatorFunc,
	validateValue validatorFunc,
) (metadataRequirement, error) {
	req := metadataRequirement{
		key:    key,
		op:     op,
		values: append([]string(nil), values...),
	}
	if err := req.validate(path, validateKey, validateValue); err != nil {
		return metadataRequirement{}, err
	}

	req.values = canonicalMetadataValues(req.values)
	return req, nil
}

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

package codecselection

import "slices"

// normalizeParametersAt validates and sorts parameters at path.
func normalizeParametersAt(path string, items []Parameter) (Parameters, error) {
	if len(items) == 0 {
		return Parameters{}, nil
	}

	normalized := make([]Parameter, 0, len(items))
	for i, item := range items {
		parameter, err := normalizeParameterAt(parameterPath(path, i), item)
		if err != nil {
			return Parameters{}, err
		}
		normalized = append(normalized, parameter)
	}
	slices.SortFunc(normalized, compareParameters)

	for i := 1; i < len(normalized); i++ {
		if normalized[i-1].name == normalized[i].name {
			return Parameters{}, errorfAt(
				parameterPath(path, i),
				ErrInvalidParameters,
				ErrorReasonInvalidParameters,
				"parameter name %q is duplicated",
				normalized[i].name,
			)
		}
	}

	return Parameters{items: normalized}, nil
}

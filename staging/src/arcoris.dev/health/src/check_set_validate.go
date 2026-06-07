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

package health

import "errors"

// prepareCheckSet validates checks and builds the immutable storage used by
// CheckSet.
//
// The returned checker slice is freshly allocated when checks are present. The
// names map is lookup-only storage for CheckSet.Has. Duplicate diagnostics keep
// caller indexes so registry and test code can point at the conflicting input.
func prepareCheckSet(checks []Checker) ([]Checker, map[string]struct{}, error) {
	if len(checks) == 0 {
		return nil, nil, nil
	}

	prepared := make([]Checker, 0, len(checks))
	names := make(map[string]struct{}, len(checks))
	indexes := make(map[string]int, len(checks))

	for index, checker := range checks {
		if err := ValidateChecker(checker); err != nil {
			if errors.Is(err, ErrNilChecker) {
				return nil, nil, NilCheckError{Index: index}
			}

			return nil, nil, InvalidCheckNameError{
				Index: index,
				Name:  checker.Name(),
				Err:   err,
			}
		}

		name := checker.Name()
		if previous, exists := indexes[name]; exists {
			return nil, nil, DuplicateCheckNameError{
				Name:          name,
				Index:         index,
				PreviousIndex: previous,
			}
		}

		indexes[name] = index
		names[name] = struct{}{}
		prepared = append(prepared, checker)
	}

	return prepared, names, nil
}

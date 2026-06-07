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

package healthregistry

import (
	"errors"

	"arcoris.dev/health"
)

type preparedCheck struct {
	// index is the caller-supplied batch position used in diagnostics.
	index int

	// name is the validated stable checker name.
	name string

	// checker is the validated non-nil checker.
	checker health.Checker
}

// prepareChecks validates one registration batch without mutating Builder.
//
// The returned slice preserves caller order. All validation failures are joined
// so callers can inspect multiple malformed entries while the batch still fails
// atomically.
func prepareChecks(target health.Target, checks []health.Checker) ([]preparedCheck, error) {
	if len(checks) == 0 {
		return nil, nil
	}

	prepared := make([]preparedCheck, 0, len(checks))
	seen := make(map[string]int, len(checks))
	var failures []error

	for index, checker := range checks {
		name, err := health.CheckerName(checker)
		if err != nil {
			failures = append(failures, registrationValidationError(target, index, name, err))
			continue
		}

		if previous, ok := seen[name]; ok {
			failures = append(failures, DuplicateCheckError{
				Target:        target,
				Name:          name,
				Index:         index,
				PreviousIndex: previous,
			})
			continue
		}

		seen[name] = index
		prepared = append(prepared, preparedCheck{
			index:   index,
			name:    name,
			checker: checker,
		})
	}

	if len(failures) > 0 {
		return nil, errors.Join(failures...)
	}

	return prepared, nil
}

// registrationValidationError converts root checker validation into
// registry-owned target/index diagnostics.
func registrationValidationError(
	target health.Target,
	index int,
	name string,
	err error,
) error {
	if errors.Is(err, health.ErrNilChecker) {
		return NilCheckerError{
			Target: target,
			Index:  index,
		}
	}

	return InvalidCheckNameError{
		Target: target,
		Index:  index,
		Name:   name,
		Err:    err,
	}
}

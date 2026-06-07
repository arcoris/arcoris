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
	"reflect"

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
		if nilChecker(checker) {
			failures = append(failures, NilCheckerError{
				Target: target,
				Index:  index,
			})
			continue
		}

		name := checker.Name()
		if err := health.ValidateCheckName(name); err != nil {
			failures = append(failures, InvalidCheckNameError{
				Target: target,
				Index:  index,
				Name:   name,
				Err:    err,
			})
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

// nilChecker reports whether checker is nil, including typed nil interface
// values.
func nilChecker(checker health.Checker) bool {
	if checker == nil {
		return true
	}

	val := reflect.ValueOf(checker)
	switch val.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package health

import (
	"errors"
	"reflect"
)

// preparedCheck carries validated batch metadata until Register can atomically
// publish the checker under the registry lock.
type preparedCheck struct {
	Index   int
	Name    string
	Checker Checker
}

// prepareChecks validates a registration batch before the registry is mutated.
//
// It collects all batch-local validation failures so callers can fix invalid
// setup in one pass. Existing-registry conflicts are checked by Register while
// holding the registry lock because those conflicts depend on current state.
func prepareChecks(target Target, checks []Checker) ([]preparedCheck, error) {
	if len(checks) == 0 {
		return nil, nil
	}

	prepared := make([]preparedCheck, 0, len(checks))
	seen := make(map[string]int, len(checks))
	var failures []error

	for index, check := range checks {
		if nilChecker(check) {
			failures = append(failures, NilCheckerError{
				Target: target,
				Index:  index,
			})
			continue
		}

		name := check.Name()
		if err := ValidateCheckName(name); err != nil {
			failures = append(failures, InvalidCheckNameError{
				Target: target,
				Index:  index,
				Name:   name,
				Err:    err,
			})
			continue
		}

		if previous, exists := seen[name]; exists {
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
			Index:   index,
			Name:    name,
			Checker: check,
		})
	}

	if len(failures) > 0 {
		return nil, errors.Join(failures...)
	}

	return prepared, nil
}

// nilChecker reports whether check is nil, including typed nil interface values.
//
// Typed nil values can occur when a nil pointer to a concrete checker type is
// assigned to Checker. Register rejects them because method calls on such values
// can panic later during evaluation.
func nilChecker(check Checker) bool {
	if check == nil {
		return true
	}

	value := reflect.ValueOf(check)
	switch value.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

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

// CheckerName returns the validated stable name for checker.
//
// CheckerName rejects nil and typed-nil checkers with ErrNilChecker, then reads
// Checker.Name exactly once and validates it with ValidateCheckName. The
// returned name is the checker identity that registries, check sets, and
// evaluators should use for the rest of that validation or execution boundary.
//
// If the checker name is invalid, CheckerName returns the invalid name together
// with the root check-name validation error. It does not execute Check and does
// not recover panics from Name because an unstable checker identity is a checker
// implementation bug.
func CheckerName(checker Checker) (string, error) {
	if nilChecker(checker) {
		return "", ErrNilChecker
	}

	name := checker.Name()
	if err := ValidateCheckName(name); err != nil {
		return name, err
	}

	return name, nil
}

// ValidateChecker validates the root Checker contract without executing Check.
//
// ValidateChecker is a convenience wrapper around CheckerName for callers that
// only need the validation outcome. Callers that need to keep using the checked
// name should call CheckerName so they do not read Checker.Name more than once.
func ValidateChecker(checker Checker) error {
	_, err := CheckerName(checker)
	return err
}

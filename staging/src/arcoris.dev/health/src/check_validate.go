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

// ValidateChecker validates the root Checker contract without executing Check.
//
// ValidateChecker rejects nil and typed-nil checkers with ErrNilChecker, then
// validates Checker.Name with ValidateCheckName. It does not recover panics from
// Name because an unstable checker identity is a checker implementation bug.
func ValidateChecker(checker Checker) error {
	if nilChecker(checker) {
		return ErrNilChecker
	}

	return ValidateCheckName(checker.Name())
}

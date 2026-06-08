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

package retrybudgetadmission

// requireBudget validates that a retry-budget adapter has a concrete budget
// before attempting to convert decisions into generic admission results.
//
// A missing budget is configuration misuse. It must remain a panic instead of a
// denied admission result so callers do not confuse broken wiring with normal
// retry-budget pressure.
func (a Admitter) requireBudget() {
	if a.budget == nil {
		panic(ErrNilRetryAdmitter)
	}
}

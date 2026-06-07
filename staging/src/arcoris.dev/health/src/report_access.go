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

// IsObserved reports whether r has a report-level observation timestamp.
func (r Report) IsObserved() bool {
	return !r.Observed.IsZero()
}

// Empty reports whether r contains no check results.
func (r Report) Empty() bool {
	return len(r.Checks) == 0
}

// Passed reports whether r.Status passes policy.
func (r Report) Passed(policy TargetPolicy) bool {
	return policy.Passes(r.Status)
}

// Failed reports whether r.Status fails policy.
func (r Report) Failed(policy TargetPolicy) bool {
	return policy.Fails(r.Status)
}

// Check returns the first check result with name.
//
// Reports produced from a CheckSet preserve resolver order and check names are
// unique per target, so a report produced by Evaluator can contain at most one
// result for a valid check name. The method still returns the first match
// defensively because Report is a plain caller-owned value.
func (r Report) Check(name string) (Result, bool) {
	for _, res := range r.Checks {
		if res.Name == name {
			return res, true
		}
	}

	return Result{}, false
}

// HasReason reports whether at least one check result has reason.
//
// HasReason performs exact reason matching. It intentionally does not interpret
// reason categories or status severity. Use Reason category helpers on individual
// results when callers need broader classification.
func (r Report) HasReason(reason Reason) bool {
	for _, res := range r.Checks {
		if res.HasReason(reason) {
			return true
		}
	}

	return false
}

// ChecksCopy returns a defensive copy of r.Checks.
func (r Report) ChecksCopy() []Result {
	if len(r.Checks) == 0 {
		return nil
	}

	checks := make([]Result, len(r.Checks))
	copy(checks, r.Checks)

	return checks
}

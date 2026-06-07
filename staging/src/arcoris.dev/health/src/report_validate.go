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

// IsValid reports whether r is structurally valid as a health report.
//
// The zero report is valid and represents "not evaluated yet." Every non-zero or
// evaluated report must use a concrete Target. This prevents caller-constructed
// reports from carrying concrete check results under TargetUnknown.
func (r Report) IsValid() bool {
	if !r.Target.IsValid() || !r.Status.IsValid() || r.Duration < 0 {
		return false
	}

	if r.Target == TargetUnknown {
		return r.Status == StatusUnknown &&
			r.Observed.IsZero() &&
			r.Duration == 0 &&
			len(r.Checks) == 0
	}

	for _, res := range r.Checks {
		if !res.IsValid() {
			return false
		}
	}

	return true
}

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

// IsValid reports whether r is structurally valid as a health result.
//
// A valid result has a known Status value, a valid Reason value, and a
// non-negative Duration. Name may be empty because the zero value is an unnamed
// StatusUnknown result and aggregators may fill the name from checker ownership.
// If Name is not empty, it must satisfy ValidCheckName.
func (r Result) IsValid() bool {
	if !r.Status.IsValid() || !r.Reason.IsValid() || r.Duration < 0 {
		return false
	}
	if r.Name != "" && !ValidCheckName(r.Name) {
		return false
	}

	return true
}

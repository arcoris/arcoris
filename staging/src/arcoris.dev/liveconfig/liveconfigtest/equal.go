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

package liveconfigtest

import (
	"maps"
	"slices"
)

// EqualConfig reports whether a and b are equal test configurations.
//
// EqualConfig compares slices and maps by contents. It is intended for tests of
// live configuration holders that accept an equality function to suppress no-op
// publications.
func EqualConfig(a, b Config) bool {
	if a.Name != b.Name || a.Version != b.Version || a.Enabled != b.Enabled || a.Timeout != b.Timeout {
		return false
	}
	return slices.Equal(a.Limits, b.Limits) && maps.Equal(a.Labels, b.Labels)
}

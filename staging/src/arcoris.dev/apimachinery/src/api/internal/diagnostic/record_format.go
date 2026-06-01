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

package diagnostic

import "strings"

// Format builds the common ARCORIS diagnostic string shape for r.
func (r Record[R]) Format(prefix string) string {
	parts := []string{prefix}

	if r.Path != "" {
		parts = append(parts, r.Path)
	}

	if r.Err != nil {
		parts = append(parts, r.Err.Error())
	}

	if r.Reason != "" {
		parts = append(parts, string(r.Reason))
	}

	if r.Detail != "" {
		parts = append(parts, r.Detail)
	}

	return strings.Join(parts, ": ")
}

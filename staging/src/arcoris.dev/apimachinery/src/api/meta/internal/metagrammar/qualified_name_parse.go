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

package metagrammar

import "strings"

// SplitQualifiedName splits "[prefix/]name" without allocating.
func SplitQualifiedName(s string) (prefix string, name string, ok bool) {
	if strings.Count(s, "/") > 1 {
		return "", "", false
	}

	prefix, name, found := strings.Cut(s, "/")
	if !found {
		return "", s, true
	}
	if prefix == "" || name == "" {
		return "", "", false
	}
	return prefix, name, true
}

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

package resource

import "testing"

func FuzzParseScope(f *testing.F) {
	for _, seed := range []string{
		"",
		scopeTextGlobal,
		scopeTextNamespaced,
		scopeTextInvalid,
		"cluster",
		"Namespaced",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		scope, err := ParseScope(input)
		if err != nil {
			return
		}

		requireNoError(t, scope.Validate())
		requireEqual(t, scope.String(), input)
	})
}

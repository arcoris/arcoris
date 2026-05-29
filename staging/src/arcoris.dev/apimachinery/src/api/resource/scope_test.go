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

func TestScopeHelpers(t *testing.T) {
	cases := []struct {
		name   string
		scope  Scope
		zero   bool
		valid  bool
		string string
	}{
		{name: "invalid", scope: ScopeInvalid, zero: true, string: scopeTextInvalid},
		{name: "global", scope: ScopeGlobal, valid: true, string: scopeTextGlobal},
		{name: "namespaced", scope: ScopeNamespaced, valid: true, string: scopeTextNamespaced},
		{name: "unknown", scope: Scope(99), string: scopeTextUnknown},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requireEqual(t, tc.scope.IsZero(), tc.zero)
			requireEqual(t, tc.scope.IsValid(), tc.valid)
			requireEqual(t, tc.scope.String(), tc.string)
		})
	}
}

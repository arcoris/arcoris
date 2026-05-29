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

func TestParseScope(t *testing.T) {
	valid := []struct {
		text string
		want Scope
	}{
		{text: scopeTextGlobal, want: ScopeGlobal},
		{text: scopeTextNamespaced, want: ScopeNamespaced},
	}

	for _, tc := range valid {
		t.Run(tc.text, func(t *testing.T) {
			got, err := ParseScope(tc.text)
			requireNoError(t, err)
			requireEqual(t, got, tc.want)
		})
	}
}

func TestParseScopeRejectsInvalidText(t *testing.T) {
	for _, text := range []string{"", "cluster", "Namespaced", "tenant"} {
		t.Run("reject/"+text, func(t *testing.T) {
			_, err := ParseScope(text)
			requireResourceError(t, err, ErrInvalidScope, pathScope, ErrorReasonInvalidScope)
		})
	}
}

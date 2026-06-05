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

package lexical

import "testing"

func TestASCIIHelpers(t *testing.T) {
	cases := []struct {
		name      string
		b         byte
		lower     bool
		upper     bool
		digit     bool
		alpha     bool
		alnum     bool
		labelChar bool
		labelEdge bool
	}{
		{name: "lowercase", b: 'a', lower: true, alpha: true, alnum: true, labelChar: true, labelEdge: true},
		{name: "uppercase", b: 'A', upper: true, alpha: true, alnum: true},
		{name: "digit", b: '7', digit: true, alnum: true, labelChar: true, labelEdge: true},
		{name: "hyphen", b: '-', labelChar: true},
		{name: "dot", b: '.'},
		{name: "underscore", b: '_'},
		{name: "slash", b: '/'},
		{name: "space", b: ' '},
		{name: "non ASCII byte", b: 0xd0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requireBool(t, IsASCIILower(tc.b), tc.lower, "IsASCIILower")
			requireBool(t, IsASCIIUpper(tc.b), tc.upper, "IsASCIIUpper")
			requireBool(t, IsASCIIDigit(tc.b), tc.digit, "IsASCIIDigit")
			requireBool(t, IsASCIIAlpha(tc.b), tc.alpha, "IsASCIIAlpha")
			requireBool(t, IsASCIIAlnum(tc.b), tc.alnum, "IsASCIIAlnum")
			requireBool(t, IsDNS1123LabelChar(tc.b), tc.labelChar, "IsDNS1123LabelChar")
			requireBool(t, IsDNS1123LabelEdge(tc.b), tc.labelEdge, "IsDNS1123LabelEdge")
		})
	}
}

func requireBool(t *testing.T, got bool, want bool, name string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}

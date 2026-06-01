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

package fieldpath

import "testing"

func TestLiteralKindString(t *testing.T) {
	testCases := []struct {
		name string
		kind LiteralKind
		want string
	}{
		{name: "invalid", kind: LiteralInvalid, want: "invalid"},
		{name: "bool", kind: LiteralBool, want: "bool"},
		{name: "integer", kind: LiteralInteger, want: "integer"},
		{name: "string", kind: LiteralString, want: "string"},
		{name: "unknown", kind: LiteralKind(255), want: "unknown"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			requireEqual(t, testCase.kind.String(), testCase.want)
		})
	}
}

func TestLiteralKindIsValid(t *testing.T) {
	requireEqual(t, LiteralInvalid.IsValid(), false)
	requireEqual(t, LiteralBool.IsValid(), true)
	requireEqual(t, LiteralInteger.IsValid(), true)
	requireEqual(t, LiteralString.IsValid(), true)
	requireEqual(t, LiteralKind(255).IsValid(), false)
}

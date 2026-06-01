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

func TestPathParserParseLiteralBool(t *testing.T) {
	testCases := []struct {
		name string
		text string
		want Literal
	}{
		{name: "true", text: "true", want: BoolLiteral(true)},
		{name: "false", text: "false", want: BoolLiteral(false)},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			p := newPathParser(testCase.text)

			value, err := p.parseLiteral()
			requireNoError(t, err)

			requireEqual(t, value.Equal(testCase.want), true)
			requireEqual(t, p.done(), true)
		})
	}
}

func TestPathParserParseLiteralInteger(t *testing.T) {
	testCases := []struct {
		name string
		text string
		want Literal
	}{
		{name: "signed", text: "-7", want: Int64Literal(-7)},
		{name: "unsigned", text: "42", want: Uint64Literal(42)},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			p := newPathParser(testCase.text)

			value, err := p.parseLiteral()
			requireNoError(t, err)

			requireEqual(t, value.Equal(testCase.want), true)
			requireEqual(t, p.done(), true)
		})
	}
}

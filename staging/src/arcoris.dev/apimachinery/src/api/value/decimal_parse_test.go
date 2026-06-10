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

package value

import "testing"

func TestParseDecimalCanonicalizesValidInput(t *testing.T) {
	tests := []struct {
		input string
		text  string
	}{
		{input: "0", text: "0"},
		{input: "-0", text: "0"},
		{input: "123", text: "123"},
		{input: "-123", text: "-123"},
		{input: "123.45", text: "123.45"},
		{input: "-0.01", text: "-0.01"},
		{input: "0.001", text: "0.001"},
		{input: "001", text: "1"},
		{input: "001.20", text: "1.20"},
		{input: "-000.010", text: "-0.010"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			decimal, err := ParseDecimal(tt.input)
			requireNoError(t, err)
			requireEqual(t, decimal.String(), tt.text)
		})
	}
}

func TestMustParseDecimalPanicsOnMalformedInput(t *testing.T) {
	requirePanic(t, func() {
		MustParseDecimal("1e3")
	})
}

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

func TestNewDecimalRejectsMalformedInput(t *testing.T) {
	tests := []string{
		"",
		"+",
		"+1",
		"-",
		".",
		".1",
		"-.",
		"1.2.3",
		"abc",
		"1e3",
		"1E-3",
		"NaN",
		"Inf",
		" 1",
		"1 ",
		"1_000",
		"١٢٣",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := NewDecimal(input)
			requireValueError(t, err, ErrInvalidDecimal, pathDecimal, ErrorReasonInvalidDecimal)
			requireErrorIs(t, err, ErrInvalidValue)
		})
	}
}

func TestNewDecimalRejectsTrailingDecimalPoint(t *testing.T) {
	tests := []string{
		"1.",
		"0.",
		"-1.",
		"-0.",
		"001.",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := NewDecimal(input)

			requireValueError(
				t,
				err,
				ErrInvalidDecimal,
				pathDecimal,
				ErrorReasonInvalidDecimal,
			)
		})
	}
}

func TestCanonicalDecimalPartsPreservesFractionalScale(t *testing.T) {
	coefficient, scale := canonicalDecimalParts("001", "20")

	requireEqual(t, coefficient, "120")
	requireEqual(t, scale, uint32(2))
}

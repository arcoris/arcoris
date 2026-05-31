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

func TestDecimalAccessorsAndString(t *testing.T) {
	tests := []struct {
		decimal     Decimal
		text        string
		negative    bool
		coefficient string
		scale       uint32
	}{
		{decimal: Decimal{coefficient: "0"}, text: "0", coefficient: "0"},
		{decimal: Decimal{coefficient: "123"}, text: "123", coefficient: "123"},
		{
			decimal:     Decimal{negative: true, coefficient: "123"},
			text:        "-123",
			negative:    true,
			coefficient: "123",
		},
		{
			decimal:     Decimal{coefficient: "12345", scale: 2},
			text:        "123.45",
			coefficient: "12345",
			scale:       2,
		},
		{
			decimal:     Decimal{negative: true, coefficient: "1", scale: 2},
			text:        "-0.01",
			negative:    true,
			coefficient: "1",
			scale:       2,
		},
		{
			decimal:     Decimal{coefficient: "1", scale: 3},
			text:        "0.001",
			coefficient: "1",
			scale:       3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			requireEqual(t, tt.decimal.String(), tt.text)
			requireEqual(t, tt.decimal.IsNegative(), tt.negative)
			requireEqual(t, tt.decimal.Coefficient(), tt.coefficient)
			requireEqual(t, tt.decimal.Scale(), tt.scale)
		})
	}
}

func TestDecimalEqual(t *testing.T) {
	decimal := Decimal{coefficient: "120", scale: 2}

	requireEqual(t, decimal.Equal(Decimal{coefficient: "120", scale: 2}), true)
	requireEqual(t, decimal.Equal(Decimal{coefficient: "12", scale: 1}), false)
}

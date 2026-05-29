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

package types

import (
	"math"
	"testing"
)

func TestFloat32ValidateRejectsInvalidRules(t *testing.T) {
	tests := []Type{
		Float32().Range(10, 1).Type(),
		Float32().Min(1).Enum(0).Type(),
		Float32().Max(1).Enum(2).Type(),
		Float32().Enum(1, 1).Type(),
		Float32().Enum(float32(math.Inf(1))).Type(),
		Float32().Enum(float32(math.NaN())).Type(),
	}
	for _, typ := range tests {
		requireErrorIs(t, ValidateType(typ, nil), ErrInvalidType)
	}
}

func TestFloat32ValidateEnumDuplicatesUseNumericEquality(t *testing.T) {
	negativeZero := float32(math.Copysign(0, -1))

	requireErrorIs(t, ValidateType(Float32().Enum(0, negativeZero).Type(), nil), ErrInvalidType)
}

func TestInvalidFloat32(t *testing.T) {
	requireEqual(t, invalidFloat32(float32(math.NaN())), true)
	requireEqual(t, invalidFloat32(float32(math.Inf(1))), true)
	requireEqual(t, invalidFloat32(float32(math.Inf(-1))), true)
	requireEqual(t, invalidFloat32(1), false)
}

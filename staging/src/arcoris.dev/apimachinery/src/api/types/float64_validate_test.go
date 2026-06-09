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

func TestFloat64ValidateRejectsInvalidRules(t *testing.T) {
	tests := []Descriptor{
		Float64().Range(10, 1).Descriptor(),
		Float64().Min(1).Enum(0).Descriptor(),
		Float64().Max(1).Enum(2).Descriptor(),
		Float64().Enum(1, 1).Descriptor(),
		Float64().Enum(math.Inf(1)).Descriptor(),
		Float64().Enum(math.NaN()).Descriptor(),
	}
	for _, desc := range tests {
		requireErrorIs(t, ValidateLocal(desc), ErrInvalidDescriptor)
	}
}

func TestFloat64ValidateEnumDuplicatesUseNumericEquality(t *testing.T) {
	negativeZero := math.Copysign(0, -1)

	requireErrorIs(t, ValidateLocal(Float64().Enum(0, negativeZero).Descriptor()), ErrInvalidDescriptor)
}

func TestInvalidFloat64(t *testing.T) {
	requireEqual(t, invalidFloat64(math.NaN()), true)
	requireEqual(t, invalidFloat64(math.Inf(1)), true)
	requireEqual(t, invalidFloat64(math.Inf(-1)), true)
	requireEqual(t, invalidFloat64(1), false)
}

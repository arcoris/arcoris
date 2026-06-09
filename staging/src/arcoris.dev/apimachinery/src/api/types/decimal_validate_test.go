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

import "testing"

func TestDecimalValidateRejectsInvalidRules(t *testing.T) {
	tests := []Descriptor{
		Decimal().Precision(0).Descriptor(),
		Decimal().Scale(-1).Descriptor(),
		Decimal().Precision(2).Scale(3).Descriptor(),
	}
	for _, desc := range tests {
		requireErrorIs(t, ValidateLocal(desc), ErrInvalidDescriptor)
	}
}

func TestDecimalValidateAllowsScaleWithoutPrecision(t *testing.T) {
	desc := Decimal().Scale(2).Descriptor()

	requireNoError(t, ValidateLocal(desc))

	view := requireDecimalView(t, desc)
	scale, ok := view.Scale()
	requireEqual(t, ok, true)
	requireEqual(t, scale, 2)
	_, ok = view.Precision()
	requireEqual(t, ok, false)
}

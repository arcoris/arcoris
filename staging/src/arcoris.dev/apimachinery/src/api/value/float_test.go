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

import (
	"math"
	"testing"
)

func TestFloatValueAcceptsFiniteValues(t *testing.T) {
	value, err := FloatValue(1.5)
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindFloat)
	requireEqual(t, value.floatValue, 1.5)
}

func TestFloatValueRejectsNonFiniteValues(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{name: "NaN", value: math.NaN()},
		{name: "+Inf", value: math.Inf(1)},
		{name: "-Inf", value: math.Inf(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FloatValue(tt.value)
			requireValueError(t, err, ErrInvalidFloat, pathFloat, ErrorReasonInvalidFloat)
			requireErrorIs(t, err, ErrInvalidValue)
		})
	}
}

func TestMustFloatValuePanicsOnInvalidInput(t *testing.T) {
	requirePanic(t, func() {
		MustFloatValue(math.NaN())
	})
}

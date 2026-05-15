/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package fixedwindow

import (
	"errors"
	"math"
	"testing"
)

func TestValidateRatio(t *testing.T) {
	tests := []struct {
		name  string
		ratio float64
		want  error
	}{
		{name: "zero", ratio: 0},
		{name: "fraction", ratio: 0.2},
		{name: "one", ratio: 1},
		{name: "negative", ratio: -0.1, want: ErrInvalidRatio},
		{name: "greater than one", ratio: 1.1, want: ErrInvalidRatio},
		{name: "nan", ratio: math.NaN(), want: ErrInvalidRatio},
		{name: "positive infinity", ratio: math.Inf(1), want: ErrInvalidRatio},
		{name: "negative infinity", ratio: math.Inf(-1), want: ErrInvalidRatio},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRatio(tt.ratio)
			if !errors.Is(err, tt.want) {
				t.Fatalf("validateRatio() error = %v, want %v", err, tt.want)
			}
		})
	}
}

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

package codecselection

// Weight is an integer encode preference weight.
//
// Weight is deliberately not a floating-point HTTP q-value. Future protocol
// parsers can map q-values into this deterministic internal scale.
type Weight uint16

const (
	// WeightDefault is the default highest ordinary preference weight.
	WeightDefault Weight = 1000

	// WeightMax is the largest supported preference weight.
	WeightMax Weight = 1000
)

// NewWeight validates value and returns a Weight.
func NewWeight(value int) (Weight, error) {
	weight := Weight(value)
	if err := weight.Validate(); err != nil {
		return 0, err
	}

	return weight, nil
}

// MustWeight returns a valid Weight or panics when value is invalid.
func MustWeight(value int) Weight {
	weight, err := NewWeight(value)
	if err != nil {
		panic(err)
	}

	return weight
}

// IsZero reports whether w is absent.
func (w Weight) IsZero() bool {
	return w == 0
}

// Validate checks that w is within the supported preference range.
func (w Weight) Validate() error {
	if w == 0 || w > WeightMax {
		return errorfAt(
			pathWeight,
			ErrInvalidPreference,
			ErrorReasonInvalidPreference,
			"preference weight %d must be between 1 and %d",
			w,
			WeightMax,
		)
	}

	return nil
}

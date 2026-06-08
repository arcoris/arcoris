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

import "testing"

func TestNewWeight(t *testing.T) {
	weight, err := NewWeight(750)
	requireNoError(t, err)

	if weight != Weight(750) {
		t.Fatalf("weight = %d; want 750", weight)
	}
	if weight.IsZero() {
		t.Fatalf("IsZero() = true; want false")
	}
}

func TestWeightValidateRejectsOutOfRangeValues(t *testing.T) {
	for _, value := range []Weight{0, WeightMax + 1} {
		t.Run("value", func(t *testing.T) {
			err := value.Validate()

			requireErrorIs(t, err, ErrInvalidPreference)
			requireSelectionError(t, err, "codecselection.weight", ErrorReasonInvalidPreference)
		})
	}
}

func TestNewWeightRejectsNegativeValue(t *testing.T) {
	_, err := NewWeight(-1)

	requireErrorIs(t, err, ErrInvalidPreference)
	requireSelectionError(t, err, "codecselection.weight", ErrorReasonInvalidPreference)
}

func TestMustWeightPanicsOnInvalidInput(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("MustWeight did not panic")
		}
	}()

	_ = MustWeight(0)
}

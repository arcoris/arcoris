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

// validateDecimal checks TypeDecimal precision and scale descriptor rules.
//
// Decimal deliberately has no min/max or enum rules in this package. Exact
// decimal values need a future value representation design before those rules
// can be added without encoding policy leaking into descriptors.
func validateDecimal(t Type, path string) error {
	if t.decimal.precision.set && t.decimal.precision.value <= 0 {
		return typeError(path+".precision", ErrInvalidType)
	}
	if t.decimal.scale.set && t.decimal.scale.value < 0 {
		return typeError(path+".scale", ErrInvalidType)
	}
	if t.decimal.precision.set && t.decimal.scale.set && t.decimal.scale.value > t.decimal.precision.value {
		return typeError(path+".scale", ErrInvalidType)
	}
	return nil
}

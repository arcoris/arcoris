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

// DecimalView exposes read-only TypeDecimal payload data.
type DecimalView struct {
	// payload is a detached copy of the decimal descriptor payload.
	payload decimalPayload
}

// Decimal returns a decimal view when t is TypeDecimal.
func (t Type) Decimal() (DecimalView, bool) {
	if t.code != TypeDecimal {
		return DecimalView{}, false
	}

	return DecimalView{payload: cloneDecimalPayload(t.decimal)}, true
}

// Precision returns the decimal precision rule.
func (v DecimalView) Precision() (int, bool) {
	return v.payload.precision.value, v.payload.precision.set
}

// Scale returns the decimal scale rule.
func (v DecimalView) Scale() (int, bool) {
	return v.payload.scale.value, v.payload.scale.set
}

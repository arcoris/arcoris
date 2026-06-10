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

// AsDecimal returns the decimal payload when v is KindDecimal.
//
// For every other kind, AsDecimal returns the zero Decimal and ok=false. The zero
// Decimal is a valid payload value, so callers must check ok.
func (v Value) AsDecimal() (Decimal, bool) {
	if v.kind != KindDecimal {
		return Decimal{}, false
	}

	return v.decimalValue, true
}

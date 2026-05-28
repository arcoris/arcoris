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

// Decimal starts a field descriptor for exact decimal values.
//
// Decimal descriptors currently record precision and scale only. Concrete
// decimal value representation and decimal min/max semantics are intentionally
// deferred to a later value-model design.
//
// Field builder flow:
//
//	Field("price").Decimal().
//		Optional().
//		Nullable().
//		Precision(12).
//		Scale(2).
//		Description("Exact decimal price.")
func (b FieldBuilder) Decimal() DecimalField {
	return DecimalField{field: b.state(), typ: Decimal()}
}

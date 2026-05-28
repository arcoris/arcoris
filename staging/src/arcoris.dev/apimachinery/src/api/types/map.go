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

// MapType builds descriptors for dynamic string-keyed maps.
//
// MapType is for dictionaries with dynamic keys and one shared value type. It
// is intentionally separate from ObjectType, which models fixed schema fields.
// Only string keys are supported in this design pass.
type MapType struct {
	// header stores the descriptor family and cross-family flags under construction.
	header typeHeader
	// payload stores the exact map shape under construction.
	payload mapPayload
}

// MapOf returns a string-keyed map descriptor builder for value.
//
// A nil TypeExpr is recorded as an invalid zero value descriptor so
// ValidateType can classify the error at map.value. The builder itself stays
// allocation-light and panic-free.
//
// Typical reusable declaration:
//
//	labelValue := String()
//	labelValue = labelValue.MinLen(1)
//
//	labelsType := MapOf(
//		labelValue,
//	).
//		MaxLen(64)
func MapOf(value TypeExpr) MapType {
	valueType := typeFromExpr(value)
	return MapType{
		header: newHeader(TypeMap),
		payload: mapPayload{
			key:   MapKeyString,
			value: &valueType,
		},
	}
}

// Nullable returns a map descriptor that admits null values.
func (t MapType) Nullable() MapType { t.header = t.header.withNullable(); return t }

// MinLen sets the inclusive minimum number of map entries.
//
// The limit is structural metadata only. Concrete map entry counts are checked
// by future value-validation layers.
func (t MapType) MinLen(n int) MapType { t.payload.minLen = intLimit{n, true}; return t }

// MaxLen sets the inclusive maximum number of map entries.
//
// The limit uses intLimit so an explicit zero maximum can be represented
// without a pointer allocation.
func (t MapType) MaxLen(n int) MapType { t.payload.maxLen = intLimit{n, true}; return t }

// Type returns a detached Type descriptor.
func (t MapType) Type() Type {
	out := typeFromHeader(t.header)
	out.mapType = cloneMapPayload(t.payload)
	return out
}

// typeExpr marks MapType as a sealed TypeExpr implementation.
func (t MapType) typeExpr() {}

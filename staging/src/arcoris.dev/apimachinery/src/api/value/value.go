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

import "time"

// Value is the closed concrete ARCORIS API payload value algebra.
//
// A Value stores actual payload data, not a descriptor, constraint, or decoded
// Go struct wrapper. The zero value is intentionally invalid and represents
// missing initialization. Use Null to represent an explicit API null literal.
//
// All payload slots are private so constructors can preserve numeric
// invariants, container ordering, and copy-on-boundary behavior. Public
// accessors return payloads only for the matching Kind and never expose mutable
// internal slices.
type Value struct {
	// kind selects the active payload slot.
	kind Kind

	// boolValue stores KindBool payload data.
	boolValue bool

	// stringValue stores KindString payload data.
	stringValue string
	// bytesValue stores KindBytes payload data and is always owned by Value.
	bytesValue []byte

	// integerValue stores KindInteger payload data.
	integerValue Integer
	// floatValue stores KindFloat payload data.
	floatValue float64
	// decimalValue stores KindDecimal payload data.
	decimalValue Decimal

	// timestampValue stores KindTimestamp payload data with monotonic time stripped.
	timestampValue time.Time
	// dateValue stores KindDate payload data.
	dateValue Date
	// timeOfDayValue stores KindTimeOfDay payload data.
	timeOfDayValue TimeOfDay
	// durationValue stores KindDuration payload data.
	durationValue time.Duration

	// objectValue stores KindObject payload data.
	objectValue objectPayload
	// listValue stores KindList payload data.
	listValue listPayload
	// mapValue stores KindMap payload data.
	mapValue mapPayload
}

// Kind returns the concrete payload category stored in v.
//
// The zero Value reports KindInvalid. Kind is a value discriminator only; it
// does not carry descriptor information such as integer width, field
// requirements, nullability, or collection constraints.
func (v Value) Kind() Kind {
	return v.kind
}

// IsZero reports whether v is the invalid uninitialized zero value.
//
// IsZero is deliberately not the same as IsNull. Null is a real API payload
// value, while the zero Value means construction did not happen.
func (v Value) IsZero() bool {
	return v.kind == KindInvalid
}

// IsNull reports whether v stores the explicit null literal.
//
// Explicit null is a payload value that future descriptor-aware validation may
// accept or reject according to field nullability rules.
func (v Value) IsNull() bool {
	return v.kind == KindNull
}

// IsScalar reports whether v stores a non-composite payload category.
//
// Scalars include primitive, numeric, and temporal values. The invalid zero
// Value is not considered scalar.
func (v Value) IsScalar() bool {
	return v.kind.IsPrimitive()
}

// IsComposite reports whether v stores nested payload values.
//
// Composite values are object, list, and map payloads. The method does not
// inspect nested values because constructors already reject invalid children.
func (v Value) IsComposite() bool {
	return v.kind.IsComposite()
}

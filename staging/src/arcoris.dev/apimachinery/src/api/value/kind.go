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

// Kind identifies the concrete category stored in a Value.
//
// Kind is a value discriminator, not a descriptor. It deliberately omits exact
// numeric widths, nullability, type references, required/optional field state,
// unknown-field policy, and collection semantics because those are constraints
// owned by descriptor packages and future validation layers.
type Kind uint8

const (
	// KindInvalid is the zero value and is never a valid concrete payload.
	KindInvalid Kind = iota

	// KindNull represents the explicit null literal.
	KindNull
	// KindBool represents a boolean payload.
	KindBool
	// KindString represents Go string payload data.
	KindString
	// KindBytes represents arbitrary binary payload data.
	KindBytes

	// KindInteger represents one exact integer in the int64/uint64 union.
	KindInteger
	// KindFloat represents one finite float64 payload.
	KindFloat
	// KindDecimal represents one exact base-10 decimal payload.
	KindDecimal

	// KindTimestamp represents one time.Time timestamp payload.
	KindTimestamp
	// KindDate represents a calendar date without time-of-day or timezone.
	KindDate
	// KindTimeOfDay represents a wall-clock time without date or timezone.
	KindTimeOfDay
	// KindDuration represents one elapsed interval payload.
	KindDuration

	// KindObject represents a schema-defined record value with named fields.
	KindObject
	// KindList represents an ordered sequence of values.
	KindList
	// KindMap represents a dynamic key/value mapping.
	KindMap
)

// IsValid reports whether k identifies a supported concrete value category.
//
// KindInvalid is reserved for the zero/uninitialized Value and is not a valid
// payload category.
func (k Kind) IsValid() bool {
	return k > KindInvalid && k <= KindMap
}

// String returns a stable diagnostic name for k.
//
// The returned text is intended for diagnostics and tests. It is not a
// wire-format contract and should not be used as a codec discriminator.
func (k Kind) String() string {
	switch k {
	case KindInvalid:
		return "invalid"
	case KindNull:
		return "null"
	case KindBool:
		return "bool"
	case KindString:
		return "string"
	case KindBytes:
		return "bytes"
	case KindInteger:
		return "integer"
	case KindFloat:
		return "float"
	case KindDecimal:
		return "decimal"
	case KindTimestamp:
		return "timestamp"
	case KindDate:
		return "date"
	case KindTimeOfDay:
		return "timeOfDay"
	case KindDuration:
		return "duration"
	case KindObject:
		return "object"
	case KindList:
		return "list"
	case KindMap:
		return "map"
	default:
		return "unknown"
	}
}

// IsPrimitive reports whether k stores one scalar payload value.
//
// Primitive here means "not object/list/map" for the value algebra. It includes
// numbers and temporal values even though descriptors may classify those more
// narrowly.
func (k Kind) IsPrimitive() bool {
	switch k {
	case KindNull,
		KindBool,
		KindString,
		KindBytes,
		KindInteger,
		KindFloat,
		KindDecimal,
		KindTimestamp,
		KindDate,
		KindTimeOfDay,
		KindDuration:
		return true
	default:
		return false
	}
}

// IsNumber reports whether k stores a numeric payload value.
//
// Integer, float, and decimal are separate value categories because they carry
// different portability and precision guarantees.
func (k Kind) IsNumber() bool {
	return k == KindInteger || k == KindFloat || k == KindDecimal
}

// IsTemporal reports whether k stores time-related payload data.
//
// Temporal values are concrete data values only; this package does not impose
// ordering, clock-source, or timezone policy.
func (k Kind) IsTemporal() bool {
	return k == KindTimestamp || k == KindDate || k == KindTimeOfDay || k == KindDuration
}

// IsComposite reports whether k stores nested payload values.
//
// Composite values preserve caller order and expose read-only views rather than
// mutable maps or slices.
func (k Kind) IsComposite() bool {
	return k == KindObject || k == KindList || k == KindMap
}

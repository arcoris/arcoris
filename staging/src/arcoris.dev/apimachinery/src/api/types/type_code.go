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

// TypeCode identifies the structural category of an API value type.
//
// TypeCode is the discriminator used by Type. It describes logical API value
// shapes, not Go runtime storage and not a particular wire format.
//
// A JSON codec, OpenAPI exporter, binary codec, code generator, validator,
// pruning engine, defaulting engine, and future patch/apply engine may all
// interpret the same TypeCode differently, but they MUST preserve the same
// logical value contract represented by the containing Type.
//
// TypeCode deliberately avoids Go platform-sized and runtime-only concepts such
// as int, uint, uintptr, error, fmt.Stringer, reflection values, and encoder
// namespaces. Those concepts are useful for logging/event fields, but they are
// not stable distributed API schema types.
type TypeCode uint8

const (
	// TypeInvalid is the zero value and is never a valid structural type.
	//
	// It exists to make missing initialization visible. Constructors, type
	// validators, registries, code generators, and exporters MUST reject
	// TypeInvalid.
	TypeInvalid TypeCode = iota

	// TypeNull represents the null literal.
	//
	// TypeNull means that the only valid value is null. It does not make another
	// type nullable. Nullable fields are represented by field/type nullability
	// semantics, not by replacing the field type with TypeNull.
	//
	// A nullable string is a string Type with nullable semantics, not a TypeNull.
	TypeNull

	// TypeBool represents a boolean value.
	TypeBool

	// TypeString represents UTF-8 text.
	//
	// TypeString is for textual data. Binary payloads MUST use TypeBytes.
	// Semantic string types such as names, namespaces, UUIDs, resource names,
	// versions, hostnames, cron expressions, or other domain identifiers SHOULD
	// be represented as named custom types over TypeString.
	TypeString

	// TypeBytes represents arbitrary binary data.
	//
	// TypeBytes does not define a wire encoding. JSON codecs may encode bytes as
	// base64 strings, while binary codecs may write the raw byte sequence.
	TypeBytes

	// TypeInt8 represents an int8 value.
	TypeInt8

	// TypeInt16 represents an int16 value.
	TypeInt16

	// TypeInt32 represents an int32 value.
	TypeInt32

	// TypeInt64 represents an int64 value.
	TypeInt64

	// TypeUint8 represents a uint8 value.
	TypeUint8

	// TypeUint16 represents a uint16 value.
	TypeUint16

	// TypeUint32 represents a uint32 value.
	TypeUint32

	// TypeUint64 represents a uint64 value.
	//
	// Codecs and generators MUST account for target formats that cannot
	// precisely represent all uint64 values, such as JSON clients backed by
	// IEEE-754 double precision numbers.
	TypeUint64

	// TypeFloat32 represents an IEEE-754 binary32 floating-point number.
	//
	// NaN and infinities SHOULD NOT be considered valid portable API values
	// unless a future Type rule explicitly permits them.
	TypeFloat32

	// TypeFloat64 represents an IEEE-754 binary64 floating-point number.
	//
	// NaN and infinities SHOULD NOT be considered valid portable API values
	// unless a future Type rule explicitly permits them.
	TypeFloat64

	// TypeDecimal represents an exact base-10 numeric value.
	//
	// TypeDecimal is intended for values where binary floating-point semantics
	// are not acceptable, such as money, exact quotas, ratios, limits, or
	// human-authored decimal configuration. Precision, scale, rounding policy,
	// and encoding are Type rules or codec concerns, not TypeCode concerns.
	TypeDecimal

	// TypeTimestamp represents an absolute point in time.
	//
	// TypeTimestamp is a semantic API type. It does not prescribe whether a
	// codec stores the value as RFC 3339 text, Unix nanoseconds, Unix
	// milliseconds, or another canonical representation.
	TypeTimestamp

	// TypeDate represents a calendar date without a time of day.
	//
	// TypeDate is distinct from TypeTimestamp. It SHOULD be used for values such
	// as billing dates, schedule dates, report dates, or other date-only
	// concepts.
	TypeDate

	// TypeTime represents a time of day without a calendar date.
	//
	// TypeTime is distinct from TypeTimestamp and TypeDuration. It SHOULD be
	// used for values such as daily schedule boundaries, local opening times, or
	// other time-of-day concepts.
	TypeTime

	// TypeDuration represents an elapsed interval.
	//
	// TypeDuration is a semantic API type. It does not prescribe whether a codec
	// stores the value as nanoseconds, ISO-8601 duration text, Go-style duration
	// text, or another canonical representation.
	TypeDuration

	// TypeObject represents a structural object with a fixed set of named fields.
	//
	// TypeObject is for schema-defined records. Dynamic key/value dictionaries
	// MUST use TypeMap. Unknown-field behavior is controlled by object Type
	// rules, not by TypeCode.
	TypeObject

	// TypeList represents an ordered sequence of values of one element type.
	//
	// Merge/apply behavior, such as atomic list, set-like list, or map-like
	// list, is controlled by list Type rules, not by TypeCode.
	TypeList

	// TypeMap represents a dynamic key/value mapping.
	//
	// TypeMap is distinct from TypeObject. TypeObject has fixed schema fields;
	// TypeMap has dynamic keys and a common value type. Key restrictions are
	// controlled by map Type rules.
	TypeMap

	// TypeRef represents a reference to a named custom type definition.
	//
	// TypeRef is the extension mechanism for reusable semantic types. A TypeRef
	// MUST resolve through an owner-created resolver. It MUST NOT refer to
	// arbitrary Go implementations.
	TypeRef
)

// IsValid reports whether c identifies a supported structural type category.
func (c TypeCode) IsValid() bool {
	return c > TypeInvalid && c <= TypeRef
}

// String returns the stable diagnostic name of c.
//
// String is intended for diagnostics, tests, and developer-facing output. It is
// not a wire-format contract.
func (c TypeCode) String() string {
	switch c {
	case TypeInvalid:
		return "invalid"
	case TypeNull:
		return "null"
	case TypeBool:
		return "bool"
	case TypeString:
		return "string"
	case TypeBytes:
		return "bytes"
	case TypeInt8:
		return "int8"
	case TypeInt16:
		return "int16"
	case TypeInt32:
		return "int32"
	case TypeInt64:
		return "int64"
	case TypeUint8:
		return "uint8"
	case TypeUint16:
		return "uint16"
	case TypeUint32:
		return "uint32"
	case TypeUint64:
		return "uint64"
	case TypeFloat32:
		return "float32"
	case TypeFloat64:
		return "float64"
	case TypeDecimal:
		return "decimal"
	case TypeTimestamp:
		return "timestamp"
	case TypeDate:
		return "date"
	case TypeTime:
		return "time"
	case TypeDuration:
		return "duration"
	case TypeObject:
		return "object"
	case TypeList:
		return "list"
	case TypeMap:
		return "map"
	case TypeRef:
		return "ref"
	default:
		return "unknown"
	}
}

// IsPrimitive reports whether c is a leaf value descriptor.
//
// Primitive includes null, bool, string, bytes, numbers, and temporal values.
// It excludes composites and references. TypeRef is a named indirection rather
// than a primitive payload.
func (c TypeCode) IsPrimitive() bool {
	return c == TypeNull ||
		c == TypeBool ||
		c == TypeString ||
		c == TypeBytes ||
		c.IsNumber() ||
		c.IsTemporal()
}

// IsNumber reports whether c belongs to the numeric descriptor category.
//
// Decimal is included because it is a number even though it is not an integer or
// binary floating-point type. Decimal exact values intentionally do not have
// min/max support in this package until a decimal value representation exists.
func (c TypeCode) IsNumber() bool {
	return c.IsInteger() || c.IsFloat() || c == TypeDecimal
}

// IsInteger reports whether c is a fixed-width int or uint descriptor.
func (c TypeCode) IsInteger() bool {
	return c.IsInt() || c.IsUint()
}

// IsInt reports whether c is a fixed-width int descriptor.
func (c TypeCode) IsInt() bool {
	return c >= TypeInt8 && c <= TypeInt64
}

// IsUint reports whether c is a fixed-width uint descriptor.
func (c TypeCode) IsUint() bool {
	return c >= TypeUint8 && c <= TypeUint64
}

// IsFloat reports whether c is a binary floating-point type.
func (c TypeCode) IsFloat() bool {
	return c == TypeFloat32 || c == TypeFloat64
}

// IsTemporal reports whether c belongs to the semantic temporal kind.
//
// Temporal TypeCode values do not prescribe a codec representation. RFC 3339,
// ISO-8601 durations, Unix timestamps, and Go-style duration text are encoding
// choices for future codec packages.
func (c TypeCode) IsTemporal() bool {
	return c >= TypeTimestamp && c <= TypeDuration
}

// IsComposite reports whether c directly contains nested type descriptors.
//
// TypeRef is intentionally excluded. It points at a Resolver definition but is
// not itself a composite payload.
func (c TypeCode) IsComposite() bool {
	return c == TypeObject || c == TypeList || c == TypeMap
}

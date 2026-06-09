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

// DescriptorKind identifies the structural category of an API value descriptor.
//
// DescriptorKind is the discriminator used by Descriptor. It describes logical
// API value shapes, not Go runtime storage and not a particular wire format.
//
// A JSON codec, OpenAPI exporter, binary codec, code generator, validator,
// pruning engine, defaulting engine, and future patch/apply engine may all
// interpret the same DescriptorKind differently, but they MUST preserve the same
// logical value contract represented by the containing Descriptor.
//
// DescriptorKind deliberately avoids Go platform-sized and runtime-only
// concepts such as int, uint, uintptr, error, fmt.Stringer, reflection values,
// and encoder namespaces. Those concepts are useful for logging/event fields,
// but they are not stable distributed API schema descriptors.
type DescriptorKind uint8

const (
	// DescriptorInvalid is the zero value and is never a valid structural descriptor.
	//
	// It exists to make missing initialization visible. Constructors, descriptor
	// validators, registries, code generators, and exporters MUST reject
	// DescriptorInvalid.
	DescriptorInvalid DescriptorKind = iota

	// DescriptorNull represents the null literal.
	//
	// DescriptorNull means that the only valid value is null. It does not make another
	// descriptor nullable. Nullable fields are represented by field descriptor nullability
	// semantics, not by replacing the field descriptor with DescriptorNull.
	//
	// A nullable string is a string Descriptor with nullable semantics, not a DescriptorNull.
	DescriptorNull

	// DescriptorBool represents a boolean value.
	DescriptorBool

	// DescriptorString represents UTF-8 text.
	//
	// DescriptorString is for textual data. Binary payloads MUST use DescriptorBytes.
	// Semantic string descriptors such as names, namespaces, UUIDs, resource names,
	// versions, hostnames, cron expressions, or other domain identifiers SHOULD
	// be represented as named custom descriptors over DescriptorString.
	DescriptorString

	// DescriptorBytes represents arbitrary binary data.
	//
	// DescriptorBytes does not define a wire encoding. JSON codecs may encode bytes as
	// base64 strings, while binary codecs may write the raw byte sequence.
	DescriptorBytes

	// DescriptorInt8 represents an int8 value.
	DescriptorInt8

	// DescriptorInt16 represents an int16 value.
	DescriptorInt16

	// DescriptorInt32 represents an int32 value.
	DescriptorInt32

	// DescriptorInt64 represents an int64 value.
	DescriptorInt64

	// DescriptorUint8 represents a uint8 value.
	DescriptorUint8

	// DescriptorUint16 represents a uint16 value.
	DescriptorUint16

	// DescriptorUint32 represents a uint32 value.
	DescriptorUint32

	// DescriptorUint64 represents a uint64 value.
	//
	// Codecs and generators MUST account for target formats that cannot
	// precisely represent all uint64 values, such as JSON clients backed by
	// IEEE-754 double precision numbers.
	DescriptorUint64

	// DescriptorFloat32 represents an IEEE-754 binary32 floating-point number.
	//
	// NaN and infinities SHOULD NOT be considered valid portable API values
	// unless a future descriptor rule explicitly permits them.
	DescriptorFloat32

	// DescriptorFloat64 represents an IEEE-754 binary64 floating-point number.
	//
	// NaN and infinities SHOULD NOT be considered valid portable API values
	// unless a future descriptor rule explicitly permits them.
	DescriptorFloat64

	// DescriptorDecimal represents an exact base-10 numeric value.
	//
	// DescriptorDecimal is intended for values where binary floating-point semantics
	// are not acceptable, such as money, exact quotas, ratios, limits, or
	// human-authored decimal configuration. Precision, scale, rounding policy,
	// and encoding are descriptor rules or codec concerns, not DescriptorKind concerns.
	DescriptorDecimal

	// DescriptorTimestamp represents an absolute point in time.
	//
	// DescriptorTimestamp is a semantic API descriptor. It does not prescribe whether a
	// codec stores the value as RFC 3339 text, Unix nanoseconds, Unix
	// milliseconds, or another canonical representation.
	DescriptorTimestamp

	// DescriptorDate represents a calendar date without a time of day.
	//
	// DescriptorDate is distinct from DescriptorTimestamp. It SHOULD be used for values such
	// as billing dates, schedule dates, report dates, or other date-only
	// concepts.
	DescriptorDate

	// DescriptorTime represents a time of day without a calendar date.
	//
	// DescriptorTime is distinct from DescriptorTimestamp and DescriptorDuration.
	// It SHOULD be used for values such as daily schedule boundaries, local opening
	// times, or other time-of-day concepts.
	DescriptorTime

	// DescriptorDuration represents an elapsed interval.
	//
	// DescriptorDuration is a semantic API descriptor. It does not prescribe
	// whether a codec stores the value as nanoseconds, ISO-8601 duration text,
	// Go-style duration text, or another canonical representation.
	DescriptorDuration

	// DescriptorObject represents a structural object with a fixed set of named fields.
	//
	// DescriptorObject is for schema-defined records. Dynamic key/value dictionaries
	// MUST use DescriptorMap. Unknown-field behavior is controlled by object descriptor
	// rules, not by DescriptorKind.
	DescriptorObject

	// DescriptorList represents an ordered sequence of values of one element descriptor.
	//
	// Merge/apply behavior, such as atomic list, set-like list, or map-like
	// list, is controlled by list descriptor rules, not by DescriptorKind.
	DescriptorList

	// DescriptorMap represents a dynamic key/value mapping.
	//
	// DescriptorMap is distinct from DescriptorObject. DescriptorObject has fixed
	// schema fields; DescriptorMap has dynamic keys and a common value descriptor.
	// Key restrictions are controlled by map descriptor rules.
	DescriptorMap

	// DescriptorRef represents a reference to a named custom descriptor definition.
	//
	// DescriptorRef is the extension mechanism for reusable semantic descriptors.
	// A DescriptorRef MUST resolve through an owner-created resolver. It MUST NOT
	// refer to arbitrary Go implementations.
	DescriptorRef
)

// IsValid reports whether c identifies a supported structural descriptor category.
func (c DescriptorKind) IsValid() bool {
	return c > DescriptorInvalid && c <= DescriptorRef
}

// String returns the stable diagnostic name of c.
//
// String is intended for diagnostics, tests, and developer-facing output. It is
// not a wire-format contract.
func (c DescriptorKind) String() string {
	switch c {
	case DescriptorInvalid:
		return "invalid"
	case DescriptorNull:
		return "null"
	case DescriptorBool:
		return "bool"
	case DescriptorString:
		return "string"
	case DescriptorBytes:
		return "bytes"
	case DescriptorInt8:
		return "int8"
	case DescriptorInt16:
		return "int16"
	case DescriptorInt32:
		return "int32"
	case DescriptorInt64:
		return "int64"
	case DescriptorUint8:
		return "uint8"
	case DescriptorUint16:
		return "uint16"
	case DescriptorUint32:
		return "uint32"
	case DescriptorUint64:
		return "uint64"
	case DescriptorFloat32:
		return "float32"
	case DescriptorFloat64:
		return "float64"
	case DescriptorDecimal:
		return "decimal"
	case DescriptorTimestamp:
		return "timestamp"
	case DescriptorDate:
		return "date"
	case DescriptorTime:
		return "time"
	case DescriptorDuration:
		return "duration"
	case DescriptorObject:
		return "object"
	case DescriptorList:
		return "list"
	case DescriptorMap:
		return "map"
	case DescriptorRef:
		return "ref"
	default:
		return "unknown"
	}
}

// IsPrimitive reports whether c is a leaf value descriptor.
//
// Primitive includes null, bool, string, bytes, numbers, and temporal values.
// It excludes composites and references. DescriptorRef is a named indirection rather
// than a primitive payload.
func (c DescriptorKind) IsPrimitive() bool {
	return c == DescriptorNull ||
		c == DescriptorBool ||
		c == DescriptorString ||
		c == DescriptorBytes ||
		c.IsNumber() ||
		c.IsTemporal()
}

// IsNumber reports whether c belongs to the numeric descriptor category.
//
// Decimal is included because it is a number even though it is not an integer or
// binary floating-point type. Decimal exact values intentionally do not have
// min/max support in this package until a decimal value representation exists.
func (c DescriptorKind) IsNumber() bool {
	return c.IsInteger() || c.IsFloat() || c == DescriptorDecimal
}

// IsInteger reports whether c is a fixed-width int or uint descriptor.
func (c DescriptorKind) IsInteger() bool {
	return c.IsInt() || c.IsUint()
}

// IsInt reports whether c is a fixed-width int descriptor.
func (c DescriptorKind) IsInt() bool {
	return c >= DescriptorInt8 && c <= DescriptorInt64
}

// IsUint reports whether c is a fixed-width uint descriptor.
func (c DescriptorKind) IsUint() bool {
	return c >= DescriptorUint8 && c <= DescriptorUint64
}

// IsFloat reports whether c is a binary floating-point type.
func (c DescriptorKind) IsFloat() bool {
	return c == DescriptorFloat32 || c == DescriptorFloat64
}

// IsTemporal reports whether c belongs to the semantic temporal kind.
//
// Temporal DescriptorKind values do not prescribe a codec representation. RFC 3339,
// ISO-8601 durations, Unix timestamps, and Go-style duration text are encoding
// choices for future codec packages.
func (c DescriptorKind) IsTemporal() bool {
	return c >= DescriptorTimestamp && c <= DescriptorDuration
}

// IsComposite reports whether c directly contains nested structural descriptors.
//
// DescriptorRef is intentionally excluded. It points at a Resolver definition but is
// not itself a composite payload.
func (c DescriptorKind) IsComposite() bool {
	return c == DescriptorObject || c == DescriptorList || c == DescriptorMap
}

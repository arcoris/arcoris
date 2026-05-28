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

// validateInactivePayloads rejects impossible package-internal descriptor state.
//
// Public builders normalize exactly one payload slot into Type. This check
// guards future internal changes and package-local tests from accidentally
// creating descriptors where the TypeCode and populated payload slot disagree.
func validateInactivePayloads(t Type, path string) error {
	if t.code != TypeString && !emptyStringPayload(t.string) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeBytes && !emptyBytesPayload(t.bytes) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeInt8 && !emptyInt8Payload(t.int8) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeInt16 && !emptyInt16Payload(t.int16) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeInt32 && !emptyInt32Payload(t.int32) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeInt64 && !emptyInt64Payload(t.int64) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeUint8 && !emptyUint8Payload(t.uint8) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeUint16 && !emptyUint16Payload(t.uint16) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeUint32 && !emptyUint32Payload(t.uint32) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeUint64 && !emptyUint64Payload(t.uint64) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeFloat32 && !emptyFloat32Payload(t.float32) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeFloat64 && !emptyFloat64Payload(t.float64) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeDecimal && !emptyDecimalPayload(t.decimal) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeObject && !emptyObjectPayload(t.object) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeList && !emptyListPayload(t.list) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeMap && !emptyMapPayload(t.mapType) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeRef && !emptyRefPayload(t.ref) {
		return typeError(path+".payload", ErrInvalidType)
	}
	return nil
}

// emptyStringPayload reports whether p has no configured TypeString state.
func emptyStringPayload(p stringPayload) bool {
	return !p.minLen.set && !p.maxLen.set && !p.hasPattern && p.pattern == "" && len(p.enum) == 0
}

// emptyBytesPayload reports whether p has no configured TypeBytes state.
func emptyBytesPayload(p bytesPayload) bool {
	return !p.minLen.set && !p.maxLen.set
}

// emptyInt8Payload reports whether p has no configured TypeInt8 state.
func emptyInt8Payload(p int8Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}

// emptyInt16Payload reports whether p has no configured TypeInt16 state.
func emptyInt16Payload(p int16Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}

// emptyInt32Payload reports whether p has no configured TypeInt32 state.
func emptyInt32Payload(p int32Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}

// emptyInt64Payload reports whether p has no configured TypeInt64 state.
func emptyInt64Payload(p int64Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}

// emptyUint8Payload reports whether p has no configured TypeUint8 state.
func emptyUint8Payload(p uint8Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}

// emptyUint16Payload reports whether p has no configured TypeUint16 state.
func emptyUint16Payload(p uint16Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}

// emptyUint32Payload reports whether p has no configured TypeUint32 state.
func emptyUint32Payload(p uint32Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}

// emptyUint64Payload reports whether p has no configured TypeUint64 state.
func emptyUint64Payload(p uint64Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}

// emptyFloat32Payload reports whether p has no configured TypeFloat32 state.
func emptyFloat32Payload(p float32Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}

// emptyFloat64Payload reports whether p has no configured TypeFloat64 state.
func emptyFloat64Payload(p float64Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}

// emptyDecimalPayload reports whether p has no configured TypeDecimal state.
func emptyDecimalPayload(p decimalPayload) bool {
	return !p.precision.set && !p.scale.set
}

// emptyObjectPayload reports whether p has no configured TypeObject state.
func emptyObjectPayload(p objectPayload) bool {
	return len(p.fields) == 0 && p.unknown == UnknownReject
}

// emptyListPayload reports whether p has no configured TypeList state.
func emptyListPayload(p listPayload) bool {
	return p.elem == nil && !p.minLen.set && !p.maxLen.set && p.semantics == ListAtomic && len(p.mapKeys) == 0
}

// emptyMapPayload reports whether p has no configured TypeMap state.
func emptyMapPayload(p mapPayload) bool {
	return p.key == MapKeyString && p.value == nil && !p.minLen.set && !p.maxLen.set
}

// emptyRefPayload reports whether p has no configured TypeRef state.
func emptyRefPayload(p refPayload) bool {
	return p.name == ""
}

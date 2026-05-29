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
		return inactivePayloadError(path, t.code, "string")
	}
	if t.code != TypeBytes && !emptyBytesPayload(t.bytes) {
		return inactivePayloadError(path, t.code, "bytes")
	}
	if t.code != TypeInt8 && !emptyInt8Payload(t.int8) {
		return inactivePayloadError(path, t.code, "int8")
	}
	if t.code != TypeInt16 && !emptyInt16Payload(t.int16) {
		return inactivePayloadError(path, t.code, "int16")
	}
	if t.code != TypeInt32 && !emptyInt32Payload(t.int32) {
		return inactivePayloadError(path, t.code, "int32")
	}
	if t.code != TypeInt64 && !emptyInt64Payload(t.int64) {
		return inactivePayloadError(path, t.code, "int64")
	}
	if t.code != TypeUint8 && !emptyUint8Payload(t.uint8) {
		return inactivePayloadError(path, t.code, "uint8")
	}
	if t.code != TypeUint16 && !emptyUint16Payload(t.uint16) {
		return inactivePayloadError(path, t.code, "uint16")
	}
	if t.code != TypeUint32 && !emptyUint32Payload(t.uint32) {
		return inactivePayloadError(path, t.code, "uint32")
	}
	if t.code != TypeUint64 && !emptyUint64Payload(t.uint64) {
		return inactivePayloadError(path, t.code, "uint64")
	}
	if t.code != TypeFloat32 && !emptyFloat32Payload(t.float32) {
		return inactivePayloadError(path, t.code, "float32")
	}
	if t.code != TypeFloat64 && !emptyFloat64Payload(t.float64) {
		return inactivePayloadError(path, t.code, "float64")
	}
	if t.code != TypeDecimal && !emptyDecimalPayload(t.decimal) {
		return inactivePayloadError(path, t.code, "decimal")
	}
	// Temporal payload slots are currently empty, but they stay in this matrix
	// so future temporal constraints inherit the same exact-slot protection.
	if t.code != TypeTimestamp && !emptyTimestampPayload(t.timestamp) {
		return inactivePayloadError(path, t.code, "timestamp")
	}
	if t.code != TypeDate && !emptyDatePayload(t.date) {
		return inactivePayloadError(path, t.code, "date")
	}
	if t.code != TypeTime && !emptyTimePayload(t.timeOfDay) {
		return inactivePayloadError(path, t.code, "timeOfDay")
	}
	if t.code != TypeDuration && !emptyDurationPayload(t.duration) {
		return inactivePayloadError(path, t.code, "duration")
	}
	if t.code != TypeObject && !emptyObjectPayload(t.object) {
		return inactivePayloadError(path, t.code, "object")
	}
	if t.code != TypeList && !emptyListPayload(t.list) {
		return inactivePayloadError(path, t.code, "list")
	}
	if t.code != TypeMap && !emptyMapPayload(t.mapType) {
		return inactivePayloadError(path, t.code, "mapType")
	}
	if t.code != TypeRef && !emptyRefPayload(t.ref) {
		return inactivePayloadError(path, t.code, "ref")
	}
	return nil
}

// inactivePayloadError reports a populated payload slot outside its TypeCode.
func inactivePayloadError(path string, code TypeCode, slot string) error {
	return typeErrorf(
		path+".payload",
		ErrInvalidType,
		TypeErrorReasonInactivePayload,
		"descriptor uses %s but %s payload is populated",
		code,
		slot,
	)
}

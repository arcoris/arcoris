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

// validateInactiveDescriptorPayloads rejects impossible package-internal descriptor state.
//
// Public builders normalize exactly one payload slot into Descriptor. This check
// guards future internal changes and package-local tests from accidentally
// creating descriptors where the DescriptorKind and populated payload slot disagree.
func validateInactiveDescriptorPayloads(desc Descriptor, path string) error {
	if desc.code != DescriptorString && !emptyStringPayload(desc.string) {
		return inactivePayloadError(path, desc.code, "string")
	}

	if desc.code != DescriptorBytes && !emptyBytesPayload(desc.bytes) {
		return inactivePayloadError(path, desc.code, "bytes")
	}

	if desc.code != DescriptorInt8 && !emptyInt8Payload(desc.int8) {
		return inactivePayloadError(path, desc.code, "int8")
	}

	if desc.code != DescriptorInt16 && !emptyInt16Payload(desc.int16) {
		return inactivePayloadError(path, desc.code, "int16")
	}

	if desc.code != DescriptorInt32 && !emptyInt32Payload(desc.int32) {
		return inactivePayloadError(path, desc.code, "int32")
	}

	if desc.code != DescriptorInt64 && !emptyInt64Payload(desc.int64) {
		return inactivePayloadError(path, desc.code, "int64")
	}

	if desc.code != DescriptorUint8 && !emptyUint8Payload(desc.uint8) {
		return inactivePayloadError(path, desc.code, "uint8")
	}

	if desc.code != DescriptorUint16 && !emptyUint16Payload(desc.uint16) {
		return inactivePayloadError(path, desc.code, "uint16")
	}

	if desc.code != DescriptorUint32 && !emptyUint32Payload(desc.uint32) {
		return inactivePayloadError(path, desc.code, "uint32")
	}

	if desc.code != DescriptorUint64 && !emptyUint64Payload(desc.uint64) {
		return inactivePayloadError(path, desc.code, "uint64")
	}

	if desc.code != DescriptorFloat32 && !emptyFloat32Payload(desc.float32) {
		return inactivePayloadError(path, desc.code, "float32")
	}

	if desc.code != DescriptorFloat64 && !emptyFloat64Payload(desc.float64) {
		return inactivePayloadError(path, desc.code, "float64")
	}

	if desc.code != DescriptorDecimal && !emptyDecimalPayload(desc.decimal) {
		return inactivePayloadError(path, desc.code, "decimal")
	}
	// Temporal payload slots are currently empty, but they stay in this matrix
	// so future temporal constraints inherit the same exact-slot protection.
	if desc.code != DescriptorTimestamp && !emptyTimestampPayload(desc.timestamp) {
		return inactivePayloadError(path, desc.code, "timestamp")
	}

	if desc.code != DescriptorDate && !emptyDatePayload(desc.date) {
		return inactivePayloadError(path, desc.code, "date")
	}

	if desc.code != DescriptorTime && !emptyTimePayload(desc.timeOfDay) {
		return inactivePayloadError(path, desc.code, "timeOfDay")
	}

	if desc.code != DescriptorDuration && !emptyDurationPayload(desc.duration) {
		return inactivePayloadError(path, desc.code, "duration")
	}

	if desc.code != DescriptorObject && !emptyObjectPayload(desc.object) {
		return inactivePayloadError(path, desc.code, "object")
	}

	if desc.code != DescriptorList && !emptyListPayload(desc.list) {
		return inactivePayloadError(path, desc.code, "list")
	}

	if desc.code != DescriptorMap && !emptyMapPayload(desc.mapType) {
		return inactivePayloadError(path, desc.code, "mapType")
	}

	if desc.code != DescriptorRef && !emptyRefPayload(desc.ref) {
		return inactivePayloadError(path, desc.code, "ref")
	}

	return nil
}

// inactivePayloadError reports a populated payload slot outside its DescriptorKind.
func inactivePayloadError(path string, code DescriptorKind, slot string) error {
	return descriptorErrorf(
		path+".payload",
		ErrInvalidDescriptor,
		DescriptorErrorReasonInactivePayload,
		"descriptor uses %s but %s payload is populated",
		code,
		slot,
	)
}

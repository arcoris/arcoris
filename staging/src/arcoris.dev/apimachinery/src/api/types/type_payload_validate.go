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
	if t.code != TypeTimestamp && !emptyTimestampPayload(t.timestamp) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeDate && !emptyDatePayload(t.date) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeTime && !emptyTimePayload(t.timeOfDay) {
		return typeError(path+".payload", ErrInvalidType)
	}
	if t.code != TypeDuration && !emptyDurationPayload(t.duration) {
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

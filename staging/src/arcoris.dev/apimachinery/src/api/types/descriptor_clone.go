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

// cloneDescriptor detaches all slice-bearing payloads from caller-owned state.
//
// Descriptor is a value object, but several exact payload slots contain slices to
// preserve declaration order. Every public boundary that returns or stores a
// Descriptor calls cloneDescriptor so later caller mutation cannot rewrite
// descriptors held by catalogs, fields, or views.
func cloneDescriptor(desc Descriptor) Descriptor {
	desc.string = cloneStringPayload(desc.string)
	desc.bytes = cloneBytesPayload(desc.bytes)
	desc.int8 = cloneInt8Payload(desc.int8)
	desc.int16 = cloneInt16Payload(desc.int16)
	desc.int32 = cloneInt32Payload(desc.int32)
	desc.int64 = cloneInt64Payload(desc.int64)
	desc.uint8 = cloneUint8Payload(desc.uint8)
	desc.uint16 = cloneUint16Payload(desc.uint16)
	desc.uint32 = cloneUint32Payload(desc.uint32)
	desc.uint64 = cloneUint64Payload(desc.uint64)
	desc.float32 = cloneFloat32Payload(desc.float32)
	desc.float64 = cloneFloat64Payload(desc.float64)
	desc.decimal = cloneDecimalPayload(desc.decimal)
	desc.timestamp = cloneTimestampPayload(desc.timestamp)
	desc.date = cloneDatePayload(desc.date)
	desc.timeOfDay = cloneTimePayload(desc.timeOfDay)
	desc.duration = cloneDurationPayload(desc.duration)
	desc.object = cloneObjectPayload(desc.object)
	desc.list = cloneListPayload(desc.list)
	desc.mapType = cloneMapPayload(desc.mapType)
	desc.ref = cloneRefPayload(desc.ref)

	return desc
}

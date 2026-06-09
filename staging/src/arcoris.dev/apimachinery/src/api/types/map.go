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

// MapDescriptor builds descriptors for dynamic string-keyed maps.
//
// MapDescriptor is for dictionaries with dynamic keys and one shared value descriptor. It
// is intentionally separate from ObjectDescriptor, which models fixed schema fields.
// Concrete map keys are strings; Keys can constrain those string tokens with a
// string-like descriptor or a reference resolving to one.
type MapDescriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores the exact map shape under construction.
	payload mapPayload
}

// MapOf returns a string-keyed map descriptor builder for value.
//
// A nil DescriptorExpr is recorded as an invalid zero value descriptor so
// ValidateResolved can classify the error at map.value. The builder itself stays
// allocation-light and panic-free.
//
// Typical reusable declaration:
//
//	labelValue := String().MinBytes(1)
//
//	labelsType := MapOf(
//		labelValue,
//	).Keys(String().MinBytes(1)).MaxEntries(64)
func MapOf(value DescriptorExpr) MapDescriptor {
	keyType := String().Descriptor()
	valueType := descriptorFromExpr(value)

	return MapDescriptor{
		header: newHeader(DescriptorMap),
		payload: mapPayload{
			key:   &keyType,
			value: &valueType,
		},
	}
}

// Keys constrains concrete string keys for the map.
//
// key must be a non-nullable string descriptor or a reference that resolves to
// a non-nullable string descriptor during resolved validation. The descriptor
// validates key tokens only; it does not turn map keys into objects, routes, or
// storage paths.
func (desc MapDescriptor) Keys(key DescriptorExpr) MapDescriptor {
	keyType := descriptorFromExpr(key)
	desc.payload.key = &keyType

	return desc
}

// Nullable returns a map descriptor that admits null values.
func (desc MapDescriptor) Nullable() MapDescriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// MinEntries sets the inclusive minimum number of map entries.
//
// The limit is structural metadata only. Concrete map entry counts are checked
// by future value-validation layers.
func (desc MapDescriptor) MinEntries(n int) MapDescriptor {
	desc.payload.minLen = limit[int]{n, true}

	return desc
}

// MaxEntries sets the inclusive maximum number of map entries.
//
// The limit uses limit[int] so an explicit zero maximum can be represented
// without a pointer allocation.
func (desc MapDescriptor) MaxEntries(n int) MapDescriptor {
	desc.payload.maxLen = limit[int]{n, true}

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc MapDescriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.mapType = cloneMapPayload(desc.payload)

	return out
}

// descriptorExpr marks MapDescriptor as a sealed DescriptorExpr implementation.
func (desc MapDescriptor) descriptorExpr() {}

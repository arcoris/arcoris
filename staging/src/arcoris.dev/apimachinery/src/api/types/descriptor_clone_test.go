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

import "testing"

func TestCloneDescriptorDetachesExactPayloadSlots(t *testing.T) {
	desc := Object(
		Field("name").String().Required().Enum("alpha", "beta"),
		Field("items").ListOf(
			Object(Field("key").String().Required()),
		).Required().Map("key"),
		Field("lookup").MapOf(String().Enum("one", "two")).Optional(),
	).Descriptor()

	cloned := cloneDescriptor(desc)
	cloned.object.fields[0] = Field("changed").String().Required().Field()
	cloned.object.fields[1].descriptor.list.mapKeys[0] = "changed"
	cloned.object.fields[2].descriptor.mapType.value.string.enum[0] = "changed"

	fields := desc.object.fields
	requireEqual(t, fields[0].Name(), FieldName("name"))
	requireEqual(t, fields[1].descriptor.list.mapKeys[0], FieldName("key"))
	requireEqual(t, fields[2].descriptor.mapType.value.string.enum[0], "one")
}

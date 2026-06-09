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

func TestFieldDescriptorCloneHelpersDetachDescriptorAndSlices(t *testing.T) {
	field := Field("name").String().Required().Enum("a").Field()
	cloned := cloneField(field)
	cloned.descriptor.string.enum[0] = "b"

	view := requireStringView(t, field.Descriptor())
	requireEqual(t, view.Enum()[0], "a")

	fields := []FieldDescriptor{field}
	clonedFields := cloneFields(fields)
	clonedFields[0] = Field("changed").String().Required().Field()
	requireEqual(t, fields[0].Name(), FieldName("name"))
	requireEqual(t, cloneFields(nil) == nil, true)
}

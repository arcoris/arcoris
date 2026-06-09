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

func TestDefinitionAccessorsAndDetach(t *testing.T) {
	def := Define("example.Name", String().Enum("alpha")).WithDescription("Name")
	desc := def.Descriptor()
	desc.string.enum[0] = "changed"

	requireEqual(t, def.Name(), TypeName("example.Name"))
	requireEqual(t, def.Description(), "Name")
	view, _ := def.Descriptor().AsString()
	requireEqual(t, view.Enum()[0], "alpha")
}

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

package fieldpath

import "testing"

func TestPathElementsReturnsClone(t *testing.T) {
	path := RootPath().Field("spec")
	elements := path.Elements()

	elements[0] = FieldElement("status")

	requireEqual(t, path.String(), "$.spec")
}

func TestPathAppendDoesNotMutateReceiver(t *testing.T) {
	root := RootPath().Field("spec")
	child := root.Append(FieldElement("replicas"))

	requireEqual(t, root.String(), "$.spec")
	requireEqual(t, child.String(), "$.spec.replicas")
}

func TestPathAppendClonesSelectorElement(t *testing.T) {
	selector := MustSelector(NewSelectorEntry("type", StringLiteral("Ready")))
	element := SelectorElement(selector)

	path := RootPath().Append(element)
	element.selector.entries[0] = NewSelectorEntry("type", StringLiteral("Changed"))

	requireEqual(t, path.String(), `$[{"type":"Ready"}]`)
}

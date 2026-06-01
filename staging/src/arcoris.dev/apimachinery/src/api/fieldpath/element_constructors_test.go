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

func TestFieldElement(t *testing.T) {
	element := FieldElement("spec")

	requireEqual(t, element.Kind(), ElementField)
	requireEqual(t, element.Name(), "spec")
}

func TestKeyElement(t *testing.T) {
	element := KeyElement("app")

	requireEqual(t, element.Kind(), ElementKey)
	requireEqual(t, element.Name(), "app")
}

func TestIndexElement(t *testing.T) {
	element := IndexElement(3)

	requireEqual(t, element.Kind(), ElementIndex)
	requireEqual(t, element.Index(), 3)
}

func TestSelectorElement(t *testing.T) {
	selector := MustSelector(NewSelectorEntry("type", StringLiteral("Ready")))
	element := SelectorElement(selector)

	requireEqual(t, element.Kind(), ElementSelector)
	requireEqual(t, element.Selector().Equal(selector), true)
}

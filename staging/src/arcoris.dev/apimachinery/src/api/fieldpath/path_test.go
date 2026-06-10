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

func TestPathRootState(t *testing.T) {
	path := Root()

	requireEqual(t, path.Len(), 0)
	requireEqual(t, path.IsRoot(), true)
}

func TestPathElementsReturnsClone(t *testing.T) {
	path := Root().Field(testField("spec"))
	elements := path.Elements()

	elements[0] = testFieldElement("status")

	requireEqual(t, path.String(), "$.spec")
}

func TestPathElementReturnsElementByIndex(t *testing.T) {
	path := Root().Field(testField("spec"))

	element, ok := path.Element(0)
	requireEqual(t, ok, true)
	requireEqual(t, element.String(), ".spec")

	_, ok = path.Element(1)
	requireEqual(t, ok, false)
}

func TestPathForEachStopsEarly(t *testing.T) {
	path := Root().Field(testField("spec")).Field(testField("replicas"))
	visited := 0

	path.ForEach(func(index int, element Element) bool {
		visited++
		return false
	})

	requireEqual(t, visited, 1)
}

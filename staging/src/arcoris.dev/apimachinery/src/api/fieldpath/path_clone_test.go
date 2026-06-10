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

func TestCloneElementsDetachesSelectorElements(t *testing.T) {
	selector := MustSelector(testSelectorEntry("type", StringLiteral("Ready")))
	elements := []Element{testSelectorElement(selector)}

	cloned := cloneElements(elements)
	elements[0].selector.entries[0] = testSelectorEntry("type", StringLiteral("Changed"))

	requireEqual(t, cloned[0].String(), `[{"type":"Ready"}]`)
}

func TestAppendClonedElementsPreservesExistingDestination(t *testing.T) {
	dst := []Element{testFieldElement("spec")}
	got := appendClonedElements(dst, []Element{testFieldElement("replicas")})

	requireEqual(t, len(got), 2)
	requireEqual(t, got[0].String(), ".spec")
	requireEqual(t, got[1].String(), ".replicas")
}

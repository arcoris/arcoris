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

func TestElementAsAccessorsRejectWrongKind(t *testing.T) {
	element := NewFieldElement(testField("spec"))

	_, keyOK := element.AsKey()
	_, indexOK := element.AsIndex()
	_, selectorOK := element.AsSelector()

	requireEqual(t, element.Kind(), ElementField)
	requireEqual(t, keyOK, false)
	requireEqual(t, indexOK, false)
	requireEqual(t, selectorOK, false)
}

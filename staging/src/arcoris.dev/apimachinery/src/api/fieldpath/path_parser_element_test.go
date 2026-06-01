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

func TestPathParserParseFieldElementQuoted(t *testing.T) {
	p := newPathParser(`."api-version"`)

	element, err := p.parseFieldElement()
	requireNoError(t, err)

	requireEqual(t, element.Kind(), ElementField)
	requireEqual(t, element.Name(), "api-version")
	requireEqual(t, p.done(), true)
}

func TestPathParserParseBracketElementKey(t *testing.T) {
	p := newPathParser(`["app.kubernetes.io/name"]`)

	element, err := p.parseBracketElement()
	requireNoError(t, err)

	requireEqual(t, element.Kind(), ElementKey)
	requireEqual(t, element.Name(), "app.kubernetes.io/name")
	requireEqual(t, p.done(), true)
}

func TestPathParserParseBracketElementIndex(t *testing.T) {
	p := newPathParser(`[17]`)

	element, err := p.parseBracketElement()
	requireNoError(t, err)

	requireEqual(t, element.Kind(), ElementIndex)
	requireEqual(t, element.Index(), 17)
	requireEqual(t, p.done(), true)
}

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

func TestPathParserParseBracketElementSelector(t *testing.T) {
	p := newPathParser(`[{"port":443,"host":"api.example.com"}]`)

	element, err := p.parseBracketElement()
	requireNoError(t, err)

	requireEqual(t, element.Kind(), ElementSelector)
	requireEqual(
		t,
		element.Selector().String(),
		`{"host":"api.example.com","port":443}`,
	)
	requireEqual(t, p.done(), true)
}

func TestPathParserParseSelectorEntryRequiresColon(t *testing.T) {
	p := newPathParser(`"type" "Ready"`)

	_, err := p.parseSelectorEntry()

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorIs(t, err, ErrInvalidSyntax)
}

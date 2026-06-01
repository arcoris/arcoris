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

func TestSelectorEntryString(t *testing.T) {
	requireEqual(
		t,
		NewSelectorEntry("type", StringLiteral("Ready")).String(),
		`"type":"Ready"`,
	)
}

func TestSelectorString(t *testing.T) {
	selector := MustSelector(
		NewSelectorEntry("port", Uint64Literal(443)),
		NewSelectorEntry("host", StringLiteral("api.example.com")),
	)

	requireEqual(t, selector.String(), `{"host":"api.example.com","port":443}`)
}

func TestSelectorStringEscapesValues(t *testing.T) {
	selector := MustSelector(NewSelectorEntry(`na"me`, StringLiteral(`a"b`)))

	requireEqual(t, selector.String(), `{"na\"me":"a\"b"}`)
}

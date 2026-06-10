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

func TestSelectorEqualIgnoresInputOrder(t *testing.T) {
	left := MustSelector(
		testSelectorEntry("port", Uint64Literal(443)),
		testSelectorEntry("host", StringLiteral("api.example.com")),
	)

	right := MustSelector(
		testSelectorEntry("host", StringLiteral("api.example.com")),
		testSelectorEntry("port", Uint64Literal(443)),
	)

	requireEqual(t, left.Equal(right), true)
}

func TestSelectorCompare(t *testing.T) {
	left := MustSelector(testSelectorEntry("a", StringLiteral("x")))
	right := MustSelector(testSelectorEntry("b", StringLiteral("x")))

	requireEqual(t, left.Compare(right), -1)
}

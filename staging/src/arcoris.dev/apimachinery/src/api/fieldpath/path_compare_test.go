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

func TestPathEqual(t *testing.T) {
	left := Root().Field(testField("spec")).Field(testField("replicas"))
	right := Root().Field(testField("spec")).Field(testField("replicas"))

	requireEqual(t, left.Equal(right), true)
}

func TestPathCompareStableOrder(t *testing.T) {
	testCases := []struct {
		name  string
		left  Path
		right Path
		want  int
	}{
		{
			name:  "shorter before child",
			left:  Root().Field(testField("spec")),
			right: Root().Field(testField("spec")).Field(testField("replicas")),
			want:  -1,
		},
		{
			name:  "field before key",
			left:  MustPath(testFieldElement("a")),
			right: MustPath(testKeyElement("a")),
			want:  -1,
		},
		{
			name:  "key before index",
			left:  MustPath(testKeyElement("a")),
			right: MustPath(MustIndexElement(0)),
			want:  -1,
		},
		{
			name:  "index before selector",
			left:  MustPath(MustIndexElement(0)),
			right: MustPath(testSelectorElement(MustSelector(testSelectorEntry("a", StringLiteral("x"))))),
			want:  -1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			requireEqual(t, testCase.left.Compare(testCase.right), testCase.want)
		})
	}
}

func TestPathCompareOrdersSelectorsDeterministically(t *testing.T) {
	left := Root().
		Field("conditions").
		Select(MustSelector(testSelectorEntry("type", StringLiteral("Ready"))))

	right := Root().
		Field("conditions").
		Select(MustSelector(testSelectorEntry("type", StringLiteral("Scheduled"))))

	requireEqual(t, left.Compare(right), -1)
}

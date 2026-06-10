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

func TestEmptySet(t *testing.T) {
	set := EmptySet()

	requireEqual(t, set.Len(), 0)
	requireEqual(t, set.IsEmpty(), true)
	requireEqual(t, len(set.Paths()), 0)
	requireEqual(t, set.String(), "{}")
}

func TestSetPathsReturnsDetachedSlice(t *testing.T) {
	set := MustSet(setReplicasPath(), setImagePath())
	paths := set.Paths()

	paths[0] = Root().Field(testField("status"))

	requireEqual(t, set.Has(Root().Field(testField("status"))), false)
	requireEqual(t, set.Len(), 2)
}

func TestSetPathsReturnsDetachedPathElements(t *testing.T) {
	set := MustSet(setReplicasPath())
	paths := set.Paths()

	paths[0].elements[0] = testFieldElement("status")

	requireEqual(t, set.Has(setReplicasPath()), true)
	requireEqual(t, set.Has(Root().Field(testField("status")).Field(testField("replicas"))), false)
}

func TestSetForEachStopsEarly(t *testing.T) {
	set := MustSet(setReplicasPath(), setImagePath())
	visited := 0

	set.ForEach(func(index int, path Path) bool {
		visited++
		return false
	})

	requireEqual(t, visited, 1)
}

func TestSetHas(t *testing.T) {
	set := MustSet(setReplicasPath(), setImagePath(), setLabelPath())

	requireEqual(t, set.Has(setReplicasPath()), true)
	requireEqual(t, set.Has(Root().Field(testField("status"))), false)
}

func TestSetHasUsesStructuralPaths(t *testing.T) {
	field := Root().Field(testField("api-version"))
	key := Root().Key(testKey("api-version"))
	set := MustSet(field)

	requireEqual(t, set.Has(field), true)
	requireEqual(t, set.Has(key), false)
}

func TestSetHasOnEmptySet(t *testing.T) {
	requireEqual(t, EmptySet().Has(Root()), false)
}

func setSpecPath() Path {
	return Root().Field(testField("spec"))
}

func setReplicasPath() Path {
	return Root().Field(testField("spec")).Field(testField("replicas"))
}

func setImagePath() Path {
	return Root().Field(testField("spec")).Field(testField("image"))
}

func setLabelPath() Path {
	return Root().Field(testField("metadata")).Field(testField("labels")).Key(testKey("app"))
}

func setIndexPath() Path {
	return Root().Field(testField("items")).Index(0)
}

func setReadyStatusPath() Path {
	ready := MustSelector(
		testSelectorEntry("type", StringLiteral("Ready")),
	)

	return Root().
		Field("conditions").
		Select(ready).
		Field("status")
}

func setPathStrings(paths []Path) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		out[i] = p.String()
	}

	return out
}

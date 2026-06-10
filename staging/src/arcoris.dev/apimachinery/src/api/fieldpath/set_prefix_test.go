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

func TestSetHasAnyUnder(t *testing.T) {
	set := MustSet(
		setSpecPath(),
		setReplicasPath(),
		Root().Field(testField("spec")).Field(testField("template")).Field(testField("metadata")).Field(testField("labels")).Key(testKey("app")),
		Root().Field(testField("metadata")).Field(testField("name")),
	)

	requireEqual(t, set.Has(setSpecPath()), true)
	requireEqual(t, set.HasAnyUnder(setSpecPath()), true)
	requireEqual(t, set.HasAnyUnder(setReplicasPath()), true)
	requireEqual(t, set.HasAnyUnder(Root().Field(testField("status"))), false)
}

func TestSetHasAnyUnderMatchesDescendantWithoutExactPrefix(t *testing.T) {
	set := MustSet(setReplicasPath())

	requireEqual(t, set.Has(setSpecPath()), false)
	requireEqual(t, set.HasAnyUnder(setSpecPath()), true)
}

func TestSetHasDescendant(t *testing.T) {
	exactOnly := MustSet(setSpecPath())
	withDescendant := MustSet(setSpecPath(), setReplicasPath())

	requireEqual(t, exactOnly.HasDescendant(setSpecPath()), false)
	requireEqual(t, withDescendant.HasDescendant(setSpecPath()), true)
}

func TestSetUnder(t *testing.T) {
	templateLabel := Root().
		Field("spec").
		Field("template").
		Field("metadata").
		Field("labels").
		Key("app")
	set := MustSet(
		setSpecPath(),
		setReplicasPath(),
		templateLabel,
		Root().Field(testField("metadata")).Field(testField("name")),
	)

	underSpec := set.Under(setSpecPath())

	requireStringSliceEqual(t, setPathStrings(underSpec.Paths()), []string{
		"$.spec",
		"$.spec.replicas",
		`$.spec.template.metadata.labels["app"]`,
	})
}

func TestSetRemoveDescendantsPreservesPrefix(t *testing.T) {
	set := MustSet(setSpecPath(), setReplicasPath(), setImagePath())

	got := set.RemoveDescendants(setSpecPath())

	requireStringSliceEqual(t, setPathStrings(got.Paths()), []string{"$.spec"})
}

func TestSetCompactSubtreesRemovesDescendants(t *testing.T) {
	set := MustSet(setSpecPath(), setReplicasPath(), setImagePath())

	got := set.CompactSubtrees()

	requireStringSliceEqual(t, setPathStrings(got.Paths()), []string{"$.spec"})
}

func TestSetUnderReturnsEmptySetWhenPrefixAbsent(t *testing.T) {
	set := MustSet(setReplicasPath(), setImagePath())

	underStatus := set.Under(Root().Field(testField("status")))

	requireEqual(t, underStatus.IsEmpty(), true)
	requireEqual(t, underStatus.String(), "{}")
}

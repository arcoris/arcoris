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

func TestNewSetSortsAndDeduplicates(t *testing.T) {
	set, err := NewSet(
		setReplicasPath(),
		setLabelPath(),
		setReadyStatusPath(),
		setImagePath(),
		setIndexPath(),
		setReplicasPath(),
	)
	requireNoError(t, err)

	requireStringSliceEqual(t, setPathStrings(set.Paths()), []string{
		`$.conditions[{"type":"Ready"}].status`,
		"$.items[0]",
		`$.metadata.labels["app"]`,
		"$.spec.image",
		"$.spec.replicas",
	})
}

func TestNewSetKeepsAncestorAndDescendant(t *testing.T) {
	set := MustSet(setSpecPath(), setReplicasPath())

	requireStringSliceEqual(t, setPathStrings(set.Paths()), []string{
		"$.spec",
		"$.spec.replicas",
	})
}

func TestNewSetRejectsInvalidPath(t *testing.T) {
	_, err := NewSet(Path{elements: []Element{{kind: ElementField}}})

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorIs(t, err, ErrInvalidElement)
}

func TestNewSetDoesNotRetainInputSlice(t *testing.T) {
	input := []Path{setReplicasPath(), setImagePath()}
	set, err := NewSet(input...)
	requireNoError(t, err)

	input[0] = Root().Field(testField("status"))

	requireEqual(t, set.Has(setReplicasPath()), true)
	requireEqual(t, set.Has(Root().Field(testField("status"))), false)
}

func TestNewSetDoesNotRetainInputPathElements(t *testing.T) {
	input := setReplicasPath()
	set, err := NewSet(input)
	requireNoError(t, err)

	input.elements[0] = testFieldElement("status")

	requireEqual(t, set.Has(setReplicasPath()), true)
	requireEqual(t, set.Has(Root().Field(testField("status")).Field(testField("replicas"))), false)
}

func TestMustSetPanicsOnInvalidPath(t *testing.T) {
	requirePanic(t, func() {
		MustSet(Path{elements: []Element{{kind: ElementField}}})
	})
}

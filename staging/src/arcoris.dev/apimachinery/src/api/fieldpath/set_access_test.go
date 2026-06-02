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

func TestSetPathsReturnsDetachedSlice(t *testing.T) {
	set := MustSet(setReplicasPath(), setImagePath())
	paths := set.Paths()

	paths[0] = RootPath().Field("status")

	requireEqual(t, set.Has(RootPath().Field("status")), false)
	requireEqual(t, set.Len(), 2)
}

func TestSetPathsReturnsDetachedPathElements(t *testing.T) {
	set := MustSet(setReplicasPath())
	paths := set.Paths()

	paths[0].elements[0] = FieldElement("status")

	requireEqual(t, set.Has(setReplicasPath()), true)
	requireEqual(t, set.Has(RootPath().Field("status").Field("replicas")), false)
}

func TestSetHas(t *testing.T) {
	set := MustSet(setReplicasPath(), setImagePath(), setLabelPath())

	requireEqual(t, set.Has(setReplicasPath()), true)
	requireEqual(t, set.Has(RootPath().Field("status")), false)
}

func TestSetHasUsesStructuralPaths(t *testing.T) {
	field := RootPath().Field("api-version")
	key := RootPath().Key("api-version")
	set := MustSet(field)

	requireEqual(t, set.Has(field), true)
	requireEqual(t, set.Has(key), false)
}

func TestSetHasOnEmptySet(t *testing.T) {
	requireEqual(t, EmptySet().Has(RootPath()), false)
}

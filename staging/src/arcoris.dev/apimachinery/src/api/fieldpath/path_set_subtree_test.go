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

func TestPathSetHasDescendant(t *testing.T) {
	set := MustPathSet(
		RootPath().Field("spec"),
		RootPath().Field("spec").Field("replicas"),
		RootPath().Field("status"),
	)

	requireEqual(t, set.HasDescendant(RootPath().Field("spec")), true)
	requireEqual(t, set.HasDescendant(RootPath().Field("status")), false)
	requireEqual(t, set.HasDescendant(RootPath().Field("spec").Field("replicas")), false)
}

func TestPathSetIntersectsSubtree(t *testing.T) {
	set := MustPathSet(
		RootPath().Field("spec").Field("replicas"),
		RootPath().Field("status"),
	)

	requireEqual(t, set.IntersectsSubtree(RootPath().Field("spec")), true)
	requireEqual(t, set.IntersectsSubtree(RootPath().Field("metadata")), false)
}

func TestPathSetRemoveDescendants(t *testing.T) {
	set := MustPathSet(
		RootPath().Field("spec"),
		RootPath().Field("spec").Field("replicas"),
		RootPath().Field("spec").Field("template").Field("image"),
		RootPath().Field("status"),
	)

	filtered := set.RemoveDescendants(RootPath().Field("spec"))
	paths := filtered.Paths()

	requireEqual(t, len(paths), 2)
	requireEqual(t, paths[0].String(), "$.spec")
	requireEqual(t, paths[1].String(), "$.status")
}

func TestPathSetRemoveDescendantsReturnsDetachedSet(t *testing.T) {
	original := MustPathSet(
		RootPath().Field("spec"),
		RootPath().Field("spec").Field("replicas"),
	)

	filtered := original.RemoveDescendants(RootPath().Field("spec"))

	requireEqual(t, original.Len(), 2)
	requireEqual(t, filtered.Len(), 1)
	requireEqual(t, original.Paths()[1].String(), "$.spec.replicas")
}

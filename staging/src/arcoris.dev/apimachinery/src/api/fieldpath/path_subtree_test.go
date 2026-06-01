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

func TestPathIsDescendantOf(t *testing.T) {
	ancestor := RootPath().Field("spec")
	descendant := RootPath().Field("spec").Field("replicas")

	requireEqual(t, descendant.IsDescendantOf(ancestor), true)
	requireEqual(t, ancestor.IsDescendantOf(ancestor), false)
	requireEqual(t, RootPath().Field("status").IsDescendantOf(ancestor), false)
}

func TestPathIntersectsSubtree(t *testing.T) {
	left := RootPath().Field("spec")
	right := RootPath().Field("spec").Field("replicas")
	other := RootPath().Field("status")

	requireEqual(t, left.IntersectsSubtree(right), true)
	requireEqual(t, right.IntersectsSubtree(left), true)
	requireEqual(t, left.IntersectsSubtree(other), false)
}

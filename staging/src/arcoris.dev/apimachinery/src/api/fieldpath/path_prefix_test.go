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

func TestRootParentIsAbsent(t *testing.T) {
	_, ok := RootPath().Parent()
	requireEqual(t, ok, false)
}

func TestPathParent(t *testing.T) {
	parent, ok := RootPath().Field("spec").Field("replicas").Parent()

	requireEqual(t, ok, true)
	requireEqual(t, parent.String(), "$.spec")
}

func TestPathHasPrefix(t *testing.T) {
	path := RootPath().Field("spec").Field("replicas")

	requireEqual(t, path.HasPrefix(RootPath()), true)
	requireEqual(t, path.HasPrefix(RootPath().Field("spec")), true)
	requireEqual(t, path.HasPrefix(RootPath().Field("status")), false)
}

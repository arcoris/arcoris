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

func TestNewPathSetSortsAndDeduplicates(t *testing.T) {
	set, err := NewPathSet(
		RootPath().Field("status"),
		RootPath().Field("spec").Field("replicas"),
		RootPath().Field("status"),
		RootPath().Field("spec"),
	)
	requireNoError(t, err)

	paths := set.Paths()
	requireEqual(t, len(paths), 3)
	requireEqual(t, paths[0].String(), "$.spec")
	requireEqual(t, paths[1].String(), "$.spec.replicas")
	requireEqual(t, paths[2].String(), "$.status")
	requireEqual(t, cap(set.paths), len(set.paths))
}

func TestNewPathSetRejectsInvalidPath(t *testing.T) {
	_, err := NewPathSet(RootPath().Field(""), RootPath().Field("status"))

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorIs(t, err, ErrInvalidElement)
}

func TestPathSetPathsReturnsClone(t *testing.T) {
	set := MustPathSet(RootPath().Field("spec"))
	paths := set.Paths()

	paths[0] = RootPath().Field("status")

	requireEqual(t, set.Paths()[0].String(), "$.spec")
}

func TestMustPathSetPanicsOnInvalidPath(t *testing.T) {
	requirePanic(t, func() {
		MustPathSet(Path{}.Field(""))
	})
}

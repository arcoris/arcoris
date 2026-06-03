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

func setSpecPath() Path {
	return RootPath().Field("spec")
}

func setReplicasPath() Path {
	return RootPath().Field("spec").Field("replicas")
}

func setImagePath() Path {
	return RootPath().Field("spec").Field("image")
}

func setLabelPath() Path {
	return RootPath().Field("metadata").Field("labels").Key("app")
}

func setIndexPath() Path {
	return RootPath().Field("items").Index(0)
}

func setReadyStatusPath() Path {
	ready := MustSelector(
		NewSelectorEntry("type", StringLiteral("Ready")),
	)

	return RootPath().
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

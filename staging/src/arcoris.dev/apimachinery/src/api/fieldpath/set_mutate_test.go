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

func TestSetInsert(t *testing.T) {
	original := MustSet(setReplicasPath())
	updated := original.Insert(setImagePath()).Insert(setReplicasPath())

	requireEqual(t, original.Len(), 1)
	requireEqual(t, updated.Len(), 2)
	requireStringSliceEqual(t, setPathStrings(updated.Paths()), []string{
		"$.spec.image",
		"$.spec.replicas",
	})
}

func TestSetInsertPanicsOnInvalidPath(t *testing.T) {
	requirePanic(t, func() {
		EmptySet().Insert(Path{elements: []Element{{kind: ElementField}}})
	})
}

func TestSetDelete(t *testing.T) {
	original := MustSet(setImagePath(), setReplicasPath())
	updated := original.Delete(setImagePath())
	missingDeleted := updated.Delete(Root().Field(testField("status")))

	requireEqual(t, original.Len(), 2)
	requireEqual(t, updated.Len(), 1)
	requireEqual(t, updated.Has(setImagePath()), false)
	requireEqual(t, updated.Has(setReplicasPath()), true)
	requireEqual(t, missingDeleted.Equal(updated), true)
}

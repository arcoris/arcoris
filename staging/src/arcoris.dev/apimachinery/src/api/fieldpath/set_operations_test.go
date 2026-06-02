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

func TestSetUnion(t *testing.T) {
	left := MustSet(setImagePath(), setReplicasPath())
	right := MustSet(setLabelPath(), setReplicasPath(), setReadyStatusPath())

	union := left.Union(right)

	requireStringSliceEqual(t, setPathStrings(union.Paths()), []string{
		`$.conditions[{"type":"Ready"}].status`,
		`$.metadata.labels["app"]`,
		"$.spec.image",
		"$.spec.replicas",
	})
	requireEqual(t, left.Len(), 2)
	requireEqual(t, right.Len(), 3)
}

func TestSetIntersection(t *testing.T) {
	left := MustSet(setImagePath(), setReplicasPath(), setReadyStatusPath())
	right := MustSet(setLabelPath(), setReplicasPath(), setReadyStatusPath())

	intersection := left.Intersection(right)

	requireStringSliceEqual(t, setPathStrings(intersection.Paths()), []string{
		`$.conditions[{"type":"Ready"}].status`,
		"$.spec.replicas",
	})
}

func TestSetDifference(t *testing.T) {
	left := MustSet(setImagePath(), setReplicasPath(), setReadyStatusPath())
	right := MustSet(setReplicasPath())

	difference := left.Difference(right)

	requireStringSliceEqual(t, setPathStrings(difference.Paths()), []string{
		`$.conditions[{"type":"Ready"}].status`,
		"$.spec.image",
	})
}

func TestSetOperationsDoNotMutateReceivers(t *testing.T) {
	left := MustSet(setImagePath())
	right := MustSet(setReplicasPath())

	_ = left.Union(right)
	_ = left.Intersection(right)
	_ = left.Difference(right)

	requireStringSliceEqual(t, setPathStrings(left.Paths()), []string{"$.spec.image"})
	requireStringSliceEqual(t, setPathStrings(right.Paths()), []string{"$.spec.replicas"})
}

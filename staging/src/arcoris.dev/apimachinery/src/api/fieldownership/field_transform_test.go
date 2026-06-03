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

package fieldownership

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestUnionTransform(t *testing.T) {
	got := unionTransform(set(imagePath()), set(replicasPath()))

	requireSet(t, got, "$.spec.image", "$.spec.replicas")
}

func TestUnionTransformKeepsEmptyIdentity(t *testing.T) {
	fields := set(imagePath())

	requireEqual(t, unionTransform(fieldpath.EmptySet(), fields).Equal(fields), true)
	requireEqual(t, unionTransform(fields, fieldpath.EmptySet()).Equal(fields), true)
}

func TestRemoveExactTransform(t *testing.T) {
	got := removeExactTransform(set(specPath(), replicasPath()), set(specPath()))

	requireSet(t, got, "$.spec.replicas")
}

func TestRemoveExactTransformDoesNotRemoveOverlaps(t *testing.T) {
	got := removeExactTransform(set(specPath(), replicasPath()), set(replicasPath()))

	requireSet(t, got, "$.spec")
}

func TestRemoveOverlapTransform(t *testing.T) {
	got := removeOverlapTransform(set(specPath(), namePath()), set(replicasPath()))

	requireSet(t, got, "$.metadata.name")
}

func TestOverlapsAny(t *testing.T) {
	requireEqual(t, overlapsAny(specPath(), set(replicasPath())), true)
	requireEqual(t, overlapsAny(namePath(), set(replicasPath())), false)
}

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

package objectindex

import (
	"testing"

	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/objectquery"
)

func TestPlanAnnotationsNarrowsExistsEqualsAndIn(t *testing.T) {
	idx := Build(testItems())
	tests := []struct {
		name string
		req  objectquery.AnnotationRequirement
		want []int
	}{
		{
			name: "exists",
			req:  mustAnnotationExists(t, "zone"),
			want: []int{0, 1, 2},
		},
		{
			name: "equals",
			req:  mustAnnotationEquals(t, "team", "core"),
			want: []int{0, 2},
		},
		{
			name: "in",
			req:  mustAnnotationIn(t, "team", "core", "tools"),
			want: []int{0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set, ok := idx.annotationRequirementCandidates(tt.req)
			if !ok {
				t.Fatal("annotationRequirementCandidates ok = false; want true")
			}
			requireCandidateMembership(t, set, len(testItems()), tt.want...)
		})
	}
}

func TestPlanAnnotationsLeavesNegativeOperatorsResidualOnly(t *testing.T) {
	idx := Build(testItems())
	for _, req := range []objectquery.AnnotationRequirement{
		mustAnnotationDoesNotExist(t, "team"),
		mustAnnotationNotEquals(t, "team", "core"),
		mustAnnotationNotIn(t, "team", "core", "tools"),
	} {
		if _, ok := idx.annotationRequirementCandidates(req); ok {
			t.Fatalf("annotationRequirementCandidates(%s) ok = true; want false", req.Operator())
		}
	}
}

func TestPlanAnnotationsIntersectsRequirements(t *testing.T) {
	requirements := mustAnnotationSelector(
		t,
		mustAnnotationIn(t, "team", "core", "tools"),
		mustAnnotationEquals(t, "zone", "west"),
	).Requirements()

	plan := Build(testItems()).planAnnotations(candidatePlan{}, requirements)

	requirePlanIncludesOnly(t, plan, len(testItems()), 1, 2)
}

func TestPlanAnnotationInCandidatesUsesValueUnion(t *testing.T) {
	idx := Build(testItems())

	set := idx.annotationInCandidates(annotations.Key("team"), []string{"core", "tools"})

	requireCandidateMembership(t, set, len(testItems()), 0, 1, 2)
}

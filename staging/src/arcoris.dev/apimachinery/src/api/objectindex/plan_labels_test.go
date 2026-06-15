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

	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectquery"
)

func TestPlanLabelsNarrowsExistsEqualsAndIn(t *testing.T) {
	idx := Build(testItems())
	tests := []struct {
		name string
		req  objectquery.LabelRequirement
		want []int
	}{
		{
			name: "exists",
			req:  mustLabelExists(t, "tier"),
			want: []int{0, 1, 2},
		},
		{
			name: "equals",
			req:  mustLabelEquals(t, "tier", "backend"),
			want: []int{0, 1},
		},
		{
			name: "in",
			req:  mustLabelIn(t, "env", "qa", "prod"),
			want: []int{0, 1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set, ok := idx.labelRequirementCandidates(tt.req)
			if !ok {
				t.Fatal("labelRequirementCandidates ok = false; want true")
			}
			requireCandidateMembership(t, set, len(testItems()), tt.want...)
		})
	}
}

func TestPlanLabelsLeavesNegativeOperatorsResidualOnly(t *testing.T) {
	idx := Build(testItems())
	for _, req := range []objectquery.LabelRequirement{
		mustLabelDoesNotExist(t, "env"),
		mustLabelNotEquals(t, "env", "prod"),
		mustLabelNotIn(t, "env", "prod", "qa"),
	} {
		if _, ok := idx.labelRequirementCandidates(req); ok {
			t.Fatalf("labelRequirementCandidates(%s) ok = true; want false", req.Operator())
		}
	}
}

func TestPlanLabelsIntersectsRequirements(t *testing.T) {
	requirements := mustLabelSelector(
		t,
		mustLabelIn(t, "env", "prod", "qa"),
		mustLabelEquals(t, "tier", "backend"),
	).Requirements()

	plan := Build(testItems()).planLabels(candidatePlan{}, requirements)

	requirePlanIncludesOnly(t, plan, len(testItems()), 0, 1)
}

func TestPlanLabelInCandidatesUsesValueUnion(t *testing.T) {
	idx := Build(testItems())

	set := idx.labelInCandidates(labels.Key("env"), []string{"prod", "qa"})

	requireCandidateMembership(t, set, len(testItems()), 0, 1, 2, 3)
}

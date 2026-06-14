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

package objectquery

import "testing"

func TestAnnotationSelectorCanonicalizesRequirementOrderAndDuplicates(t *testing.T) {
	selector := mustAnnotationSelector(
		t,
		mustAnnotationEquals(t, "team", "platform"),
		mustAnnotationIn(t, "note", "qa rollout", "prod rollout"),
		mustAnnotationIn(t, "note", "prod rollout", "qa rollout"),
	)

	if len(selector.requirements) != 2 {
		t.Fatalf("len = %d; want 2", len(selector.requirements))
	}
	requireRequirement(t, selector.requirements[0].req, "note", OperatorIn, "prod rollout", "qa rollout")
	requireRequirement(t, selector.requirements[1].req, "team", OperatorEquals, "platform")
}

func TestAnnotationSelectorAllowsDifferentRequirementsForSameKey(t *testing.T) {
	selector := mustAnnotationSelector(
		t,
		mustAnnotationIn(t, "note", "prod rollout", "qa rollout"),
		mustAnnotationNotEquals(t, "note", "qa rollout"),
	)

	if len(selector.requirements) != 2 {
		t.Fatalf("len = %d; want 2", len(selector.requirements))
	}
	if !selector.match(testItem("system", "worker", nil, map[string]string{"note": "prod rollout"})) {
		t.Fatal("selector did not match prod rollout")
	}
	if selector.match(testItem("system", "worker", nil, map[string]string{"note": "qa rollout"})) {
		t.Fatal("selector matched qa rollout")
	}
}

func TestAnnotationSelectorDoesNotMutateInputRequirements(t *testing.T) {
	req := AnnotationRequirement{req: metadataRequirement{key: "note", op: OperatorIn, values: []string{"qa rollout", "prod rollout"}}}

	selector := mustAnnotationSelector(t, req)
	selector.requirements[0].req.values[0] = "mutated"

	requireRequirement(t, req.req, "note", OperatorIn, "qa rollout", "prod rollout")
}

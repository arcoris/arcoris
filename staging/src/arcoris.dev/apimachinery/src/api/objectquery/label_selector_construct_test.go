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

func TestLabelSelectorCanonicalizesRequirementOrderAndDuplicates(t *testing.T) {
	selector := mustLabelSelector(
		t,
		mustLabelEquals(t, "tier", "backend"),
		mustLabelIn(t, "env", "qa", "prod"),
		mustLabelIn(t, "env", "prod", "qa"),
	)

	if len(selector.requirements) != 2 {
		t.Fatalf("len = %d; want 2", len(selector.requirements))
	}
	requireRequirement(t, selector.requirements[0].req, "env", OperatorIn, "prod", "qa")
	requireRequirement(t, selector.requirements[1].req, "tier", OperatorEquals, "backend")
}

func TestLabelSelectorAllowsDifferentRequirementsForSameKey(t *testing.T) {
	selector := mustLabelSelector(
		t,
		mustLabelIn(t, "env", "prod", "qa"),
		mustLabelNotEquals(t, "env", "qa"),
	)

	if len(selector.requirements) != 2 {
		t.Fatalf("len = %d; want 2", len(selector.requirements))
	}
	if !selector.match(testItem("system", "worker", map[string]string{"env": "prod"}, nil)) {
		t.Fatal("selector did not match env=prod")
	}
	if selector.match(testItem("system", "worker", map[string]string{"env": "qa"}, nil)) {
		t.Fatal("selector matched env=qa")
	}
}

func TestLabelSelectorDoesNotMutateInputRequirements(t *testing.T) {
	req := LabelRequirement{req: metadataRequirement{key: "env", op: OperatorIn, values: []string{"qa", "prod"}}}

	selector := mustLabelSelector(t, req)
	selector.requirements[0].req.values[0] = "mutated"

	requireRequirement(t, req.req, "env", OperatorIn, "qa", "prod")
}

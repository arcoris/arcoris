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

func TestLabelSelectorAndsRequirements(t *testing.T) {
	selector := mustLabelSelector(
		t,
		mustLabelEquals(t, "env", "prod"),
		mustLabelEquals(t, "tier", "backend"),
	)

	if !selector.match(testItem("system", "worker", map[string]string{"env": "prod", "tier": "backend"}, nil)) {
		t.Fatal("selector did not match both requirements")
	}
	if selector.match(testItem("system", "worker", map[string]string{"env": "prod"}, nil)) {
		t.Fatal("selector matched one requirement")
	}
}

func TestLabelRequirementMatching(t *testing.T) {
	item := testItem("system", "worker", map[string]string{"env": "prod"}, nil)
	tests := []struct {
		name  string
		req   LabelRequirement
		match bool
	}{
		{name: "exists present", req: mustLabelExists(t, "env"), match: true},
		{name: "exists absent", req: mustLabelExists(t, "tier")},
		{name: "does not exist absent", req: mustLabelDoesNotExist(t, "tier"), match: true},
		{name: "does not exist present", req: mustLabelDoesNotExist(t, "env")},
		{name: "equals exact", req: mustLabelEquals(t, "env", "prod"), match: true},
		{name: "equals absent", req: mustLabelEquals(t, "tier", "backend")},
		{name: "equals different", req: mustLabelEquals(t, "env", "qa")},
		{name: "not equals absent", req: mustLabelNotEquals(t, "tier", "backend"), match: true},
		{name: "not equals different", req: mustLabelNotEquals(t, "env", "qa"), match: true},
		{name: "not equals equal", req: mustLabelNotEquals(t, "env", "prod")},
		{name: "in member", req: mustLabelIn(t, "env", "prod", "qa"), match: true},
		{name: "in absent", req: mustLabelIn(t, "tier", "backend")},
		{name: "in outside", req: mustLabelIn(t, "env", "qa")},
		{name: "not in absent", req: mustLabelNotIn(t, "tier", "backend"), match: true},
		{name: "not in outside", req: mustLabelNotIn(t, "env", "qa"), match: true},
		{name: "not in member", req: mustLabelNotIn(t, "env", "prod", "qa")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := mustLabelSelector(t, tt.req)
			if got := selector.match(item); got != tt.match {
				t.Fatalf("match = %v; want %v", got, tt.match)
			}
		})
	}
}

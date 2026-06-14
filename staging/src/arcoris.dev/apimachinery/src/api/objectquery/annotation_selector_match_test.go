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

func TestAnnotationSelectorAndsRequirements(t *testing.T) {
	selector := mustAnnotationSelector(
		t,
		mustAnnotationEquals(t, "note", "prod rollout"),
		mustAnnotationEquals(t, "team", "platform"),
	)

	if !selector.match(testItem("system", "worker", nil, map[string]string{"note": "prod rollout", "team": "platform"})) {
		t.Fatal("selector did not match both requirements")
	}
	if selector.match(testItem("system", "worker", nil, map[string]string{"note": "prod rollout"})) {
		t.Fatal("selector matched one requirement")
	}
}

func TestAnnotationRequirementMatching(t *testing.T) {
	item := testItem("system", "worker", nil, map[string]string{"note": "prod rollout"})
	tests := []struct {
		name  string
		req   AnnotationRequirement
		match bool
	}{
		{name: "exists present", req: mustAnnotationExists(t, "note"), match: true},
		{name: "exists absent", req: mustAnnotationExists(t, "owner")},
		{name: "does not exist absent", req: mustAnnotationDoesNotExist(t, "owner"), match: true},
		{name: "does not exist present", req: mustAnnotationDoesNotExist(t, "note")},
		{name: "equals exact", req: mustAnnotationEquals(t, "note", "prod rollout"), match: true},
		{name: "equals absent", req: mustAnnotationEquals(t, "owner", "team")},
		{name: "equals different", req: mustAnnotationEquals(t, "note", "qa rollout")},
		{name: "not equals absent", req: mustAnnotationNotEquals(t, "owner", "team"), match: true},
		{name: "not equals different", req: mustAnnotationNotEquals(t, "note", "qa rollout"), match: true},
		{name: "not equals equal", req: mustAnnotationNotEquals(t, "note", "prod rollout")},
		{name: "in member", req: mustAnnotationIn(t, "note", "prod rollout", "qa rollout"), match: true},
		{name: "in absent", req: mustAnnotationIn(t, "owner", "team")},
		{name: "in outside", req: mustAnnotationIn(t, "note", "qa rollout")},
		{name: "not in absent", req: mustAnnotationNotIn(t, "owner", "team"), match: true},
		{name: "not in outside", req: mustAnnotationNotIn(t, "note", "qa rollout"), match: true},
		{name: "not in member", req: mustAnnotationNotIn(t, "note", "prod rollout", "qa rollout")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := mustAnnotationSelector(t, tt.req)
			if got := selector.match(item); got != tt.match {
				t.Fatalf("match = %v; want %v", got, tt.match)
			}
		})
	}
}

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

func TestAnnotationRequirementConstructors(t *testing.T) {
	tests := []struct {
		name   string
		build  func() (AnnotationRequirement, error)
		key    string
		op     Operator
		values []string
	}{
		{name: "exists", build: func() (AnnotationRequirement, error) { return AnnotationExists("note") }, key: "note", op: OperatorExists},
		{name: "does not exist", build: func() (AnnotationRequirement, error) { return AnnotationDoesNotExist("note") }, key: "note", op: OperatorDoesNotExist},
		{name: "equals", build: func() (AnnotationRequirement, error) { return AnnotationEquals("note", "prod rollout") }, key: "note", op: OperatorEquals, values: []string{"prod rollout"}},
		{name: "not equals", build: func() (AnnotationRequirement, error) { return AnnotationNotEquals("note", "prod rollout") }, key: "note", op: OperatorNotEquals, values: []string{"prod rollout"}},
		{name: "in", build: func() (AnnotationRequirement, error) {
			return AnnotationIn("note", "qa rollout", "prod rollout", "qa rollout")
		}, key: "note", op: OperatorIn, values: []string{"prod rollout", "qa rollout"}},
		{name: "not in", build: func() (AnnotationRequirement, error) {
			return AnnotationNotIn("note", "qa rollout", "prod rollout", "qa rollout")
		}, key: "note", op: OperatorNotIn, values: []string{"prod rollout", "qa rollout"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.build()
			requireNoError(t, err)
			requireRequirement(t, req.req, tt.key, tt.op, tt.values...)
		})
	}
}

func TestAnnotationRequirementRejectsInvalidShape(t *testing.T) {
	tests := []struct {
		name string
		req  AnnotationRequirement
	}{
		{name: "empty key", req: AnnotationRequirement{req: metadataRequirement{op: OperatorExists}}},
		{name: "invalid operator", req: AnnotationRequirement{req: metadataRequirement{key: "note", op: Operator(255)}}},
		{name: "exists values", req: AnnotationRequirement{req: metadataRequirement{key: "note", op: OperatorExists, values: []string{"prod"}}}},
		{name: "does not exist values", req: AnnotationRequirement{req: metadataRequirement{key: "note", op: OperatorDoesNotExist, values: []string{"prod"}}}},
		{name: "equals no value", req: AnnotationRequirement{req: metadataRequirement{key: "note", op: OperatorEquals}}},
		{name: "equals too many values", req: AnnotationRequirement{req: metadataRequirement{key: "note", op: OperatorEquals, values: []string{"prod", "qa"}}}},
		{name: "not equals no value", req: AnnotationRequirement{req: metadataRequirement{key: "note", op: OperatorNotEquals}}},
		{name: "in empty", req: AnnotationRequirement{req: metadataRequirement{key: "note", op: OperatorIn}}},
		{name: "not in empty", req: AnnotationRequirement{req: metadataRequirement{key: "note", op: OperatorNotIn}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAnnotationSelector(tt.req)
			requireErrorIs(t, err, ErrInvalidRequirement)
		})
	}
}

func TestAnnotationRequirementDoesNotRetainCallerValues(t *testing.T) {
	values := []string{"qa rollout", "prod rollout"}
	req, err := AnnotationIn("note", values...)
	requireNoError(t, err)
	values[0] = "mutated"

	requireRequirement(t, req.req, "note", OperatorIn, "prod rollout", "qa rollout")
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

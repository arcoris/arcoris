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

func TestLabelRequirementConstructors(t *testing.T) {
	tests := []struct {
		name   string
		build  func() (LabelRequirement, error)
		key    string
		op     Operator
		values []string
	}{
		{name: "exists", build: func() (LabelRequirement, error) { return LabelExists("env") }, key: "env", op: OperatorExists},
		{name: "does not exist", build: func() (LabelRequirement, error) { return LabelDoesNotExist("env") }, key: "env", op: OperatorDoesNotExist},
		{name: "equals", build: func() (LabelRequirement, error) { return LabelEquals("env", "prod") }, key: "env", op: OperatorEquals, values: []string{"prod"}},
		{name: "not equals", build: func() (LabelRequirement, error) { return LabelNotEquals("env", "prod") }, key: "env", op: OperatorNotEquals, values: []string{"prod"}},
		{name: "in", build: func() (LabelRequirement, error) { return LabelIn("env", "qa", "prod", "qa") }, key: "env", op: OperatorIn, values: []string{"prod", "qa"}},
		{name: "not in", build: func() (LabelRequirement, error) { return LabelNotIn("env", "qa", "prod", "qa") }, key: "env", op: OperatorNotIn, values: []string{"prod", "qa"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.build()
			requireNoError(t, err)
			requireRequirement(t, req.req, tt.key, tt.op, tt.values...)
		})
	}
}

func TestLabelRequirementRejectsInvalidShape(t *testing.T) {
	tests := []struct {
		name string
		req  LabelRequirement
	}{
		{name: "empty key", req: LabelRequirement{req: metadataRequirement{op: OperatorExists}}},
		{name: "invalid operator", req: LabelRequirement{req: metadataRequirement{key: "env", op: Operator(255)}}},
		{name: "exists values", req: LabelRequirement{req: metadataRequirement{key: "env", op: OperatorExists, values: []string{"prod"}}}},
		{name: "does not exist values", req: LabelRequirement{req: metadataRequirement{key: "env", op: OperatorDoesNotExist, values: []string{"prod"}}}},
		{name: "equals no value", req: LabelRequirement{req: metadataRequirement{key: "env", op: OperatorEquals}}},
		{name: "equals too many values", req: LabelRequirement{req: metadataRequirement{key: "env", op: OperatorEquals, values: []string{"prod", "qa"}}}},
		{name: "not equals no value", req: LabelRequirement{req: metadataRequirement{key: "env", op: OperatorNotEquals}}},
		{name: "in empty", req: LabelRequirement{req: metadataRequirement{key: "env", op: OperatorIn}}},
		{name: "not in empty", req: LabelRequirement{req: metadataRequirement{key: "env", op: OperatorNotIn}}},
		{name: "invalid value", req: LabelRequirement{req: metadataRequirement{key: "env", op: OperatorEquals, values: []string{"bad value"}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLabelSelector(tt.req)
			requireErrorIs(t, err, ErrInvalidRequirement)
		})
	}
}

func TestLabelRequirementDoesNotRetainCallerValues(t *testing.T) {
	values := []string{"qa", "prod"}
	req, err := LabelIn("env", values...)
	requireNoError(t, err)
	values[0] = "mutated"

	requireRequirement(t, req.req, "env", OperatorIn, "prod", "qa")
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

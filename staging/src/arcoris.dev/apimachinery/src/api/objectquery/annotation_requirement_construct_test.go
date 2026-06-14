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

func TestAnnotationRequirementDoesNotRetainCallerValues(t *testing.T) {
	values := []string{"qa rollout", "prod rollout"}
	req, err := AnnotationIn("note", values...)
	requireNoError(t, err)
	values[0] = "mutated"

	requireRequirement(t, req.req, "note", OperatorIn, "prod rollout", "qa rollout")
}

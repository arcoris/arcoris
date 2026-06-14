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

func TestLabelRequirementDoesNotRetainCallerValues(t *testing.T) {
	values := []string{"qa", "prod"}
	req, err := LabelIn("env", values...)
	requireNoError(t, err)
	values[0] = "mutated"

	requireRequirement(t, req.req, "env", OperatorIn, "prod", "qa")
}

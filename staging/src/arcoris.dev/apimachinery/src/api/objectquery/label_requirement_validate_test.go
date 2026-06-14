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

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

func TestMetadataRequirementValidateRejectsInvalidOperator(t *testing.T) {
	req := metadataRequirement{key: "env", op: Operator(255)}
	err := req.validate("query.labels.requirements[0]", validateLabelKey, validateLabelValue)

	requireErrorIs(t, err, ErrInvalidRequirement)
	requireErrorIs(t, err, ErrInvalidOperator)
}

func TestValidateMetadataValueCount(t *testing.T) {
	tests := []struct {
		name  string
		op    Operator
		count int
		valid bool
	}{
		{name: "exists none", op: OperatorExists, valid: true},
		{name: "exists one", op: OperatorExists, count: 1},
		{name: "equals one", op: OperatorEquals, count: 1, valid: true},
		{name: "equals none", op: OperatorEquals},
		{name: "in one", op: OperatorIn, count: 1, valid: true},
		{name: "in none", op: OperatorIn},
		{name: "invalid operator", op: Operator(255)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMetadataValueCount("query.labels.values", tt.op, tt.count)
			if tt.valid {
				requireNoError(t, err)
				return
			}
			requireErrorIs(t, err, ErrInvalidRequirement)
		})
	}
}

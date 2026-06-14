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

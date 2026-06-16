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

import (
	"testing"

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

// TestMetadataTerms verifies label and annotation matching semantics.
func TestMetadataTerms(t *testing.T) {
	items := testItems()
	tests := []struct {
		name  string
		query Query
		want  []metaidentity.Name
	}{
		{name: "label exists", query: mustQ(LabelExists("env")), want: []metaidentity.Name{"worker-1", "worker-2", "worker-3"}},
		{name: "label does not exist", query: mustQ(LabelDoesNotExist("env")), want: []metaidentity.Name{"worker-4"}},
		{name: "label equals", query: mustQ(LabelEquals("env", "prod")), want: []metaidentity.Name{"worker-1", "worker-3"}},
		{name: "label not equals", query: mustQ(LabelNotEquals("env", "prod")), want: []metaidentity.Name{"worker-2", "worker-4"}},
		{name: "label in", query: mustQ(LabelIn("env", "qa", "prod", "qa")), want: []metaidentity.Name{"worker-1", "worker-2", "worker-3"}},
		{name: "label not in", query: mustQ(LabelNotIn("env", "prod")), want: []metaidentity.Name{"worker-2", "worker-4"}},
		{name: "annotation exists", query: mustQ(AnnotationExists("team")), want: []metaidentity.Name{"worker-1", "worker-2"}},
		{name: "annotation equals", query: mustQ(AnnotationEquals("team", "core")), want: []metaidentity.Name{"worker-1"}},
		{name: "annotation not in", query: mustQ(AnnotationNotIn("team", "core")), want: []metaidentity.Name{"worker-2", "worker-3", "worker-4"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireNames(t, mustPredicate(t, tt.query).Filter(items), tt.want...)
		})
	}
}

// TestMetadataValidationErrors verifies key, arity, and value validation.
func TestMetadataValidationErrors(t *testing.T) {
	_, err := LabelEquals("", "prod")
	requireErrorIs(t, err, ErrInvalidTerm)

	_, err = LabelIn("env")
	requireErrorIs(t, err, ErrInvalidTerm)

	_, err = LabelEquals("env", "bad value")
	requireErrorIs(t, err, ErrInvalidTerm)
}

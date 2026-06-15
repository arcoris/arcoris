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

func TestQueryZeroValue(t *testing.T) {
	var query Query
	if !query.Identity.IsZero() || !query.Labels.IsZero() || !query.Annotations.IsZero() {
		t.Fatal("zero Query contains non-zero section")
	}
}

func TestQueryValidateZero(t *testing.T) {
	requireNoError(t, (Query{}).Validate())
}

func TestQueryValidateInvalidIdentity(t *testing.T) {
	err := (Query{
		Identity: IdentitySelector{Name: NameRequirement{set: true}},
	}).Validate()

	requireErrorIs(t, err, ErrInvalidQuery)
}

func TestQueryValidateInvalidLabels(t *testing.T) {
	err := (Query{
		Labels: LabelSelector{requirements: []LabelRequirement{{req: metadataRequirement{key: "env", op: OperatorIn}}}},
	}).Validate()

	requireErrorIs(t, err, ErrInvalidQuery)
}

func TestQueryValidateInvalidAnnotations(t *testing.T) {
	err := (Query{
		Annotations: AnnotationSelector{requirements: []AnnotationRequirement{{req: metadataRequirement{key: "note", op: OperatorIn}}}},
	}).Validate()

	requireErrorIs(t, err, ErrInvalidQuery)
}

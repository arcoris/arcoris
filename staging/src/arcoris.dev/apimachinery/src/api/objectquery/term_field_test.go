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

	"arcoris.dev/apimachinery/api/value"
)

// TestFieldConstructorsBuildFieldTerms verifies every public field constructor
// records the requested operator in the field domain.
func TestFieldConstructorsBuildFieldTerms(t *testing.T) {
	ref := fieldRef("spec.image")
	queries := []Query{
		mustQ(FieldExists(ref)),
		mustQ(FieldDoesNotExist(ref)),
		mustQ(FieldEquals(ref, value.StringValue("api"))),
		mustQ(FieldNotEquals(ref, value.StringValue("api"))),
		mustQ(FieldIn(ref, value.StringValue("api"))),
		mustQ(FieldNotIn(ref, value.StringValue("api"))),
		mustQ(FieldLessThan(ref, value.StringValue("b"))),
		mustQ(FieldLessOrEqual(ref, value.StringValue("api"))),
		mustQ(FieldGreaterThan(ref, value.StringValue("a"))),
		mustQ(FieldGreaterOrEqual(ref, value.StringValue("api"))),
		mustQ(FieldHasPrefix(ref, "a")),
		mustQ(FieldHasSuffix(ref, "i")),
		mustQ(FieldContains(ref, "p")),
	}

	for _, query := range queries {
		if query.expr.term.kind != termField {
			t.Fatalf("term kind = %v; want field", query.expr.term.kind)
		}
	}
}

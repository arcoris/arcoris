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

// TestPredicateMatchFieldOperators verifies selectable field evaluation,
// including missing/null and string/ordered operators.
func TestPredicateMatchFieldOperators(t *testing.T) {
	fields := mustFieldSet(t,
		selectable(fieldRef("spec.image"), value.KindString, Operators(OperatorEquals, OperatorNotEquals, OperatorIn, OperatorHasPrefix, OperatorContains)),
		selectable(fieldRef("spec.replicas"), value.KindInteger, Operators(OperatorGreaterThan, OperatorLessOrEqual)),
		selectable(fieldRef("spec.nullable"), value.KindString, Operators(OperatorEquals, OperatorNotEquals)),
	)
	items := testItems()

	tests := []struct {
		name  string
		query Query
		want  []string
	}{
		{name: "equals", query: mustQ(FieldEquals(fieldRef("spec.image"), value.StringValue("api"))), want: []string{"worker-1", "worker-2"}},
		{name: "in", query: mustQ(FieldIn(fieldRef("spec.image"), value.StringValue("api"), value.StringValue("web"))), want: []string{"worker-1", "worker-2", "worker-3"}},
		{name: "prefix", query: mustQ(FieldHasPrefix(fieldRef("spec.image"), "a")), want: []string{"worker-1", "worker-2"}},
		{name: "contains", query: mustQ(FieldContains(fieldRef("spec.image"), "e")), want: []string{"worker-3"}},
		{name: "greater", query: mustQ(FieldGreaterThan(fieldRef("spec.replicas"), value.Int64Value(2))), want: []string{"worker-1", "worker-3"}},
		{name: "lessOrEqual", query: mustQ(FieldLessOrEqual(fieldRef("spec.replicas"), value.Int64Value(1))), want: []string{"worker-2"}},
		{name: "missing negative", query: mustQ(FieldNotEquals(fieldRef("spec.image"), value.StringValue("api"))), want: []string{"worker-3", "worker-4"}},
		{name: "null equals", query: mustQ(FieldEquals(fieldRef("spec.nullable"), value.NullValue())), want: []string{"worker-1", "worker-2", "worker-3"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustPredicate(t, tt.query, WithSelectableFields(fields)).Filter(items)
			requireNamesFromStrings(t, got, tt.want...)
		})
	}
}

// TestPredicateMatchObservedField verifies observed selectable fields compile
// and evaluate against the observed payload, not desired state.
func TestPredicateMatchObservedField(t *testing.T) {
	item := testItem("system", "worker-1", 1, nil, nil, desiredRecord("api", "desired", 1))
	observed := value.MustRecordValue(
		value.MustRecordMember("status", value.MustRecordValue(
			value.MustRecordMember("phase", value.StringValue("Ready")),
		)),
	)
	item.State.Object.Observed = &observed
	ref := observedFieldRef("status.phase")
	fields := mustFieldSet(t, selectable(ref, value.KindString, Operators(OperatorEquals)))
	predicate := mustPredicate(t, mustQ(FieldEquals(ref, value.StringValue("Ready"))), WithSelectableFields(fields))

	if !predicate.Match(item) {
		t.Fatal("observed field predicate did not match")
	}
}

// BenchmarkPredicateMatchFieldTerm covers registered field evaluation.
func BenchmarkPredicateMatchFieldTerm(b *testing.B) {
	fields := mustFieldSet(b, selectable(fieldRef("spec.image"), value.KindString, Operators(OperatorEquals)))
	predicate := mustPredicate(b, mustQ(FieldEquals(fieldRef("spec.image"), value.StringValue("api"))), WithSelectableFields(fields))
	item := testItems()[0]

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = predicate.Match(item)
	}
}

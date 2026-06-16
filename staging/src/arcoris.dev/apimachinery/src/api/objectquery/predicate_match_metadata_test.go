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

// TestMatchMetadataTermAbsentKeySemantics verifies negative metadata operators
// intentionally match absent keys while positive membership does not.
func TestMatchMetadataTermAbsentKeySemantics(t *testing.T) {
	item := testItems()[3]
	tests := []struct {
		name string
		term term
		want bool
	}{
		{name: "exists", term: term{kind: termMetadata, metadataDomain: metadataLabels, metadataKey: "env", operator: OperatorExists}, want: false},
		{name: "doesNotExist", term: term{kind: termMetadata, metadataDomain: metadataLabels, metadataKey: "env", operator: OperatorDoesNotExist}, want: true},
		{name: "equals", term: term{kind: termMetadata, metadataDomain: metadataLabels, metadataKey: "env", operator: OperatorEquals, stringValues: []string{"prod"}}, want: false},
		{name: "notEquals", term: term{kind: termMetadata, metadataDomain: metadataLabels, metadataKey: "env", operator: OperatorNotEquals, stringValues: []string{"prod"}}, want: true},
		{name: "in", term: term{kind: termMetadata, metadataDomain: metadataLabels, metadataKey: "env", operator: OperatorIn, stringValues: []string{"prod"}}, want: false},
		{name: "notIn", term: term{kind: termMetadata, metadataDomain: metadataLabels, metadataKey: "env", operator: OperatorNotIn, stringValues: []string{"prod"}}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchMetadataTerm(tt.term, item); got != tt.want {
				t.Fatalf("matchMetadataTerm = %v; want %v", got, tt.want)
			}
		})
	}
}

// BenchmarkPredicateMatchMetadataOnly covers the common label/annotation path.
func BenchmarkPredicateMatchMetadataOnly(b *testing.B) {
	predicate := mustPredicate(b, mustQ(LabelEquals("env", "prod")))
	item := testItems()[0]

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = predicate.Match(item)
	}
}

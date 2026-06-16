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

// TestBooleanConstructors verifies public boolean constructor normalization.
func TestBooleanConstructors(t *testing.T) {
	envProd := mustQ(LabelEquals("env", "prod"))
	tierBackend := mustQ(LabelEquals("tier", "backend"))

	tests := []struct {
		name string
		got  Query
		want Query
	}{
		{name: "and all", got: mustAnd(t, All(), envProd), want: envProd},
		{name: "or none", got: mustOr(t, None(), envProd), want: envProd},
		{name: "and none", got: mustAnd(t, None(), envProd), want: None()},
		{name: "or all", got: mustOr(t, All(), envProd), want: All()},
		{name: "double not", got: mustNot(t, mustNot(t, envProd)), want: envProd},
		{
			name: "nested and flattened sorted deduped",
			got:  mustAnd(t, tierBackend, mustAnd(t, envProd, tierBackend)),
			want: mustAnd(t, envProd, tierBackend),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustPredicate(t, tt.got).Query().expr.canonicalKey()
			want := mustPredicate(t, tt.want).Query().expr
			wantKey := "all"
			if want != nil {
				wantKey = want.canonicalKey()
			}
			if got != wantKey {
				t.Fatalf("canonical = %s; want %s", got, wantKey)
			}
		})
	}
}

// TestTermQueryBuildsCanonicalLeaf verifies termQuery stores a leaf key.
func TestTermQueryBuildsCanonicalLeaf(t *testing.T) {
	query := termQuery(term{kind: termName, name: "worker-1", operator: OperatorEquals})
	if query.expr == nil || query.expr.key == "" {
		t.Fatal("termQuery did not build a canonical leaf")
	}
}

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

// TestPredicateMatchBooleanSemantics verifies AND, OR, NOT, All, and None
// evaluate deterministically over one item.
func TestPredicateMatchBooleanSemantics(t *testing.T) {
	item := testItems()[0]

	if !mustPredicate(t, All()).Match(item) {
		t.Fatal("All predicate did not match")
	}
	if mustPredicate(t, None()).Match(item) {
		t.Fatal("None predicate matched")
	}
	and := mustPredicate(t, mustAnd(t, mustQ(LabelEquals("env", "prod")), mustQ(AnnotationEquals("team", "core"))))
	if !and.Match(item) {
		t.Fatal("AND predicate did not match")
	}
	or := mustPredicate(t, mustOr(t, mustQ(LabelEquals("env", "qa")), mustQ(LabelEquals("env", "prod"))))
	if !or.Match(item) {
		t.Fatal("OR predicate did not match")
	}
	not := mustPredicate(t, mustNot(t, mustQ(LabelEquals("env", "qa"))))
	if !not.Match(item) {
		t.Fatal("NOT predicate did not match")
	}
}

// TestMatchTermRejectsUnknownTermKind verifies malformed private terms fail
// closed at evaluation time.
func TestMatchTermRejectsUnknownTermKind(t *testing.T) {
	if matchTerm(term{}, testItems()[0]) {
		t.Fatal("unknown term matched; want false")
	}
}

// TestMatchExprRejectsUnknownExpressionKind verifies malformed private
// expressions never broaden a predicate.
func TestMatchExprRejectsUnknownExpressionKind(t *testing.T) {
	if matchExpr(&expr{}, testItems()[0]) {
		t.Fatal("unknown expression matched; want false")
	}
}

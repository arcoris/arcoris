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

// TestCompileBuildsCanonicalPredicateAndPlan verifies the public compiler
// performs both validation and planning in one pass.
func TestCompileBuildsCanonicalPredicateAndPlan(t *testing.T) {
	query := mustAnd(t, mustQ(LabelEquals("tier", "backend")), mustQ(LabelEquals("env", "prod")))

	predicate := mustPredicate(t, query)

	if predicate.Query().expr.kind != exprAnd {
		t.Fatalf("canonical query kind = %v; want AND", predicate.Query().expr.kind)
	}
	var constraints []Constraint
	predicate.RangeConstraints(func(constraint Constraint) bool {
		constraints = append(constraints, constraint)
		return true
	})
	if len(constraints) != 2 {
		t.Fatalf("constraints = %d; want 2", len(constraints))
	}
}

// TestCompileWrapsExpressionFailures keeps the broad query sentinel stable
// even when the lower failure is an expression-shape problem.
func TestCompileWrapsExpressionFailures(t *testing.T) {
	_, err := Compile(Query{expr: &expr{kind: exprKind(99)}})

	requireErrorIs(t, err, ErrInvalidQuery)
	requireErrorIs(t, err, ErrInvalidExpression)
}

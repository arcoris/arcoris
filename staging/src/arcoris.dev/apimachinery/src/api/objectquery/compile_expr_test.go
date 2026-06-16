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

// TestValidateExprRejectsMalformedNot verifies validation catches private AST
// shapes that public constructors never produce.
func TestValidateExprRejectsMalformedNot(t *testing.T) {
	_, err := validateExpr(&expr{kind: exprNot}, compileOptions{})

	requireErrorIs(t, err, ErrInvalidExpression)
}

// TestValidateCompositeExprCanonicalizesChildren verifies compiler validation
// reuses the same deterministic boolean normalization as public constructors.
func TestValidateCompositeExprCanonicalizesChildren(t *testing.T) {
	left := mustQ(LabelEquals("tier", "backend")).expr
	right := mustQ(LabelEquals("env", "prod")).expr
	raw := &expr{kind: exprAnd, children: []*expr{left, right, left}}

	got, err := validateExpr(raw, compileOptions{})
	requireNoError(t, err)

	want := mustAnd(t, queryFromExpr(left), queryFromExpr(right)).expr
	if got.canonicalKey() != want.canonicalKey() {
		t.Fatalf("canonical key = %q; want %q", got.canonicalKey(), want.canonicalKey())
	}
}

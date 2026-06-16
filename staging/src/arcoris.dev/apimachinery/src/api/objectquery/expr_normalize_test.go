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

// TestNormalizeExprConstantsAndDoubleNot verifies cheap expression identities
// without introducing CNF/DNF expansion.
func TestNormalizeExprConstantsAndDoubleNot(t *testing.T) {
	if got := normalizeExpr(&expr{kind: exprAll}); got != nil {
		t.Fatalf("normalize All = %#v; want nil", got)
	}
	if got := normalizeExpr(&expr{kind: exprNone}); !isNoneExpr(got) {
		t.Fatalf("normalize None = %#v; want none", got)
	}

	query := mustQ(LabelEquals("env", "prod"))
	double := &expr{kind: exprNot, children: []*expr{mustNot(t, query).expr}}
	if got := normalizeExpr(double); got.canonicalKey() != query.expr.canonicalKey() {
		t.Fatalf("normalize double NOT = %q; want %q", got.canonicalKey(), query.expr.canonicalKey())
	}
}

// TestNormalizeChildrenUsesBooleanRules verifies raw child lists flow through
// the same normalization path as public constructors.
func TestNormalizeChildrenUsesBooleanRules(t *testing.T) {
	child := mustQ(LabelEquals("env", "prod")).expr

	if got := normalizeChildren(exprOr, nil); !isNoneExpr(got) {
		t.Fatalf("normalize empty OR = %#v; want none", got)
	}
	if got := normalizeChildren(exprAnd, []*expr{child, child}); got.canonicalKey() != child.canonicalKey() {
		t.Fatalf("normalize duplicate AND = %q; want %q", got.canonicalKey(), child.canonicalKey())
	}
}

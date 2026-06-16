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

// TestNormalizeBoolFlattensSortsAndDeduplicates verifies AND/OR canonical
// children are deterministic regardless of constructor input order.
func TestNormalizeBoolFlattensSortsAndDeduplicates(t *testing.T) {
	env := mustQ(LabelEquals("env", "prod"))
	tier := mustQ(LabelEquals("tier", "backend"))
	nested := mustAnd(t, tier, env, tier)

	got := normalizeBool(exprAnd, []Query{nested, env})

	if got.kind != exprAnd {
		t.Fatalf("kind = %v; want AND", got.kind)
	}
	if len(got.children) != 2 {
		t.Fatalf("children = %d; want 2", len(got.children))
	}
	if got.children[0].canonicalKey() > got.children[1].canonicalKey() {
		t.Fatalf("children are not sorted: %q > %q", got.children[0].canonicalKey(), got.children[1].canonicalKey())
	}
}

// TestCanonicalBoolEmptyAndSingletonRules locks down boolean identity and
// singleton reductions.
func TestCanonicalBoolEmptyAndSingletonRules(t *testing.T) {
	child := mustQ(LabelEquals("env", "prod")).expr

	if got := canonicalBool(exprAnd, nil); got != nil {
		t.Fatalf("empty AND = %#v; want nil All", got)
	}
	if got := canonicalBool(exprOr, nil); !isNoneExpr(got) {
		t.Fatalf("empty OR = %#v; want None", got)
	}
	if got := canonicalBool(exprAnd, []*expr{child}); got.canonicalKey() != child.canonicalKey() {
		t.Fatalf("singleton AND = %q; want %q", got.canonicalKey(), child.canonicalKey())
	}
}

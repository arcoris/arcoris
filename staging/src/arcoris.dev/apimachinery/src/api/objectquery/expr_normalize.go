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

// normalizeExpr reduces cheap boolean identities without expanding into CNF or
// DNF. That keeps canonicalization deterministic and bounded.
func normalizeExpr(e *expr) *expr {
	if e == nil {
		return nil
	}

	switch e.kind {
	case exprAll:
		return nil
	case exprNone:
		return &expr{kind: exprNone, key: "none"}
	case exprAnd:
		return normalizeChildren(exprAnd, e.children)
	case exprOr:
		return normalizeChildren(exprOr, e.children)
	case exprNot:
		q, _ := Not(queryFromExpr(e.children[0]))
		return q.expr
	default:
		out := e.clone()
		out.key = canonicalExprKey(out)
		return out
	}
}

// normalizeChildren funnels raw child nodes back through Query constructors so
// AND and OR use one central normalization path.
func normalizeChildren(kind exprKind, children []*expr) *expr {
	parts := make([]Query, 0, len(children))
	for _, child := range children {
		parts = append(parts, queryFromExpr(child))
	}

	return normalizeBool(kind, parts)
}

// isAllExpr recognizes the package-wide nil-as-All representation.
func isAllExpr(e *expr) bool {
	return e == nil || e.kind == exprAll
}

// isNoneExpr recognizes the explicit false expression.
func isNoneExpr(e *expr) bool {
	return e != nil && e.kind == exprNone
}

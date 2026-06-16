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

// And returns the logical conjunction of parts.
func And(parts ...Query) (Query, error) {
	return queryFromExpr(normalizeBool(exprAnd, parts)), nil
}

// Or returns the logical disjunction of parts.
func Or(parts ...Query) (Query, error) {
	return queryFromExpr(normalizeBool(exprOr, parts)), nil
}

// Not returns the logical negation of query.
func Not(query Query) (Query, error) {
	e := normalizeExpr(query.expr)
	switch {
	case isAllExpr(e):
		return None(), nil
	case isNoneExpr(e):
		return All(), nil
	case e.kind == exprNot:
		return queryFromExpr(e.children[0]), nil
	default:
		out := &expr{kind: exprNot, children: []*expr{e}}
		out.key = canonicalExprKey(out)
		return Query{expr: out}, nil
	}
}

// termQuery wraps a constructor-produced term as an expression leaf and records
// its canonical key immediately.
func termQuery(t term) Query {
	e := &expr{kind: exprTerm, term: t}
	e.key = canonicalExprKey(e)
	return Query{expr: e}
}

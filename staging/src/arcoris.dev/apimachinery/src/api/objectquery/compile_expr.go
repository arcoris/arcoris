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

// validateExpr enforces expression shape and recursively validates leaf terms.
// It returns the canonical expression that should be stored by Predicate.
func validateExpr(e *expr, opts compileOptions) (*expr, error) {
	if e == nil || e.kind == exprAll {
		return nil, nil
	}

	switch e.kind {
	case exprNone:
		return &expr{kind: exprNone, key: "none"}, nil
	case exprTerm:
		return validateTermExpr(e, opts)
	case exprNot:
		return validateNotExpr(e, opts)
	case exprAnd, exprOr:
		return validateCompositeExpr(e, opts)
	default:
		return nil, invalidExpressionError("unknown expression kind")
	}
}

// validateTermExpr validates one leaf expression and refreshes its canonical
// key after the term has been normalized.
func validateTermExpr(e *expr, opts compileOptions) (*expr, error) {
	term, err := validateTerm(e.term, opts)
	if err != nil {
		return nil, err
	}
	out := &expr{kind: exprTerm, term: term}
	out.key = canonicalExprKey(out)
	return out, nil
}

// validateNotExpr preserves NOT as residual boolean structure. It only
// validates the child and applies cheap double-negation/constant normalization.
func validateNotExpr(e *expr, opts compileOptions) (*expr, error) {
	if len(e.children) != 1 {
		return nil, invalidExpressionError("not expression must have one child")
	}
	child, err := validateExpr(e.children[0], opts)
	if err != nil {
		return nil, err
	}
	q, err := Not(queryFromExpr(child))
	if err != nil {
		return nil, err
	}

	return q.expr, nil
}

// validateCompositeExpr validates every child and then rebuilds AND/OR through
// the public constructors so flattening, ordering, and deduplication stay
// centralized.
func validateCompositeExpr(e *expr, opts compileOptions) (*expr, error) {
	if len(e.children) == 0 {
		return normalizeExpr(e), nil
	}

	parts := make([]Query, 0, len(e.children))
	for _, child := range e.children {
		compiled, err := validateExpr(child, opts)
		if err != nil {
			return nil, err
		}
		parts = append(parts, queryFromExpr(compiled))
	}
	if e.kind == exprAnd {
		out, _ := And(parts...)
		return out.expr, nil
	}
	out, _ := Or(parts...)
	return out.expr, nil
}

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

import "sort"

// normalizeBool flattens matching boolean nodes, removes identity operands,
// collapses duplicates, and preserves annihilators such as AND(None, X).
func normalizeBool(kind exprKind, parts []Query) *expr {
	children := make([]*expr, 0, len(parts))
	seen := map[string]struct{}{}

	for _, part := range parts {
		child := normalizeExpr(part.expr)
		if kind == exprAnd && isNoneExpr(child) {
			return &expr{kind: exprNone, key: "none"}
		}
		if kind == exprOr && isAllExpr(child) {
			return nil
		}
		if isAllExpr(child) && kind == exprAnd {
			continue
		}
		if isNoneExpr(child) && kind == exprOr {
			continue
		}
		if child.kind == kind {
			appendCanonicalChildren(&children, seen, child.children)
			continue
		}
		appendCanonicalChild(&children, seen, child)
	}

	return canonicalBool(kind, children)
}

// appendCanonicalChildren appends a normalized child list while sharing the
// caller's duplicate detector.
func appendCanonicalChildren(children *[]*expr, seen map[string]struct{}, nested []*expr) {
	for _, child := range nested {
		appendCanonicalChild(children, seen, child)
	}
}

// appendCanonicalChild records one child by canonical key and clones it before
// it becomes part of the normalized expression.
func appendCanonicalChild(children *[]*expr, seen map[string]struct{}, child *expr) {
	key := child.canonicalKey()
	if _, exists := seen[key]; exists {
		return
	}
	seen[key] = struct{}{}
	*children = append(*children, child.clone())
}

// canonicalBool finalizes a normalized boolean node, including empty and
// singleton reductions.
func canonicalBool(kind exprKind, children []*expr) *expr {
	if len(children) == 0 {
		if kind == exprOr {
			return &expr{kind: exprNone, key: "none"}
		}
		return nil
	}
	if len(children) == 1 {
		return children[0].clone()
	}

	sort.Slice(children, func(i int, j int) bool {
		return children[i].canonicalKey() < children[j].canonicalKey()
	})
	out := &expr{kind: kind, children: children}
	out.key = canonicalExprKey(out)
	return out
}

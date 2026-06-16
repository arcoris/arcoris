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

// exprKind is the private boolean expression node discriminator. It remains
// internal so parsers and adapters cannot depend on AST storage details.
type exprKind uint8

// Private expression node kinds.
const (
	// exprAll is normalized to nil in public Query and Predicate values.
	exprAll exprKind = iota
	// exprNone is the explicit false expression.
	exprNone
	// exprAnd combines canonical, sorted, deduplicated conjunctive children.
	exprAnd
	// exprOr combines canonical, sorted, deduplicated disjunctive children.
	exprOr
	// exprNot stores exactly one child and is never expanded into CNF/DNF.
	exprNot
	// exprTerm wraps one leaf term over an objectstore.ListItem.
	exprTerm
)

// expr is the private immutable expression tree used by Query and Predicate.
// The compiler validates this shape before it is used for evaluation.
type expr struct {
	// kind determines which payload fields are meaningful.
	kind exprKind
	// term is populated only for exprTerm.
	term term
	// children contains boolean operands for AND, OR, and NOT nodes.
	children []*expr
	// key caches the canonical structural key for sorting and deduplication.
	key string
}

// clone returns a detached copy of e so public query values cannot mutate
// compiled predicate state.
func (e *expr) clone() *expr {
	if e == nil {
		return nil
	}

	out := &expr{kind: e.kind, term: e.term.clone(), key: e.key}
	if len(e.children) > 0 {
		out.children = make([]*expr, len(e.children))
		for i, child := range e.children {
			out.children[i] = child.clone()
		}
	}

	return out
}

// canonicalKey returns the stable structural key for e. Nil is encoded as
// "all" because nil is the canonical representation of the All expression.
func (e *expr) canonicalKey() string {
	if e == nil {
		return "all"
	}
	if e.key != "" {
		return e.key
	}

	return canonicalExprKey(e)
}

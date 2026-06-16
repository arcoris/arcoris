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

import "arcoris.dev/apimachinery/api/objectstore"

// Match reports whether item satisfies p.
//
// Match is the semantic source of truth for objectquery. Plans and future
// indexes may narrow candidates, but they must still confirm results here.
func (p Predicate) Match(item objectstore.ListItem) bool {
	return matchExpr(p.expr, item)
}

// matchExpr evaluates the canonical boolean expression tree. Nil is All.
func matchExpr(e *expr, item objectstore.ListItem) bool {
	if e == nil {
		return true
	}

	switch e.kind {
	case exprNone:
		return false
	case exprTerm:
		return matchTerm(e.term, item)
	case exprNot:
		return !matchExpr(e.children[0], item)
	case exprAnd:
		for _, child := range e.children {
			if !matchExpr(child, item) {
				return false
			}
		}
		return true
	case exprOr:
		for _, child := range e.children {
			if matchExpr(child, item) {
				return true
			}
		}
		return false
	default:
		return true
	}
}

// matchTerm evaluates one leaf term against the authoritative storage key and
// committed object state carried by the list item.
func matchTerm(t term, item objectstore.ListItem) bool {
	switch t.kind {
	case termResource:
		return item.Key.Resource == t.resource
	case termNamespace:
		return item.Key.Object.Namespace == t.namespace
	case termName:
		return item.Key.Object.Name == t.name
	case termObject:
		return item.Key.Object.Namespace == t.namespace && item.Key.Object.Name == t.name
	case termKey:
		return item.Key.Equal(t.key)
	case termMetadata:
		return matchMetadataTerm(t, item)
	case termField:
		return matchFieldTerm(t, item)
	default:
		return false
	}
}

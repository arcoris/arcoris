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

// Query is an immutable object query expression.
//
// The zero Query is valid and is equivalent to All().
type Query struct {
	// expr is nil for the canonical All query. Non-nil expressions are treated
	// as immutable; constructors and Predicate.Query clone before exposing them.
	expr *expr
}

// All returns a query that matches every list item.
func All() Query {
	return Query{}
}

// None returns a query that matches no list item.
func None() Query {
	return Query{expr: &expr{kind: exprNone}}
}

// IsZero reports whether q is the zero All query.
func (q Query) IsZero() bool {
	return q.expr == nil
}

// Validate checks whether q can be compiled with default options.
func (q Query) Validate() error {
	_, err := Compile(q)
	return err
}

// queryFromExpr converts a private expression node back into the public value
// type while preserving the nil-as-All convention.
func queryFromExpr(e *expr) Query {
	if e == nil || e.kind == exprAll {
		return Query{}
	}

	return Query{expr: e.clone()}
}

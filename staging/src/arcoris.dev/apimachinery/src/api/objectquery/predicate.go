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

// Predicate is a compiled, canonical, immutable object query evaluator.
type Predicate struct {
	// expr is the validated canonical expression used by Match and Filter.
	expr *expr
	// plan contains detached narrowing hints derived from expr.
	plan Plan
}

// IsZero reports whether p is the zero All predicate.
func (p Predicate) IsZero() bool {
	return p.expr == nil
}

// Query returns a detached canonical query.
func (p Predicate) Query() Query {
	return queryFromExpr(p.expr)
}

// Plan returns detached planning hints for p.
func (p Predicate) Plan() Plan {
	return p.plan.clone()
}

// RangeConstraints visits detached planning constraints until fn returns false.
func (p Predicate) RangeConstraints(fn func(Constraint) bool) {
	p.plan.RangeConstraints(fn)
}

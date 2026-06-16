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

// buildPlan extracts conservative positive constraints from a canonical
// expression. It intentionally ignores OR, NOT, and negative terms unless a
// safe narrowing representation exists.
func buildPlan(e *expr) Plan {
	if e == nil {
		return Plan{}
	}
	switch e.kind {
	case exprTerm:
		return Plan{constraints: constraintsForTerm(e.term)}
	case exprAnd:
		var constraints []Constraint
		for _, child := range e.children {
			constraints = append(constraints, buildPlan(child).constraints...)
		}
		return Plan{constraints: constraints}
	default:
		return Plan{}
	}
}

// constraintsForTerm maps one positive leaf term into zero or more advisory
// constraints. Returning nil means the term remains residual-only.
func constraintsForTerm(t term) []Constraint {
	switch t.kind {
	case termResource:
		return []Constraint{{
			Kind: ConstraintResource,
			Ref:  ConstraintRef{Resource: t.resource},
			Op:   OperatorEquals,
		}}
	case termNamespace:
		return []Constraint{{
			Kind: ConstraintObjectNamespace,
			Ref:  ConstraintRef{Namespace: t.namespace},
			Op:   OperatorEquals,
		}}
	case termName:
		return []Constraint{{
			Kind: ConstraintObjectName,
			Ref:  ConstraintRef{Name: t.name},
			Op:   OperatorEquals,
		}}
	case termObject:
		return []Constraint{{
			Kind: ConstraintObject,
			Ref: ConstraintRef{
				Namespace: t.namespace,
				Name:      t.name,
			},
			Op: OperatorEquals,
		}}
	case termKey:
		return []Constraint{{
			Kind: ConstraintKey,
			Ref:  ConstraintRef{Key: t.key},
			Op:   OperatorEquals,
		}}
	case termMetadata:
		return constraintsForMetadataTerm(t)
	case termField:
		return constraintsForFieldTerm(t)
	default:
		return nil
	}
}

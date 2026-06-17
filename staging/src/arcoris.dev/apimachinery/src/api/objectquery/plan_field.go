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

// constraintsForFieldTerm emits only positive field constraints. Negative
// field operators stay residual because missing-value semantics matter.
func constraintsForFieldTerm(t term) []Constraint {
	if !fieldIndexSupportsOperator(t.field.Index, t.operator) {
		return nil
	}

	switch t.operator {
	case OperatorExists,
		OperatorEquals,
		OperatorIn,
		OperatorLessThan,
		OperatorLessOrEqual,
		OperatorGreaterThan,
		OperatorGreaterOrEqual:
		return []Constraint{{
			Kind:   ConstraintField,
			Ref:    ConstraintRef{Field: t.fieldRef},
			Op:     t.operator,
			Values: cloneValues(t.values),
		}}
	default:
		return nil
	}
}

// fieldIndexSupportsOperator maps SelectableField index hints to conservative
// field constraints. Predicate.Match remains the semantic source of truth.
func fieldIndexSupportsOperator(hint IndexHint, op Operator) bool {
	switch hint {
	case IndexEquality:
		return op == OperatorExists || op == OperatorEquals || op == OperatorIn
	case IndexRange:
		switch op {
		case OperatorExists,
			OperatorEquals,
			OperatorIn,
			OperatorLessThan,
			OperatorLessOrEqual,
			OperatorGreaterThan,
			OperatorGreaterOrEqual:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

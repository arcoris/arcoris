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

import (
	"fmt"

	"arcoris.dev/apimachinery/api/value"
)

// validateFieldLiterals applies field-kind and operator-specific literal
// rules after generic arity checks have passed.
func validateFieldLiterals(field SelectableField, op Operator, literals []value.Value) error {
	if err := validateFieldArity(op, literals); err != nil {
		return err
	}

	switch op {
	case OperatorExists, OperatorDoesNotExist:
		return nil
	case OperatorHasPrefix, OperatorHasSuffix, OperatorContains:
		if field.Kind != value.KindString {
			return unsupportedOperatorError(op, "non-string field "+field.Ref.String())
		}
	case OperatorLessThan,
		OperatorLessOrEqual,
		OperatorGreaterThan,
		OperatorGreaterOrEqual:
		if !isOrderedKind(field.Kind) {
			return unsupportedOperatorError(op, "non-ordered field "+field.Ref.String())
		}
	}

	for _, literal := range literals {
		if literal.IsZero() {
			return invalidFieldError(ErrInvalidTerm, "zero field literal")
		}
		if literal.IsNull() {
			if op == OperatorEquals || op == OperatorNotEquals ||
				op == OperatorIn || op == OperatorNotIn {
				continue
			}
			return invalidFieldError(
				ErrUnsupportedOperator,
				"operator %s cannot use null literal",
				op.String(),
			)
		}
		if literal.Kind() != field.Kind {
			return invalidFieldError(
				ErrInvalidTerm,
				"literal kind %s does not match field kind %s",
				literal.Kind().String(),
				field.Kind.String(),
			)
		}
	}

	return nil
}

// isOrderedKind identifies value kinds with stable less/greater semantics.
func isOrderedKind(kind value.Kind) bool {
	return kind.IsNumber() || kind.IsTemporal()
}

// requireFieldMatchComparable rejects runtime comparisons across different
// value kinds. Compile-time validation should usually make this a guardrail.
func requireFieldMatchComparable(actual value.Value, literal value.Value) error {
	if actual.Kind() != literal.Kind() {
		return fmt.Errorf("field kind %s does not match literal kind %s",
			actual.Kind().String(),
			literal.Kind().String(),
		)
	}

	return nil
}

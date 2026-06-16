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

import "arcoris.dev/apimachinery/api/value"

// fieldQuery validates constructor-level shape and stores canonical literal
// copies. Registry-dependent validation happens later in Compile.
func fieldQuery(ref FieldRef, op Operator, literals ...value.Value) (Query, error) {
	if err := ref.Validate(); err != nil {
		return Query{}, invalidFieldError(err, "invalid field ref %s", ref.String())
	}
	if err := validateFieldArity(op, literals); err != nil {
		return Query{}, err
	}

	canonical, err := canonicalValues(literals)
	if err != nil {
		return Query{}, invalidFieldError(err, "invalid field literal")
	}

	return termQuery(term{
		kind:     termField,
		fieldRef: ref,
		values:   canonical,
		operator: op,
	}), nil
}

// validateFieldArity enforces operator arity before field-specific type checks.
func validateFieldArity(op Operator, literals []value.Value) error {
	if !op.IsValid() {
		return invalidOperatorError(op)
	}

	switch op {
	case OperatorExists, OperatorDoesNotExist:
		if len(literals) != 0 {
			return invalidTermError("%s takes no literals", op.String())
		}
	case OperatorIn, OperatorNotIn:
		if len(literals) == 0 {
			return invalidTermError("%s requires at least one literal", op.String())
		}
	default:
		if len(literals) != 1 {
			return invalidTermError("%s takes exactly one literal", op.String())
		}
	}

	return nil
}

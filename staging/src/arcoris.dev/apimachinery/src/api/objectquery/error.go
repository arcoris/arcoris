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
	"errors"
	"fmt"
)

// Package sentinel errors are intentionally broad enough for stable errors.Is
// checks while preserving lower causes through errors.Join.
var (
	// ErrInvalidQuery reports a malformed object query.
	ErrInvalidQuery = errors.New("invalid object query")
	// ErrInvalidExpression reports an invalid boolean query expression.
	ErrInvalidExpression = errors.New("invalid object query expression")
	// ErrInvalidTerm reports an invalid typed query term.
	ErrInvalidTerm = errors.New("invalid object query term")
	// ErrInvalidField reports an invalid selectable field query.
	ErrInvalidField = errors.New("invalid object query field")
	// ErrInvalidOperator reports an unknown query operator.
	ErrInvalidOperator = errors.New("invalid object query operator")
	// ErrUnsupportedOperator reports an operator that is known but not allowed
	// for the selected query domain or field kind.
	ErrUnsupportedOperator = errors.New("unsupported object query operator")
	// ErrUnresolvedField reports a field query whose FieldRef is not registered
	// in the supplied SelectableFieldSet.
	ErrUnresolvedField = errors.New("unresolved object query field")
	// ErrInvalidChange reports an objectstore.Change that cannot be projected.
	ErrInvalidChange = errors.New("invalid object query change")
)

// invalidQueryError preserves the broad query sentinel while keeping the
// lower validation cause visible to errors.Is.
func invalidQueryError(cause error) error {
	if cause == nil {
		return ErrInvalidQuery
	}

	return errors.Join(ErrInvalidQuery, cause)
}

// invalidExpressionError classifies failures in boolean expression shape.
func invalidExpressionError(format string, args ...any) error {
	return errors.Join(ErrInvalidQuery, ErrInvalidExpression, fmt.Errorf(format, args...))
}

// invalidTermError classifies failures in typed query terms.
func invalidTermError(format string, args ...any) error {
	return errors.Join(ErrInvalidQuery, ErrInvalidTerm, fmt.Errorf(format, args...))
}

// invalidFieldError keeps field failures distinguishable from generic term
// failures while preserving any lower sentinel.
func invalidFieldError(cause error, format string, args ...any) error {
	err := errors.Join(ErrInvalidQuery, ErrInvalidField, fmt.Errorf(format, args...))
	if cause != nil {
		err = errors.Join(err, cause)
	}

	return err
}

// invalidOperatorError reports an unknown operator value before domain-specific
// operator support is considered.
func invalidOperatorError(op Operator) error {
	return errors.Join(
		ErrInvalidQuery,
		ErrInvalidOperator,
		fmt.Errorf("operator %s is invalid", op.String()),
	)
}

// unsupportedOperatorError reports a known operator that is forbidden for the
// selected metadata domain or selectable field.
func unsupportedOperatorError(op Operator, domain string) error {
	return errors.Join(
		ErrInvalidQuery,
		ErrUnsupportedOperator,
		fmt.Errorf("operator %s is unsupported for %s", op.String(), domain),
	)
}

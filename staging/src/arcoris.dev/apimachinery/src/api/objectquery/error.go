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

// ErrorReason identifies the stable reason for a query validation failure.
type ErrorReason string

// Stable query diagnostic reasons.
const (
	ErrorReasonInvalidExpression   ErrorReason = "invalid_expression"
	ErrorReasonInvalidTerm         ErrorReason = "invalid_term"
	ErrorReasonInvalidField        ErrorReason = "invalid_field"
	ErrorReasonInvalidOperator     ErrorReason = "invalid_operator"
	ErrorReasonUnsupportedOperator ErrorReason = "unsupported_operator"
	ErrorReasonUnresolvedField     ErrorReason = "unresolved_field"
	ErrorReasonInvalidChange       ErrorReason = "invalid_change"
)

// Error carries a stable objectquery diagnostic while preserving lower causes.
type Error struct {
	// Path identifies the logical compiler location that failed.
	Path string
	// Reason is the stable machine-readable failure class.
	Reason ErrorReason
	// Cause is the underlying validation or construction failure.
	Cause error
}

// Error returns a compact diagnostic string.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Path == "" {
		return string(e.Reason) + ": " + e.Cause.Error()
	}

	return e.Path + ": " + string(e.Reason) + ": " + e.Cause.Error()
}

// Unwrap returns the underlying cause for errors.Is/errors.As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

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
	return newError("query.expression", ErrorReasonInvalidExpression, ErrInvalidExpression, fmt.Errorf(format, args...))
}

// invalidTermError classifies failures in typed query terms.
func invalidTermError(format string, args ...any) error {
	return newError("query.term", ErrorReasonInvalidTerm, ErrInvalidTerm, fmt.Errorf(format, args...))
}

// invalidFieldError keeps field failures distinguishable from generic term
// failures while preserving any lower sentinel.
func invalidFieldError(cause error, format string, args ...any) error {
	return newError("query.field", ErrorReasonInvalidField, ErrInvalidField, fmt.Errorf(format, args...), cause)
}

// invalidOperatorError reports an unknown operator value before domain-specific
// operator support is considered.
func invalidOperatorError(op Operator) error {
	return newError(
		"query.operator",
		ErrorReasonInvalidOperator,
		ErrInvalidOperator,
		fmt.Errorf("operator %s is invalid", op.String()),
	)
}

// unsupportedOperatorError reports a known operator that is forbidden for the
// selected metadata domain or selectable field.
func unsupportedOperatorError(op Operator, domain string) error {
	return newError(
		"query.operator",
		ErrorReasonUnsupportedOperator,
		ErrUnsupportedOperator,
		fmt.Errorf("operator %s is unsupported for %s", op.String(), domain),
	)
}

// unresolvedFieldError reports a missing selectable field declaration.
func unresolvedFieldError(ref FieldRef, cause error) error {
	return newError(
		"query.field."+ref.String(),
		ErrorReasonUnresolvedField,
		ErrUnresolvedField,
		cause,
	)
}

// invalidChangeError reports an objectstore.Change projection failure.
func invalidChangeError(cause error) error {
	return newError("query.change", ErrorReasonInvalidChange, ErrInvalidChange, cause)
}

// newError joins broad sentinels with a structured diagnostic.
func newError(path string, reason ErrorReason, sentinel error, causes ...error) error {
	joined := []error{ErrInvalidQuery, sentinel}
	var cause error
	for _, item := range causes {
		if item == nil {
			continue
		}
		if cause == nil {
			cause = item
		} else {
			cause = errors.Join(cause, item)
		}
		joined = append(joined, item)
	}
	if cause == nil {
		cause = sentinel
	}

	joined = append(joined, &Error{Path: path, Reason: reason, Cause: cause})
	return errors.Join(joined...)
}

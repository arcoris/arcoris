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

package identity

import (
	"errors"
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
	"arcoris.dev/apimachinery/api/meta/internal/metagrammar"
)

// Metadata identity sentinels classify broad validation and encoding failures.
var (
	// ErrInvalidName classifies malformed object names.
	ErrInvalidName = errors.New("invalid metadata name")
	// ErrInvalidNamePrefix classifies malformed metadata name prefixes.
	ErrInvalidNamePrefix = errors.New("invalid metadata name prefix")
	// ErrInvalidNamespace classifies malformed namespaces.
	ErrInvalidNamespace = errors.New("invalid metadata namespace")
	// ErrInvalidUID classifies malformed object UIDs.
	ErrInvalidUID = errors.New("invalid metadata UID")
	// ErrInvalidObjectName classifies malformed object name composites.
	ErrInvalidObjectName = errors.New("invalid metadata object name")
	// ErrInvalidObjectIdentity classifies malformed object identity composites.
	ErrInvalidObjectIdentity = errors.New("invalid metadata object identity")
	// ErrInvalidObjectReference classifies malformed object references.
	ErrInvalidObjectReference = errors.New("invalid metadata object reference")
	// ErrInvalidObjectIdentityReference classifies malformed UID-pinned object references.
	ErrInvalidObjectIdentityReference = errors.New("invalid metadata object identity reference")
	// ErrInvalidJSON classifies malformed JSON scalar metadata identity values.
	ErrInvalidJSON = errors.New("invalid metadata identity JSON")
	// ErrNilReceiver classifies Unmarshal calls made on nil metadata identity pointers.
	ErrNilReceiver = errors.New("nil metadata identity receiver")
)

// ErrorReason identifies a precise metadata identity validation failure.
type ErrorReason string

// Metadata identity reasons refine broad sentinel errors with stable diagnostics.
const (
	// ErrorReasonEmptyValue reports a required value that is absent.
	ErrorReasonEmptyValue ErrorReason = "empty_value"
	// ErrorReasonInvalidLength reports a length limit violation.
	ErrorReasonInvalidLength ErrorReason = "invalid_length"
	// ErrorReasonInvalidCharacter reports a byte outside the allowed grammar.
	ErrorReasonInvalidCharacter ErrorReason = "invalid_character"
	// ErrorReasonInvalidEdge reports an invalid first or last byte.
	ErrorReasonInvalidEdge ErrorReason = "invalid_edge"
	// ErrorReasonInvalidForm reports a malformed composite identity.
	ErrorReasonInvalidForm ErrorReason = "invalid_form"
	// ErrorReasonInvalidJSON reports non-string or malformed JSON.
	ErrorReasonInvalidJSON ErrorReason = "invalid_json"
	// ErrorReasonNilReceiver reports an Unmarshal call on a nil pointer.
	ErrorReasonNilReceiver ErrorReason = "nil_receiver"
)

// Error is the structured diagnostic returned by metadata identity validation.
type Error struct {
	// Value stores the rejected scalar or composite text when one is available.
	Value string

	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a human-readable metadata identity diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	record := e.Record
	if e.Value != "" && record.Detail != "" {
		record.Detail += fmt.Sprintf(" (value %q)", e.Value)
	}

	return record.Format("meta/identity")
}

// Unwrap preserves broad and nested error identity.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

// invalid reports a metadata identity validation failure.
func invalid(path string, value string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Value:  value,
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}

// invalidf reports a metadata identity validation failure with formatted detail.
func invalidf(path string, value string, err error, reason ErrorReason, format string, args ...any) error {
	return invalid(path, value, err, reason, fmt.Sprintf(format, args...))
}

// nested reports a failure caused by a nested metadata or API identity value.
func nested(path string, value string, err error, cause error) error {
	return &Error{
		Value: value,
		Record: diagnostic.WrapRecord(
			path,
			err,
			ErrorReasonInvalidForm,
			fmt.Sprintf("nested value is invalid: %v", cause),
			cause,
		),
	}
}

// fromGrammar converts internal grammar diagnostics to public metadata identity
// diagnostics without exposing metagrammar as API surface.
func fromGrammar(path string, value string, err error, v *metagrammar.Violation) error {
	if v == nil {
		return nil
	}

	return invalid(path, value, err, reasonFromGrammar(v.Reason), v.Detail)
}

// reasonFromGrammar maps internal grammar reasons to public reason values.
func reasonFromGrammar(reason metagrammar.Reason) ErrorReason {
	switch reason {
	case metagrammar.ReasonEmptyValue:
		return ErrorReasonEmptyValue
	case metagrammar.ReasonInvalidLength:
		return ErrorReasonInvalidLength
	case metagrammar.ReasonInvalidCharacter:
		return ErrorReasonInvalidCharacter
	case metagrammar.ReasonInvalidEdge:
		return ErrorReasonInvalidEdge
	default:
		return ErrorReasonInvalidForm
	}
}

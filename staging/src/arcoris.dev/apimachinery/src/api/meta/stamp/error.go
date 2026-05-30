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

package stamp

import (
	"errors"
	"fmt"
	"strings"

	"arcoris.dev/apimachinery/api/meta/internal/metagrammar"
)

// Stamp sentinels classify broad validation and encoding failures.
var (
	// ErrInvalidResourceVersion classifies malformed resource version tokens.
	ErrInvalidResourceVersion = errors.New("invalid metadata resource version")
	// ErrInvalidGeneration classifies malformed generation values.
	ErrInvalidGeneration = errors.New("invalid metadata generation")
	// ErrInvalidTimestamp classifies malformed metadata timestamps.
	ErrInvalidTimestamp = errors.New("invalid metadata timestamp")
	// ErrInvalidDeletion classifies malformed deletion metadata.
	ErrInvalidDeletion = errors.New("invalid metadata deletion")
	// ErrInvalidJSON classifies malformed JSON scalar stamp values.
	ErrInvalidJSON = errors.New("invalid metadata stamp JSON")
	// ErrNilReceiver classifies Unmarshal calls made on nil stamp pointers.
	ErrNilReceiver = errors.New("nil metadata stamp receiver")
)

// ErrorReason identifies a precise stamp validation failure.
type ErrorReason string

// Stamp reasons refine broad sentinel errors with stable diagnostics.
const (
	// ErrorReasonEmptyValue reports a required value that is absent.
	ErrorReasonEmptyValue ErrorReason = "empty_value"
	// ErrorReasonInvalidLength reports an opaque token length violation.
	ErrorReasonInvalidLength ErrorReason = "invalid_length"
	// ErrorReasonInvalidCharacter reports an unsafe token byte.
	ErrorReasonInvalidCharacter ErrorReason = "invalid_character"
	// ErrorReasonInvalidForm reports malformed scalar syntax.
	ErrorReasonInvalidForm ErrorReason = "invalid_form"
	// ErrorReasonInvalidJSON reports non-string or malformed JSON.
	ErrorReasonInvalidJSON ErrorReason = "invalid_json"
	// ErrorReasonNilReceiver reports an Unmarshal call on a nil pointer.
	ErrorReasonNilReceiver ErrorReason = "nil_receiver"
)

// Error is the structured diagnostic returned by metadata stamp validation.
type Error struct {
	// Path identifies the stamp field that failed validation.
	Path string
	// Err is the broad sentinel used with errors.Is.
	Err error
	// Reason gives stable machine-readable detail within Err.
	Reason ErrorReason
	// Detail gives human-readable context for logs and diagnostics.
	Detail string
	// Cause preserves nested validation or decoding failures.
	Cause error
}

// Error returns a human-readable stamp diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"meta/stamp"}

	if e.Path != "" {
		parts = append(parts, e.Path)
	}

	if e.Err != nil {
		parts = append(parts, e.Err.Error())
	}

	if e.Reason != "" {
		parts = append(parts, string(e.Reason))
	}

	if e.Detail != "" {
		parts = append(parts, e.Detail)
	}

	return strings.Join(parts, ": ")
}

// Unwrap preserves broad and nested error identity.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	if e.Err != nil && e.Cause != nil {
		return errors.Join(e.Err, e.Cause)
	}

	if e.Err != nil {
		return e.Err
	}

	return e.Cause
}

// invalid builds a direct stamp validation diagnostic.
func invalid(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: reason,
		Detail: detail,
	}
}

// invalidf builds a direct stamp validation diagnostic with formatted detail.
func invalidf(path string, err error, reason ErrorReason, format string, args ...any) error {
	return invalid(path, err, reason, fmt.Sprintf(format, args...))
}

// nested wraps a failure reported by a nested stamp value.
func nested(path string, err error, cause error) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: ErrorReasonInvalidForm,
		Detail: fmt.Sprintf("nested value is invalid: %v", cause),
		Cause:  cause,
	}
}

// fromGrammar maps internal grammar failures to stamp diagnostics.
func fromGrammar(path string, err error, v *metagrammar.Violation) error {
	if v == nil {
		return nil
	}

	return invalid(path, err, reasonFromGrammar(v.Reason), v.Detail)
}

// reasonFromGrammar maps internal grammar reasons to stamp reasons.
func reasonFromGrammar(reason metagrammar.Reason) ErrorReason {
	switch reason {
	case metagrammar.ReasonEmptyValue:
		return ErrorReasonEmptyValue
	case metagrammar.ReasonInvalidLength:
		return ErrorReasonInvalidLength
	case metagrammar.ReasonInvalidCharacter:
		return ErrorReasonInvalidCharacter
	default:
		return ErrorReasonInvalidForm
	}
}

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

package annotations

import (
	"errors"
	"fmt"
	"strings"

	"arcoris.dev/apimachinery/api/meta/internal/metagrammar"
)

// Annotation sentinels classify broad validation and encoding failures.
var (
	// ErrInvalidKey classifies malformed annotation keys.
	ErrInvalidKey = errors.New("invalid metadata annotation key")
	// ErrInvalidValue classifies malformed annotation values.
	ErrInvalidValue = errors.New("invalid metadata annotation value")
	// ErrInvalidSet classifies malformed annotation maps.
	ErrInvalidSet = errors.New("invalid metadata annotation set")
	// ErrInvalidJSON classifies malformed JSON annotation scalar values.
	ErrInvalidJSON = errors.New("invalid metadata annotation JSON")
	// ErrNilReceiver classifies Unmarshal calls made on nil annotation pointers.
	ErrNilReceiver = errors.New("nil metadata annotation receiver")
)

// ErrorReason identifies a precise annotation validation failure.
type ErrorReason string

// Annotation reasons refine broad sentinel errors with stable diagnostics.
const (
	// ErrorReasonEmptyValue reports a required annotation value that is absent.
	ErrorReasonEmptyValue ErrorReason = "empty_value"
	// ErrorReasonInvalidLength reports a length limit violation.
	ErrorReasonInvalidLength ErrorReason = "invalid_length"
	// ErrorReasonInvalidCharacter reports an unsafe byte in a key or value.
	ErrorReasonInvalidCharacter ErrorReason = "invalid_character"
	// ErrorReasonInvalidEdge reports an invalid leading or trailing byte.
	ErrorReasonInvalidEdge ErrorReason = "invalid_edge"
	// ErrorReasonInvalidForm reports malformed annotation structure.
	ErrorReasonInvalidForm ErrorReason = "invalid_form"
	// ErrorReasonInvalidJSON reports malformed or non-string JSON scalar input.
	ErrorReasonInvalidJSON ErrorReason = "invalid_json"
	// ErrorReasonNilReceiver reports decoding into a nil pointer receiver.
	ErrorReasonNilReceiver ErrorReason = "nil_receiver"
)

// Error is the structured diagnostic returned by annotation validation.
type Error struct {
	// Path identifies the annotation field or map entry that failed validation.
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

// Error returns a compact human-readable annotation diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	parts := []string{"meta/annotations"}
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

// Unwrap preserves both the broad sentinel and nested cause identity.
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

// invalid builds a direct annotation validation diagnostic.
func invalid(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: reason,
		Detail: detail,
	}
}

// nested wraps a failure reported by a nested annotation value.
func nested(path string, err error, cause error) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: ErrorReasonInvalidForm,
		Detail: fmt.Sprintf("nested value is invalid: %v", cause),
		Cause:  cause,
	}
}

// fromGrammar maps internal grammar failures to annotation diagnostics.
func fromGrammar(path string, err error, v *metagrammar.Violation) error {
	if v == nil {
		return nil
	}
	return invalid(path, err, reasonFromGrammar(v.Reason), v.Detail)
}

// reasonFromGrammar maps internal grammar reasons to annotation reasons.
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

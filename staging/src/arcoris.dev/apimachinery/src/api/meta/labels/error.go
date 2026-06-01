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

package labels

import (
	"errors"
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
	"arcoris.dev/apimachinery/api/meta/internal/metagrammar"
)

// Label sentinels classify broad validation and encoding failures.
var (
	// ErrInvalidKey classifies malformed label keys.
	ErrInvalidKey = errors.New("invalid metadata label key")
	// ErrInvalidValue classifies malformed label values.
	ErrInvalidValue = errors.New("invalid metadata label value")
	// ErrInvalidSet classifies malformed label sets.
	ErrInvalidSet = errors.New("invalid metadata label set")
	// ErrInvalidJSON classifies malformed JSON scalar label values.
	ErrInvalidJSON = errors.New("invalid metadata label JSON")
	// ErrNilReceiver classifies Unmarshal calls made on nil label pointers.
	ErrNilReceiver = errors.New("nil metadata label receiver")
)

// ErrorReason identifies a precise label validation failure.
type ErrorReason string

// Label reasons refine broad sentinel errors with stable diagnostics.
const (
	// ErrorReasonEmptyValue reports a required label value that is absent.
	ErrorReasonEmptyValue ErrorReason = "empty_value"
	// ErrorReasonInvalidLength reports a length limit violation.
	ErrorReasonInvalidLength ErrorReason = "invalid_length"
	// ErrorReasonInvalidCharacter reports an unsafe byte in a key or value.
	ErrorReasonInvalidCharacter ErrorReason = "invalid_character"
	// ErrorReasonInvalidEdge reports an invalid leading or trailing byte.
	ErrorReasonInvalidEdge ErrorReason = "invalid_edge"
	// ErrorReasonInvalidForm reports malformed label structure.
	ErrorReasonInvalidForm ErrorReason = "invalid_form"
	// ErrorReasonInvalidJSON reports malformed or non-string JSON scalar input.
	ErrorReasonInvalidJSON ErrorReason = "invalid_json"
	// ErrorReasonNilReceiver reports decoding into a nil pointer receiver.
	ErrorReasonNilReceiver ErrorReason = "nil_receiver"
)

// Error is the structured diagnostic returned by label validation.
type Error struct {
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable label diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("meta/labels")
}

// Unwrap preserves both the broad sentinel and nested cause identity.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

// invalid builds a direct label validation diagnostic.
func invalid(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}

// nested wraps a failure reported by a nested label value.
func nested(path string, err error, cause error) error {
	return &Error{
		Record: diagnostic.WrapRecord(
			path,
			err,
			ErrorReasonInvalidForm,
			fmt.Sprintf("nested value is invalid: %v", cause),
			cause,
		),
	}
}

// fromGrammar maps internal grammar failures to label diagnostics.
func fromGrammar(path string, err error, v *metagrammar.Violation) error {
	if v == nil {
		return nil
	}

	return invalid(path, err, reasonFromGrammar(v.Reason), v.Detail)
}

// reasonFromGrammar maps internal grammar reasons to label reasons.
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

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
)

var (
	// ErrInvalidIdentifier classifies malformed API identity values.
	ErrInvalidIdentifier = errors.New("invalid API identity")

	// ErrInvalidJSON classifies malformed JSON scalar identity encodings.
	ErrInvalidJSON = errors.New("invalid API identity JSON")

	// ErrNilReceiver classifies Unmarshal calls made on nil identity pointers.
	ErrNilReceiver = errors.New("nil API identity receiver")
)

// ErrorReason identifies the precise identity invariant that failed.
//
// Err remains the broad sentinel used with errors.Is. Reason gives callers a
// stable diagnostic without requiring them to parse the human-facing message.
type ErrorReason string

const (
	// ErrorReasonEmptyValue reports a missing required identity segment.
	ErrorReasonEmptyValue ErrorReason = "empty_value"
	// ErrorReasonInvalidForm reports the wrong canonical grammar shape.
	ErrorReasonInvalidForm ErrorReason = "invalid_form"
	// ErrorReasonInvalidLength reports a DNS or identity length violation.
	ErrorReasonInvalidLength ErrorReason = "invalid_length"
	// ErrorReasonInvalidCharacter reports a byte outside the allowed grammar.
	ErrorReasonInvalidCharacter ErrorReason = "invalid_character"
	// ErrorReasonInvalidEdge reports an invalid first or last byte.
	ErrorReasonInvalidEdge ErrorReason = "invalid_edge"
	// ErrorReasonInvalidJSON reports a non-string or malformed JSON scalar.
	ErrorReasonInvalidJSON ErrorReason = "invalid_json"
	// ErrorReasonNilReceiver reports an Unmarshal call on a nil pointer.
	ErrorReasonNilReceiver ErrorReason = "nil_receiver"
)

// Error is the structured diagnostic returned by identity validation.
//
// Name identifies the identity kind that failed, such as "group" or
// "group/version/resource". Value stores the rejected canonical text when one
// exists. Err preserves broad classification through errors.Is. Cause stores a
// nested parser, validation, or JSON error when the failure originated below
// the current identity boundary.
type Error struct {
	Name  string
	Value string

	// Record stores the shared sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a stable human-readable identity diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	subject := e.Name
	switch e.Err {
	case ErrInvalidJSON:
		subject = "invalid " + subject + " JSON"
	case ErrNilReceiver:
		subject = "nil " + subject + " receiver"
	default:
		subject = "invalid " + subject
	}
	if e.Value != "" {
		subject += fmt.Sprintf(" %q", e.Value)
	}

	return diagnostic.NewRecord(subject, nil, e.Reason, e.Detail).Format("identity")
}

// Unwrap preserves broad and nested error identity.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Record.Unwrap()
}

// invalid reports a malformed identity with a precise reason and detail.
func invalid(name, value string, reason ErrorReason, detail string) error {
	return &Error{
		Name:   name,
		Value:  value,
		Record: diagnostic.NewRecord("", ErrInvalidIdentifier, reason, detail),
	}
}

// invalidf reports a malformed identity with formatted detail.
func invalidf(name, value string, reason ErrorReason, format string, args ...any) error {
	return invalid(name, value, reason, fmt.Sprintf(format, args...))
}

// invalidValue wraps an error from a nested identity segment.
func invalidValue(name, value string, err error) error {
	return &Error{
		Name:  name,
		Value: value,
		Record: diagnostic.WrapRecord(
			"",
			ErrInvalidIdentifier,
			reasonOf(err),
			fmt.Sprintf("nested identity is invalid: %v", err),
			err,
		),
	}
}

// invalidJSON reports malformed JSON scalar encoding.
func invalidJSON(name, value string, reason ErrorReason, detail string, cause error) error {
	return &Error{
		Name:   name,
		Value:  value,
		Record: diagnostic.WrapRecord("", ErrInvalidJSON, reason, detail, cause),
	}
}

// nilReceiver reports Unmarshal calls made on nil receiver pointers.
func nilReceiver(name string) error {
	return &Error{
		Name: name,
		Record: diagnostic.NewRecord(
			"",
			ErrNilReceiver,
			ErrorReasonNilReceiver,
			"identity cannot be decoded into a nil receiver",
		),
	}
}

// reasonOf returns the precise reason from nested identity errors.
func reasonOf(err error) ErrorReason {
	var identityErr *Error
	if errors.As(err, &identityErr) {
		return identityErr.Reason
	}
	return ErrorReasonInvalidForm
}

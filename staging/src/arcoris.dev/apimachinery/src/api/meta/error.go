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

package meta

import (
	"errors"
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
	"arcoris.dev/apimachinery/api/meta/internal/metagrammar"
)

// Root metadata sentinels classify broad validation and encoding failures.
var (
	// ErrInvalidTypeMeta classifies malformed apiVersion/kind metadata.
	ErrInvalidTypeMeta = errors.New("invalid API type metadata")
	// ErrInvalidObjectMeta classifies malformed object metadata.
	ErrInvalidObjectMeta = errors.New("invalid API object metadata")
	// ErrInvalidPageMeta classifies malformed page metadata.
	ErrInvalidPageMeta = errors.New("invalid API page metadata")
	// ErrInvalidPageToken classifies malformed opaque pagination tokens.
	ErrInvalidPageToken = errors.New("invalid API page token")
	// ErrInvalidJSON classifies malformed JSON scalar metadata values.
	ErrInvalidJSON = errors.New("invalid API metadata JSON")
	// ErrNilReceiver classifies Unmarshal calls made on nil metadata pointers.
	ErrNilReceiver = errors.New("nil API metadata receiver")
)

// Error is the structured diagnostic returned by root metadata validation.
type Error struct {
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("meta")
}

// Unwrap preserves both the broad sentinel and nested cause identity.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

// invalid builds a direct root metadata validation diagnostic.
func invalid(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}

// nested wraps a failure reported by a nested metadata domain package.
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

// fromGrammar maps internal grammar failures to root metadata diagnostics.
func fromGrammar(path string, err error, v *metagrammar.Violation) error {
	if v == nil {
		return nil
	}

	return invalid(path, err, reasonFromGrammar(v.Reason), v.Detail)
}

// reasonFromGrammar maps internal grammar reasons to root metadata reasons.
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

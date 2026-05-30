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
	"strings"

	"arcoris.dev/apimachinery/api/meta/internal/metagrammar"
)

// Root metadata sentinels classify broad validation and encoding failures.
var (
	// ErrInvalidTypeMeta classifies malformed apiVersion/kind metadata.
	ErrInvalidTypeMeta = errors.New("invalid API type metadata")
	// ErrInvalidObjectMeta classifies malformed object metadata.
	ErrInvalidObjectMeta = errors.New("invalid API object metadata")
	// ErrInvalidListMeta classifies malformed list/page metadata.
	ErrInvalidListMeta = errors.New("invalid API list metadata")
	// ErrInvalidPageToken classifies malformed opaque pagination tokens.
	ErrInvalidPageToken = errors.New("invalid API page token")
	// ErrInvalidJSON classifies malformed JSON scalar metadata values.
	ErrInvalidJSON = errors.New("invalid API metadata JSON")
	// ErrNilReceiver classifies Unmarshal calls made on nil metadata pointers.
	ErrNilReceiver = errors.New("nil API metadata receiver")
)

// Error is the structured diagnostic returned by root metadata validation.
type Error struct {
	// Path identifies the metadata field that failed validation.
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

// Error returns a compact human-readable diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	parts := []string{"meta"}
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

// invalid builds a direct root metadata validation diagnostic.
func invalid(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: reason,
		Detail: detail,
	}
}

// nested wraps a failure reported by a nested metadata domain package.
func nested(path string, err error, cause error) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: ErrorReasonInvalidForm,
		Detail: fmt.Sprintf("nested value is invalid: %v", cause),
		Cause:  cause,
	}
}

// fromGrammar maps internal grammar failures to root metadata diagnostics.
func fromGrammar(path string, err error, v *metagrammar.Violation) error {
	if v == nil {
		return nil
	}
	reason := ErrorReasonInvalidForm
	if v.Reason == metagrammar.ReasonInvalidCharacter {
		reason = ErrorReasonInvalidCharacter
	}
	if v.Reason == metagrammar.ReasonEmptyValue {
		reason = ErrorReasonEmptyValue
	}
	return invalid(path, err, reason, v.Detail)
}

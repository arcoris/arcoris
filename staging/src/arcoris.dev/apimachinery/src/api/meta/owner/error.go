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

package owner

import (
	"errors"
	"fmt"
	"strings"
)

// Owner sentinels classify broad validation failures.
var (
	// ErrInvalidReference classifies malformed owner reference entries.
	ErrInvalidReference = errors.New("invalid metadata owner reference")
	// ErrInvalidList classifies malformed owner reference lists.
	ErrInvalidList = errors.New("invalid metadata owner reference list")
	// ErrMultipleControllers classifies lists with more than one controller owner.
	ErrMultipleControllers = errors.New("multiple metadata controller owners")
	// ErrDuplicateReference classifies repeated owner references.
	ErrDuplicateReference = errors.New("duplicate metadata owner reference")
)

// ErrorReason identifies a precise owner-reference validation failure.
type ErrorReason string

// Owner reasons refine broad sentinel errors with stable diagnostics.
const (
	// ErrorReasonInvalidReference reports a malformed nested object reference.
	ErrorReasonInvalidReference ErrorReason = "invalid_reference"
	// ErrorReasonMultipleControllers reports more than one controlling owner.
	ErrorReasonMultipleControllers ErrorReason = "multiple_controllers"
	// ErrorReasonDuplicateReference reports a repeated owner reference.
	ErrorReasonDuplicateReference ErrorReason = "duplicate_reference"
)

// Error is the structured diagnostic returned by owner validation.
type Error struct {
	// Path identifies the owner field or list entry that failed validation.
	Path string
	// Err is the broad sentinel used with errors.Is.
	Err error
	// Reason gives stable machine-readable detail within Err.
	Reason ErrorReason
	// Detail gives human-readable context for logs and diagnostics.
	Detail string
	// Cause preserves nested validation failures.
	Cause error
}

// Error returns a compact human-readable owner diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"meta/owner"}

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

// invalid builds a direct owner validation diagnostic.
func invalid(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: reason,
		Detail: detail,
	}
}

// nested wraps a failure reported by a nested object reference.
func nested(path string, err error, cause error) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: ErrorReasonInvalidReference,
		Detail: fmt.Sprintf("nested value is invalid: %v", cause),
		Cause:  cause,
	}
}

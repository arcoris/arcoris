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

	"arcoris.dev/apimachinery/api/internal/diagnostic"
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
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable owner diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("meta/owner")
}

// Unwrap preserves both the broad sentinel and nested cause identity.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

// invalid builds a direct owner validation diagnostic.
func invalid(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}

// nested wraps a failure reported by a nested object reference.
func nested(path string, err error, cause error) error {
	return &Error{
		Record: diagnostic.WrapRecord(
			path,
			err,
			ErrorReasonInvalidReference,
			fmt.Sprintf("nested value is invalid: %v", cause),
			cause,
		),
	}
}

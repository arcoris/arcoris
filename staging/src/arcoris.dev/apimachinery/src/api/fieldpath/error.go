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

package fieldpath

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidPath classifies malformed semantic paths.
	ErrInvalidPath = errors.New("invalid field path")
	// ErrInvalidSyntax classifies malformed field-path text.
	ErrInvalidSyntax = errors.New("invalid field path syntax")
	// ErrInvalidElement classifies malformed path elements.
	ErrInvalidElement = errors.New("invalid field path element")
	// ErrInvalidSelector classifies malformed associative-list selectors.
	ErrInvalidSelector = errors.New("invalid field path selector")
	// ErrInvalidEntry classifies malformed selector entries.
	ErrInvalidEntry = errors.New("invalid field path selector entry")
	// ErrInvalidLiteral classifies malformed selector literals.
	ErrInvalidLiteral = errors.New("invalid field path literal")

	// ErrEmptyFieldName classifies empty field names.
	ErrEmptyFieldName = errors.New("empty field name")
	// ErrEmptyKey classifies empty map keys.
	ErrEmptyKey = errors.New("empty key")
	// ErrNegativeIndex classifies negative list indexes.
	ErrNegativeIndex = errors.New("negative index")
	// ErrEmptySelector classifies selectors with no entries.
	ErrEmptySelector = errors.New("empty selector")
	// ErrDuplicateField classifies repeated selector field names.
	ErrDuplicateField = errors.New("duplicate selector field")
)

// Error is the structured diagnostic returned by field-path validation.
type Error struct {
	// Err is the broad sentinel used with errors.Is.
	Err error
	// Reason gives stable machine-readable detail within Err.
	Reason ErrorReason
	// Detail gives human-readable diagnostic context.
	Detail string
	// Cause preserves more specific nested failures.
	Cause error
}

// Error returns a compact human-readable field-path diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"fieldpath"}

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

// Unwrap preserves broad sentinels and nested causes.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	if e.Err != nil && e.Cause != nil {
		return errors.Join(e.Err, e.Cause)
	}

	if e.Cause != nil {
		return e.Cause
	}

	return e.Err
}

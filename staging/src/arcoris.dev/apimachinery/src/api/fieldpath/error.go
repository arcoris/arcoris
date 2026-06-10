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

	"arcoris.dev/apimachinery/api/internal/diagnostic"
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
	// ErrEmptyMapKey classifies empty map keys.
	ErrEmptyMapKey = errors.New("empty map key")
	// ErrNegativeIndex classifies negative list indexes.
	ErrNegativeIndex = errors.New("negative index")
	// ErrEmptySelector classifies selectors with no entries.
	ErrEmptySelector = errors.New("empty selector")
	// ErrDuplicateSelectorField classifies repeated selector field names.
	ErrDuplicateSelectorField = errors.New("duplicate selector field")
	// ErrNonCanonicalText classifies valid path text with a non-canonical spelling.
	ErrNonCanonicalText = errors.New("non-canonical field path text")
)

// Error is the structured diagnostic returned by field-path validation.
type Error struct {
	// Record stores the shared sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable field-path diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("fieldpath")
}

// Unwrap preserves broad sentinels and nested causes.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

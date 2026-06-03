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

package objectownership

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidDocument classifies malformed object ownership documents.
	ErrInvalidDocument = errors.New("invalid object ownership document")

	// ErrUnsupportedVersion classifies unknown document versions.
	ErrUnsupportedVersion = errors.New("unsupported object ownership document version")

	// ErrInvalidSurface classifies malformed surface ownership data.
	ErrInvalidSurface = errors.New("invalid object ownership surface")

	// ErrInvalidEntry classifies malformed owner/path records.
	ErrInvalidEntry = errors.New("invalid object ownership entry")

	// ErrInvalidPath classifies malformed document field path text.
	ErrInvalidPath = errors.New("invalid object ownership path")
)

// Error is the structured diagnostic returned by objectownership.
type Error struct {
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable objectownership diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("objectownership")
}

// Unwrap exposes the broad sentinel and nested cause for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

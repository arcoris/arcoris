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

package fieldownership

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidOwner classifies owner text that cannot identify ownership.
	ErrInvalidOwner = errors.New("invalid field owner")

	// ErrInvalidEntry classifies owner/field-set records that cannot be stored.
	ErrInvalidEntry = errors.New("invalid field ownership entry")

	// ErrInvalidState classifies malformed ownership state representation.
	ErrInvalidState = errors.New("invalid field ownership state")

	// ErrInvalidOwnedPath classifies malformed owner/path query records.
	ErrInvalidOwnedPath = errors.New("invalid owned field path")

	// ErrInvalidConflict classifies malformed conflict records.
	ErrInvalidConflict = errors.New("invalid field ownership conflict record")

	// ErrInvalidPath classifies malformed semantic field paths.
	ErrInvalidPath = errors.New("invalid field path")

	// ErrConflict classifies overlapping ownership discovered for an attempt.
	ErrConflict = errors.New("field ownership conflict")
)

// Error is the structured diagnostic returned for validation failures.
type Error struct {
	// Record stores the shared location, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable fieldownership diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("fieldownership")
}

// Unwrap exposes the broad sentinel and any nested cause for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

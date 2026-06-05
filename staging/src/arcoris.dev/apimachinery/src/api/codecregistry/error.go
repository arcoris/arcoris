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

package codecregistry

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidCodec classifies nil or otherwise unusable codec implementations.
	ErrInvalidCodec = errors.New("invalid codec")

	// ErrInvalidInfo classifies malformed or non-normalizable codec metadata.
	ErrInvalidInfo = errors.New("invalid codec info")

	// ErrInvalidEntryID classifies malformed registry entry identities.
	ErrInvalidEntryID = errors.New("invalid codec registry entry ID")

	// ErrDuplicateEntryID classifies duplicate configured codec identities.
	ErrDuplicateEntryID = errors.New("duplicate codec registry entry ID")

	// ErrInvalidRegistration classifies missing or incomplete registrations.
	ErrInvalidRegistration = errors.New("invalid codec registration")

	// ErrCapabilityMismatch classifies disagreement between Info.Targets and
	// implemented byte or streaming capability interfaces.
	ErrCapabilityMismatch = errors.New("codec capability mismatch")
)

// Error is the structured diagnostic returned by registry construction.
type Error struct {
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a stable human-readable registry diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("codecregistry")
}

// Unwrap exposes broad sentinels and nested codec causes for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

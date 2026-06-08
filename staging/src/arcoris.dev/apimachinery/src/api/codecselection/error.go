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

package codecselection

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidContentType classifies malformed normalized content type keys.
	ErrInvalidContentType = errors.New("invalid codec selection content type")

	// ErrInvalidParameters classifies malformed content type parameters.
	ErrInvalidParameters = errors.New("invalid codec selection parameters")

	// ErrInvalidPreference classifies malformed encode preferences.
	ErrInvalidPreference = errors.New("invalid codec selection preference")

	// ErrInvalidBinding classifies malformed decode or encode bindings.
	ErrInvalidBinding = errors.New("invalid codec selection binding")

	// ErrDuplicateDecodeBinding classifies ambiguous decode binding keys.
	ErrDuplicateDecodeBinding = errors.New("duplicate codec decode binding")

	// ErrDuplicateEncodeBinding classifies ambiguous encode binding keys.
	ErrDuplicateEncodeBinding = errors.New("duplicate codec encode binding")

	// ErrUnknownEntryID classifies bindings that reference absent registry entries.
	ErrUnknownEntryID = errors.New("unknown codec registry entry ID")

	// ErrEntryMediaTypeMismatch classifies EntryID bindings to undeclared media types.
	ErrEntryMediaTypeMismatch = errors.New("codec registry entry media type mismatch")

	// ErrEntryTargetMismatch classifies EntryID bindings to undeclared targets.
	ErrEntryTargetMismatch = errors.New("codec registry entry target mismatch")

	// ErrEntryCapabilityMismatch classifies EntryID bindings to missing capabilities.
	ErrEntryCapabilityMismatch = errors.New("codec registry entry capability mismatch")

	// ErrNoDecodeBinding classifies runtime decode misses.
	ErrNoDecodeBinding = errors.New("no matching codec decode binding")

	// ErrNoEncodePreference classifies runtime encode preference misses.
	ErrNoEncodePreference = errors.New("no supported codec encode preference")
)

// Error is the structured diagnostic returned by codec selection.
type Error struct {
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a stable human-readable selection diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("codecselection")
}

// Unwrap exposes broad sentinels and nested causes for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

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

package objectapply

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidRequest classifies malformed objectapply request shapes.
	ErrInvalidRequest = errors.New("invalid object apply request")

	// ErrInvalidOwner classifies malformed Desired field owner identities.
	ErrInvalidOwner = errors.New("invalid field owner")

	// ErrInvalidResource classifies missing or unusable resource definitions.
	ErrInvalidResource = errors.New("invalid resource")

	// ErrInvalidObject classifies live or applied object envelope failures.
	ErrInvalidObject = errors.New("invalid object")

	// ErrIdentityMismatch classifies name, namespace, group, kind, or UID mismatch.
	ErrIdentityMismatch = errors.New("object identity mismatch")

	// ErrVersionMismatch classifies cross-version apply attempts.
	ErrVersionMismatch = errors.New("object version mismatch")

	// ErrUnsupportedObservedApply classifies attempts to apply Observed.
	ErrUnsupportedObservedApply = errors.New("unsupported observed apply")

	// ErrUnsupportedMetadataChange classifies attempts to apply metadata.
	ErrUnsupportedMetadataChange = errors.New("unsupported metadata change")

	// ErrDesiredApplyFailed classifies non-conflict failures from valueapply.
	ErrDesiredApplyFailed = errors.New("desired apply failed")

	// ErrConflict classifies Desired ownership conflicts.
	ErrConflict = errors.New("field ownership conflict")
)

// Error is the structured diagnostic returned by objectapply.
//
// The type mirrors other api/* packages: Err is a broad sentinel for errors.Is,
// Reason is the stable objectapply-specific machine code, Path is a
// domain-specific object/request location, and Cause preserves lower layers.
type Error struct {
	// Record stores shared diagnostic fields while keeping ErrorReason local.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable objectapply diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("objectapply")
}

// Unwrap exposes the sentinel and nested cause for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

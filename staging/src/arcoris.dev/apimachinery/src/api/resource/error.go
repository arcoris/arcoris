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

package resource

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// Broad validation sentinels preserved for errors.Is checks.
var (
	// ErrInvalidDefinition classifies invalid resource-family definitions.
	ErrInvalidDefinition = errors.New("invalid API resource definition")
	// ErrInvalidVersion classifies invalid version-level resource descriptors.
	ErrInvalidVersion = errors.New("invalid API resource version")
	// ErrInvalidScope classifies invalid resource scope values.
	ErrInvalidScope = errors.New("invalid API resource scope")
	// ErrInvalidJSON classifies invalid JSON scalar encoding for resource values.
	ErrInvalidJSON = errors.New("invalid API resource JSON")
	// ErrNilReceiver classifies nil pointer decoding receivers.
	ErrNilReceiver = errors.New("nil API resource receiver")
)

// Error is a structured resource-definition diagnostic.
//
// Err is the broad sentinel used for errors.Is. Reason and Detail describe the
// exact invariant for humans, tests, CLIs, and future tooling. Cause preserves
// nested api/identity or api/types diagnostics.
type Error struct {
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a stable human-readable diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("resource")
}

// Unwrap preserves broad and nested error identities.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Record.Unwrap()
}

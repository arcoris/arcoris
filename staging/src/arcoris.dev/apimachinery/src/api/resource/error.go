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
	"strings"
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
	// Path identifies the descriptor location that failed validation.
	Path string

	// Err is the broad sentinel used for errors.Is classification.
	Err error

	// Reason is the precise invariant failure within Err's broad category.
	Reason ErrorReason

	// Detail is a human-facing explanation with descriptor-specific context.
	Detail string

	// Cause preserves nested api/identity, api/types, or JSON diagnostics.
	Cause error
}

// Error returns a stable human-readable diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	var parts []string
	parts = append(parts, "resource")
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

// Unwrap preserves broad and nested error identities.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	if e.Err == nil {
		return e.Cause
	}
	if e.Cause == nil {
		return e.Err
	}
	return errors.Join(e.Err, e.Cause)
}

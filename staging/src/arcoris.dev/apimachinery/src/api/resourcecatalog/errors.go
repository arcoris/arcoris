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

package resourcecatalog

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidCatalog classifies catalog registration failures caused by an
	// invalid resource definition.
	ErrInvalidCatalog = errors.New("invalid API resource catalog")

	// ErrDefinitionExists classifies a registration that conflicts with stored
	// catalog state.
	ErrDefinitionExists = errors.New("API resource definition already exists")

	// ErrDuplicateDefinition classifies duplicate identities within one
	// registration batch.
	ErrDuplicateDefinition = errors.New("duplicate API resource definition")

	// ErrNilCatalog classifies write operations on a nil Catalog receiver.
	ErrNilCatalog = errors.New("nil API resource catalog")
)

// Error is a structured resource-catalog diagnostic.
//
// Err is the catalog-level sentinel used by errors.Is. Reason and Detail
// describe the exact catalog invariant for humans, tests, CLIs, and future
// tooling. Cause preserves nested api/resource diagnostics.
type Error struct {
	// Path identifies the catalog operation location that failed.
	Path string

	// Err is the broad catalog sentinel used for errors.Is classification.
	Err error

	// Reason is the precise catalog invariant failure within Err's category.
	Reason ErrorReason

	// Detail is a human-facing explanation with identity-specific context.
	Detail string

	// Cause preserves nested api/resource diagnostics.
	Cause error
}

// Error returns a stable human-readable diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"resourcecatalog"}
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

// Unwrap preserves catalog and nested resource error identities.
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

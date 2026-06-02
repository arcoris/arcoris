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

package listmapkey

import (
	"strings"

	"arcoris.dev/apimachinery/api/fieldpath"
)

// Error reports failed ListMap key extraction.
//
// The error intentionally carries internal failure classification only. Public
// packages wrap it into their own sentinel/reason models so this internal
// package does not leak a public diagnostic contract.
type Error struct {
	// Path is the concrete payload location that prevented key extraction.
	Path fieldpath.Path

	// Kind classifies the failure for callers that need fallback behavior.
	Kind FailureKind

	// Detail gives human-readable context for diagnostics.
	Detail string

	// Cause preserves lower-level failures such as fieldpath selector errors.
	Cause error
}

// Error returns a compact diagnostic string.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"listmapkey", e.Path.String()}
	if e.Kind != "" {
		parts = append(parts, string(e.Kind))
	}
	if e.Detail != "" {
		parts = append(parts, e.Detail)
	}

	return strings.Join(parts, ": ")
}

// Unwrap preserves the lower-level cause when one exists.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

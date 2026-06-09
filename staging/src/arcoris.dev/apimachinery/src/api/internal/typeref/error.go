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

package typeref

import (
	"errors"
	"strings"

	"arcoris.dev/apimachinery/api/fieldpath"
)

// Error reports a DescriptorRef traversal failure without imposing a public error
// taxonomy on the caller.
type Error struct {
	// Path is the semantic payload path whose descriptor reference blocked
	// traversal.
	Path fieldpath.Path

	// Kind classifies the failure for caller-specific sentinel mapping.
	Kind FailureKind

	// Detail gives human-readable diagnostic context.
	Detail string

	// Cause preserves lower-level failures when traversal is blocked by one.
	Cause error
}

// Error returns a compact diagnostic string for internal resolver failures.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"typeref", e.Path.String()}
	if e.Kind != "" {
		parts = append(parts, string(e.Kind))
	}
	if e.Detail != "" {
		parts = append(parts, e.Detail)
	}

	return strings.Join(parts, ": ")
}

// Unwrap preserves lower-level traversal causes.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

// AsError extracts a DescriptorRef traversal error from an arbitrary error.
func AsError(err error) (*Error, bool) {
	var refError *Error
	if errors.As(err, &refError) {
		return refError, true
	}

	return nil, false
}

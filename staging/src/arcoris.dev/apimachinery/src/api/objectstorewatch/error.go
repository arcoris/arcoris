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

package objectstorewatch

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidCollectionRead classifies malformed list-to-watch collection reads.
	ErrInvalidCollectionRead = errors.New("invalid object store watch collection read")
	// ErrInvalidBoundary classifies malformed object store watch boundaries.
	ErrInvalidBoundary = errors.New("invalid object store watch boundary")
)

// errorFor builds the package's standard structured diagnostic.
//
// Call sites pass the path, reason, broad sentinel, and lower cause directly so
// validation code stays explicit about the exact contract location that failed.
// The returned value preserves errors.Is for the sentinel and cause, and
// errors.As for *Error.
func errorFor(path string, reason ErrorReason, sentinel error, cause error) error {
	if cause == nil {
		cause = sentinel
	}
	if cause == nil {
		cause = errors.New("object store watch contract violation")
	}

	diagnostic := &Error{Path: path, Reason: reason, Cause: cause}
	if sentinel == nil {
		return diagnostic
	}

	return errors.Join(sentinel, diagnostic, cause)
}

// Error carries one structured objectstorewatch diagnostic.
//
// Package functions usually return Error through errors.Join together with a
// broad sentinel and the lower validation failure. That shape lets callers use
// errors.Is for stable classes, errors.As for structured diagnostics, and still
// inspect the underlying objectstore or objectwatch cause.
type Error struct {
	// Path identifies the logical contract location that failed. Paths are
	// stable enough for tests and logs, but they are not a parser-facing API.
	Path string
	// Reason is the stable machine-readable failure class for this diagnostic.
	Reason ErrorReason
	// Cause is the lower validation or conversion failure. It may be nil for
	// manually constructed diagnostics; Error and Unwrap are nil-safe.
	Cause error
}

// Error returns compact diagnostic text and is safe on nil receivers.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"objectstorewatch"}
	if e.Path != "" {
		parts = append(parts, e.Path)
	}
	if e.Reason.IsValid() {
		parts = append(parts, e.Reason.String())
	}
	if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}

	return strings.Join(parts, ": ")
}

// Unwrap returns the lower cause for errors.Is and errors.As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

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

package object

import (
	"errors"
	"strings"
)

// Object sentinels classify broad metadata validation failures.
var (
	// ErrInvalidObject classifies object envelopes with invalid metadata.
	ErrInvalidObject = errors.New("invalid API object")
	// ErrInvalidList classifies list envelopes with invalid metadata.
	ErrInvalidList = errors.New("invalid API object list")
)

// Error is the structured diagnostic returned by object metadata validation.
type Error struct {
	// Path identifies the object or list metadata field that failed validation.
	Path string
	// Err is the broad sentinel used with errors.Is.
	Err error
	// Reason gives stable machine-readable detail within Err.
	Reason ErrorReason
	// Detail gives human-readable context for logs and diagnostics.
	Detail string
	// Cause preserves nested metadata validation failures.
	Cause error
}

// Error returns a compact human-readable object diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"object"}

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

// Unwrap preserves both the broad sentinel and nested metadata cause.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	if e.Err != nil && e.Cause != nil {
		return errors.Join(e.Err, e.Cause)
	}

	if e.Err != nil {
		return e.Err
	}

	return e.Cause
}

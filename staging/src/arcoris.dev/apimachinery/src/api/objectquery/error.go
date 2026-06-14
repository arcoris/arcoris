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

package objectquery

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidQuery classifies malformed top-level query values.
	ErrInvalidQuery = errors.New("invalid object query")

	// ErrInvalidSelector classifies malformed metadata selector values.
	ErrInvalidSelector = errors.New("invalid object query selector")

	// ErrInvalidRequirement classifies malformed metadata requirements.
	ErrInvalidRequirement = errors.New("invalid object query requirement")

	// ErrInvalidOperator classifies unknown requirement operators.
	ErrInvalidOperator = errors.New("invalid object query operator")
)

// Error is the structured diagnostic returned by objectquery validation.
type Error struct {
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable objectquery diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("objectquery")
}

// Unwrap exposes broad sentinels and lower validation causes for errors.Is.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

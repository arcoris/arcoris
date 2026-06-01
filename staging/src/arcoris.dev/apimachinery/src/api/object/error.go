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

	"arcoris.dev/apimachinery/api/internal/diagnostic"
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
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable object diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("object")
}

// Unwrap preserves both the broad sentinel and nested metadata cause.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

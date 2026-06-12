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

package codecjson

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidJSON classifies malformed JSON syntax or token shape.
	ErrInvalidJSON = errors.New("invalid JSON")

	// ErrDuplicateKey classifies duplicate JSON object member names.
	ErrDuplicateKey = errors.New("duplicate JSON object key")

	// ErrTrailingData classifies non-whitespace data after one JSON document.
	ErrTrailingData = errors.New("trailing JSON data")

	// ErrUnsupportedValue classifies values that generic JSON cannot round-trip.
	ErrUnsupportedValue = errors.New("unsupported JSON value")

	// ErrInvalidNumber classifies JSON numbers outside the value model.
	ErrInvalidNumber = errors.New("invalid JSON number")

	// ErrInvalidEnvelope classifies malformed object or ownership JSON shapes.
	ErrInvalidEnvelope = errors.New("invalid JSON object envelope")
)

// Error is the structured diagnostic returned by codecjson.
type Error struct {
	// Record stores JSON path, local sentinel, reason, detail, and cause.
	diagnostic.Record[ErrorReason]

	// CodecErr stores the broad api/codec classification for errors.Is.
	CodecErr error
}

// Error returns a stable human-readable JSON codec diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("codecjson")
}

// Unwrap exposes local, root codec, and nested causes for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return errors.Join(e.Err, e.CodecErr, e.Cause)
}

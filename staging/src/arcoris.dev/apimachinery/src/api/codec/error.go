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

package codec

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidFormat classifies malformed codec format identifiers.
	//
	// This sentinel is used for syntax or canonicalization failures in Format,
	// including empty values, non-lowercase text, and invalid token characters.
	ErrInvalidFormat = errors.New("invalid codec format")

	// ErrInvalidMediaType classifies malformed codec media types.
	//
	// This sentinel is used for media type shape failures, not for negotiation
	// misses. Unsupported but well-formed media types use ErrUnsupportedMediaType.
	ErrInvalidMediaType = errors.New("invalid codec media type")

	// ErrInvalidTarget classifies malformed or unknown target names in codec metadata.
	//
	// Target is closed-world in this package, so unknown target names are invalid
	// declarations rather than extension points.
	ErrInvalidTarget = errors.New("invalid codec target")

	// ErrInvalidInfo classifies malformed codec implementation metadata.
	//
	// Info errors should wrap the more specific invalid format, media type, or
	// target cause when one exists.
	ErrInvalidInfo = errors.New("invalid codec info")

	// ErrUnsupportedFormat classifies format choices unsupported by a codec.
	//
	// Concrete codecs or future registries can use this when a valid Format is
	// recognized as unavailable in the current context.
	ErrUnsupportedFormat = errors.New("unsupported codec format")

	// ErrUnsupportedMediaType classifies media types unsupported by a codec.
	ErrUnsupportedMediaType = errors.New("unsupported codec media type")

	// ErrUnsupportedTarget classifies API document targets unsupported by a codec.
	ErrUnsupportedTarget = errors.New("unsupported codec target")

	// ErrUnsupportedFeature classifies optional codec behavior that is unavailable.
	//
	// Examples include unsupported pretty output, unsupported streaming, or
	// format-specific strict behavior that a concrete implementation cannot
	// provide.
	ErrUnsupportedFeature = errors.New("unsupported codec feature")

	// ErrDecodeFailed classifies failures while reading encoded document data.
	//
	// Concrete decoders should wrap lower-level syntax, I/O, or representation
	// causes so callers can retain both the broad decode classification and the
	// specific failure.
	ErrDecodeFailed = errors.New("decode failed")

	// ErrEncodeFailed classifies failures while writing encoded document data.
	ErrEncodeFailed = errors.New("encode failed")

	// ErrInvalidDocument classifies malformed encoded document content.
	ErrInvalidDocument = errors.New("invalid encoded document")

	// ErrStrictViolation classifies input rejected by strict decode mode.
	ErrStrictViolation = errors.New("strict codec violation")

	// ErrMaxDepthExceeded classifies documents exceeding configured nesting depth.
	ErrMaxDepthExceeded = errors.New("maximum codec depth exceeded")

	// ErrInvalidNumber classifies encoded numbers outside the API value model.
	//
	// Codecs should prefer this over silent precision loss when an encoded number
	// cannot be represented faithfully.
	ErrInvalidNumber = errors.New("invalid encoded number")
)

// Error is the structured diagnostic returned by codec contracts and codecs.
//
// The embedded diagnostic record provides a stable diagnostic path, broad
// sentinel, machine-readable reason, human detail, and optional nested cause.
type Error struct {
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a stable human-readable codec diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("codec")
}

// Unwrap exposes broad sentinels and nested causes for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

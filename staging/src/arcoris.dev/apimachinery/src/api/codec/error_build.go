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
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

const (
	// pathCodecFormat names standalone Format diagnostics.
	pathCodecFormat = "codec.format"

	// pathCodecMediaType names standalone MediaType diagnostics.
	pathCodecMediaType = "codec.mediaType"

	// pathCodecTarget names standalone Target diagnostics.
	pathCodecTarget = "codec.target"

	// pathCodecInfo names Info diagnostics.
	pathCodecInfo = "codec.info"
)

// ErrorAt creates a structured codec diagnostic without a nested cause.
//
// Concrete codec packages should use this constructor when they can classify a
// failure directly with a codec sentinel. A nil sentinel returns nil so callers
// can compose builders without special-case branches.
func ErrorAt(path string, err error, reason ErrorReason, detail string) error {
	if err == nil {
		return nil
	}

	return &Error{
		Record: diagnostic.NewRecord(normalizeDiagnosticPath(path), err, reason, detail),
	}
}

// WrapAt creates a structured codec diagnostic preserving a nested cause.
//
// Use WrapAt when adapting lower-level parser, I/O, or model errors into the
// codec error vocabulary. The nested cause remains discoverable through
// errors.Is and errors.As.
func WrapAt(path string, err error, reason ErrorReason, detail string, cause error) error {
	if err == nil {
		return nil
	}

	return &Error{
		Record: diagnostic.WrapRecord(normalizeDiagnosticPath(path), err, reason, detail, cause),
	}
}

// errorfAt creates a structured codec diagnostic with formatted detail text.
//
// The helper keeps validation code readable while still returning the same
// structured Error shape as ErrorAt.
func errorfAt(path string, err error, reason ErrorReason, format string, args ...any) error {
	return ErrorAt(path, err, reason, fmt.Sprintf(format, args...))
}

// normalizeDiagnosticPath returns the stable root path for empty diagnostics.
//
// Codec paths are syntactic document paths, not api/fieldpath.Path values.
// Empty paths normalize to "$" so all public diagnostics have a location.
func normalizeDiagnosticPath(path string) string {
	if path == "" {
		return "$"
	}

	return path
}

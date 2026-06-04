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

import "strings"

// Validate checks that m is a canonical open-world media type.
//
// Validate accepts exactly "type/subtype" with lowercase ASCII token parts. It
// rejects parameters, comma-separated lists, whitespace, and non-canonical case.
func (m MediaType) Validate() error {
	return validateMediaTypeAt(pathCodecMediaType, m)
}

// validateMediaTypeAt checks m at a caller-provided diagnostic path.
//
// Info validation uses this helper so nested media type failures keep stable
// indexed paths such as codec.info.mediaTypes[0].
func validateMediaTypeAt(path string, m MediaType) error {
	text := m.String()
	if text == "" {
		return ErrorAt(
			path,
			ErrInvalidMediaType,
			ErrorReasonInvalidMediaType,
			"codec media type is required",
		)
	}
	if strings.ContainsAny(text, ";,") {
		return errorfAt(
			path,
			ErrInvalidMediaType,
			ErrorReasonInvalidMediaType,
			"codec media type %q must not contain parameters or list separators",
			text,
		)
	}

	parts := strings.Split(text, "/")
	if len(parts) != 2 {
		return errorfAt(
			path,
			ErrInvalidMediaType,
			ErrorReasonInvalidMediaType,
			"codec media type %q must be type/subtype",
			text,
		)
	}
	if parts[0] == "" {
		return errorfAt(path, ErrInvalidMediaType, ErrorReasonInvalidMediaType, "codec media type must have a type")
	}
	if parts[1] == "" {
		return errorfAt(path, ErrInvalidMediaType, ErrorReasonInvalidMediaType, "codec media type must have a subtype")
	}
	if err := validateCodecToken(
		path,
		parts[0],
		mediaTypeTokenOptions(),
		ErrInvalidMediaType,
		ErrorReasonInvalidMediaType,
		"codec media type",
	); err != nil {
		return err
	}
	if err := validateCodecToken(
		path,
		parts[1],
		mediaTypeTokenOptions(),
		ErrInvalidMediaType,
		ErrorReasonInvalidMediaType,
		"codec media subtype",
	); err != nil {
		return err
	}

	return nil
}

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

import "arcoris.dev/apimachinery/api/internal/lexical"

// Validate checks that f is a canonical open-world format token.
//
// Validate is strict: callers must pass already-normalized lowercase text.
func (f Format) Validate() error {
	return validateFormatAt(pathCodecFormat, f)
}

// validateFormatAt checks f at a caller-provided diagnostic path.
//
// Info validation uses this helper so nested format failures point at
// codec.info.format instead of the standalone codec.format path.
func validateFormatAt(path string, f Format) error {
	text := f.String()
	if err := validateCodecToken(
		path,
		text,
		formatTokenOptions(),
		ErrInvalidFormat,
		ErrorReasonInvalidFormat,
		"codec format",
	); err != nil {
		return err
	}
	if !lexical.IsASCIILower(text[0]) && !lexical.IsASCIIDigit(text[0]) {
		return errorfAt(
			path,
			ErrInvalidFormat,
			ErrorReasonInvalidFormat,
			"codec format %q must start with a lowercase letter or digit",
			text,
		)
	}

	return nil
}

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
	"strings"

	"arcoris.dev/apimachinery/api/internal/lexical"
)

// formatTokenOptions returns the open-world Format token grammar.
//
// The grammar allows lowercase ASCII letters, digits, ".", "_", and "-" after
// the first byte. The first-byte rule is enforced by validateFormatAt because
// lexical.TokenOptions only describes per-byte character sets.
func formatTokenOptions() lexical.TokenOptions {
	return lexical.TokenOptions{
		MinLength:       1,
		AllowLower:      true,
		AllowDigit:      true,
		AllowHyphen:     true,
		AllowDot:        true,
		AllowUnderscore: true,
	}
}

// mediaTypeTokenOptions returns the type/subtype token grammar for MediaType.
//
// The grammar allows "+" so vendor and structured syntax suffixes such as
// "application/vnd.arcoris.object+json" remain valid without codec-local byte
// predicates.
func mediaTypeTokenOptions() lexical.TokenOptions {
	return lexical.TokenOptions{
		MinLength:       1,
		AllowLower:      true,
		AllowDigit:      true,
		AllowHyphen:     true,
		AllowDot:        true,
		AllowUnderscore: true,
		AllowPlus:       true,
	}
}

// normalizeCodecToken trims surrounding whitespace and canonicalizes case.
//
// It is shared by Format.Normalize and MediaType.Normalize. Validate methods
// deliberately do not call this helper.
func normalizeCodecToken(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}

// validateCodecToken maps internal lexical failures into codec diagnostics.
//
// api/internal/lexical owns the byte-level ASCII grammar. The codec package
// owns public sentinel errors, reasons, diagnostic paths, and user-facing detail
// wording, so lexical violations are translated here before leaving the package.
func validateCodecToken(
	path string,
	text string,
	opts lexical.TokenOptions,
	err error,
	reason ErrorReason,
	name string,
) error {
	if violation := lexical.ValidateASCIIToken(text, opts); violation != nil {
		return errorfAt(
			path,
			err,
			reason,
			"%s is invalid: %s",
			name,
			violation.Detail,
		)
	}

	return nil
}

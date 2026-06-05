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

package lexical

import "fmt"

// TokenOptions describes a small ASCII token grammar.
//
// The options are internal mechanics for future descriptor identifiers. They do
// not define a public ARCORIS name type.
type TokenOptions struct {
	// MinLength rejects values shorter than this byte count when positive.
	MinLength int

	// MaxLength rejects values longer than this byte count when positive.
	MaxLength int

	// AllowLower allows ASCII lowercase letters.
	AllowLower bool

	// AllowUpper allows ASCII uppercase letters.
	AllowUpper bool

	// AllowDigit allows ASCII decimal digits.
	AllowDigit bool

	// AllowHyphen allows "-".
	AllowHyphen bool

	// AllowDot allows ".".
	AllowDot bool

	// AllowUnderscore allows "_".
	AllowUnderscore bool

	// AllowPlus allows "+".
	//
	// This is useful for protocol tokens that embed media-type structured suffix
	// grammar, such as "object+json", while still avoiding full media type
	// parsing in the lexical package.
	AllowPlus bool

	// AllowSlash allows "/".
	//
	// This is useful for owner-defined catalog identifiers that use path-like
	// grouping without making the lexical package understand paths.
	AllowSlash bool

	// RequireAlnumEdges requires non-empty values to start and end with an
	// ASCII letter or digit.
	RequireAlnumEdges bool

	// RejectAdjacentSeparators rejects adjacent allowed punctuation bytes.
	//
	// Separator bytes are "-", ".", "_", "+", and "/" when those bytes are
	// enabled by this options value. Letter and digit bytes reset the separator
	// run.
	RejectAdjacentSeparators bool
}

// ValidateASCIIToken validates value against opts.
//
// The helper is deliberately byte-oriented and does not trim or normalize.
func ValidateASCIIToken(value string, opts TokenOptions) *Violation {
	if opts.MinLength > 0 && len(value) < opts.MinLength {
		if value == "" {
			return violation(ReasonEmptyValue, "ASCII token must be non-empty")
		}
		return violation(
			ReasonInvalidLength,
			fmt.Sprintf("ASCII token length must be >= %d bytes", opts.MinLength),
		)
	}
	if opts.MaxLength > 0 && len(value) > opts.MaxLength {
		return violation(
			ReasonInvalidLength,
			fmt.Sprintf("ASCII token length must be <= %d bytes", opts.MaxLength),
		)
	}
	if opts.RequireAlnumEdges && len(value) > 0 {
		if !IsASCIIAlnum(value[0]) || !IsASCIIAlnum(value[len(value)-1]) {
			return violation(ReasonInvalidEdge, "ASCII token must start and end with an ASCII letter or digit")
		}
	}

	for i := 0; i < len(value); i++ {
		b := value[i]
		if !tokenByteAllowed(b, opts) {
			return violation(
				ReasonInvalidCharacter,
				fmt.Sprintf("ASCII token contains invalid byte %q at index %d", b, i),
			)
		}
		if opts.RejectAdjacentSeparators && i > 0 && tokenByteSeparator(b, opts) && tokenByteSeparator(value[i-1], opts) {
			return violation(
				ReasonInvalidForm,
				fmt.Sprintf("ASCII token contains adjacent separator bytes at index %d", i),
			)
		}
	}

	return nil
}

// tokenByteAllowed reports whether b is enabled by opts.
func tokenByteAllowed(b byte, opts TokenOptions) bool {
	switch {
	case opts.AllowLower && IsASCIILower(b):
		return true
	case opts.AllowUpper && IsASCIIUpper(b):
		return true
	case opts.AllowDigit && IsASCIIDigit(b):
		return true
	case opts.AllowHyphen && b == '-':
		return true
	case opts.AllowDot && b == '.':
		return true
	case opts.AllowUnderscore && b == '_':
		return true
	case opts.AllowPlus && b == '+':
		return true
	case opts.AllowSlash && b == '/':
		return true
	default:
		return false
	}
}

// tokenByteSeparator reports whether b is an enabled punctuation separator.
func tokenByteSeparator(b byte, opts TokenOptions) bool {
	switch b {
	case '-':
		return opts.AllowHyphen
	case '.':
		return opts.AllowDot
	case '_':
		return opts.AllowUnderscore
	case '+':
		return opts.AllowPlus
	case '/':
		return opts.AllowSlash
	default:
		return false
	}
}

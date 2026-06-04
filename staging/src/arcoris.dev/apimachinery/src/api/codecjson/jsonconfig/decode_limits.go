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

package jsonconfig

const (
	// DefaultMaxDepth is the safe default JSON nesting limit.
	DefaultMaxDepth = 128

	// DefaultMaxNumberDigits bounds exact JSON number expansion by default.
	DefaultMaxNumberDigits = 4096
)

// DecodeLimitsConfig bounds accepted JSON input.
type DecodeLimitsConfig struct {
	// MaxDepth is the inclusive document nesting limit.
	MaxDepth int

	// MaxDocumentBytes bounds raw input bytes. Zero means unlimited.
	MaxDocumentBytes int64

	// MaxStringBytes bounds decoded string and member-name byte length. Zero means unlimited.
	MaxStringBytes int

	// MaxNumberDigits bounds digit and exponent expansion for exact numbers.
	MaxNumberDigits int
}

// defaultDecodeLimitsConfig returns safe decode limits.
func defaultDecodeLimitsConfig() DecodeLimitsConfig {
	return DecodeLimitsConfig{
		MaxDepth:        DefaultMaxDepth,
		MaxNumberDigits: DefaultMaxNumberDigits,
	}
}

// resolveDecodeLimitsConfig applies non-zero required limit defaults.
func resolveDecodeLimitsConfig(limits *DecodeLimitsConfig) {
	if limits.MaxDepth == 0 {
		limits.MaxDepth = DefaultMaxDepth
	}
	if limits.MaxNumberDigits == 0 {
		limits.MaxNumberDigits = DefaultMaxNumberDigits
	}
}

// validateDecodeLimitsConfig checks input size and nesting limits.
func validateDecodeLimitsConfig(limits DecodeLimitsConfig) error {
	switch {
	case limits.MaxDepth <= 0:
		return invalidConfig("decode.limits.max_depth", "must be greater than zero")
	case limits.MaxDocumentBytes < 0:
		return invalidConfig("decode.limits.max_document_bytes", "must be zero or greater")
	case limits.MaxStringBytes < 0:
		return invalidConfig("decode.limits.max_string_bytes", "must be zero or greater")
	case limits.MaxNumberDigits <= 0:
		return invalidConfig("decode.limits.max_number_digits", "must be greater than zero")
	default:
		return nil
	}
}

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

// EncodeLimitsConfig bounds JSON output work.
type EncodeLimitsConfig struct {
	// MaxDepth is the inclusive output nesting limit.
	MaxDepth int

	// MaxOutputBytes bounds encoded JSON bytes. Zero means unlimited.
	MaxOutputBytes int64
}

// defaultEncodeLimitsConfig returns safe output limits.
func defaultEncodeLimitsConfig() EncodeLimitsConfig {
	return EncodeLimitsConfig{MaxDepth: DefaultMaxDepth}
}

// resolveEncodeLimitsConfig applies output limit defaults in place.
func resolveEncodeLimitsConfig(limits *EncodeLimitsConfig) {
	if limits.MaxDepth == 0 {
		limits.MaxDepth = DefaultMaxDepth
	}
}

// validateEncodeLimitsConfig checks output size and nesting limits.
func validateEncodeLimitsConfig(limits EncodeLimitsConfig) error {
	switch {
	case limits.MaxDepth <= 0:
		return invalidConfig("encode.limits.max_depth", "must be greater than zero")
	case limits.MaxOutputBytes < 0:
		return invalidConfig("encode.limits.max_output_bytes", "must be zero or greater")
	default:
		return nil
	}
}

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

// DecodeOwnershipConfig controls object ownership state decoding.
type DecodeOwnershipConfig struct {
	// UnknownFields controls unknown ownership JSON fields.
	UnknownFields UnknownFieldMode

	// ValidateState controls semantic ownership-state validation after decode.
	//
	// A zero Config resolves this to true for the safe default policy. Callers
	// that intentionally disable validation should set UnknownFields explicitly
	// as well, because a bare false bool is otherwise indistinguishable from
	// zero-value default construction.
	ValidateState bool
}

// defaultDecodeOwnershipConfig returns strict ownership decode policy.
func defaultDecodeOwnershipConfig() DecodeOwnershipConfig {
	return DecodeOwnershipConfig{
		UnknownFields: UnknownFieldReject,
		ValidateState: true,
	}
}

// resolveDecodeOwnershipConfig applies ownership decode defaults in place.
func resolveDecodeOwnershipConfig(config *DecodeOwnershipConfig) {
	if config.UnknownFields == UnknownFieldDefault {
		config.UnknownFields = UnknownFieldReject
		config.ValidateState = true
	}
}

// validateDecodeOwnershipConfig checks ownership-state decode policy.
func validateDecodeOwnershipConfig(config DecodeOwnershipConfig) error {
	if !isKnownUnknownFieldMode(config.UnknownFields) {
		return invalidConfig("decode.ownership.unknown_fields", "unknown field mode %d", config.UnknownFields)
	}

	return nil
}

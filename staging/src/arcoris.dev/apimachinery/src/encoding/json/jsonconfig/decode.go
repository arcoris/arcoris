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

// DecodeConfig groups JSON decode policy by document concern.
type DecodeConfig struct {
	// Limits bounds input size and nesting.
	Limits DecodeLimitsConfig

	// Numbers controls JSON number interpretation.
	Numbers DecodeNumberConfig

	// Strings controls JSON string and raw UTF-8 handling.
	Strings DecodeStringConfig

	// Objects controls generic object and object-envelope behavior.
	Objects DecodeObjectConfig

	// Ownership controls object ownership state behavior.
	Ownership DecodeOwnershipConfig
}

// defaultDecodeConfig returns the package's safe API-oriented decode policy.
func defaultDecodeConfig() DecodeConfig {
	return DecodeConfig{
		Limits:    defaultDecodeLimitsConfig(),
		Numbers:   defaultDecodeNumberConfig(),
		Strings:   defaultDecodeStringConfig(),
		Objects:   defaultDecodeObjectConfig(),
		Ownership: defaultDecodeOwnershipConfig(),
	}
}

// resolveDecodeConfig applies decode defaults in place.
func resolveDecodeConfig(config *DecodeConfig) {
	resolveDecodeLimitsConfig(&config.Limits)
	resolveDecodeNumberConfig(&config.Numbers)
	resolveDecodeStringConfig(&config.Strings)
	resolveDecodeObjectConfig(&config.Objects)
	resolveDecodeOwnershipConfig(&config.Ownership)
}

// validateDecodeConfig checks every decode policy group after Resolve.
func validateDecodeConfig(config DecodeConfig) error {
	if err := validateDecodeLimitsConfig(config.Limits); err != nil {
		return err
	}
	if err := validateDecodeNumberConfig(config.Numbers); err != nil {
		return err
	}
	if err := validateDecodeStringConfig(config.Strings); err != nil {
		return err
	}
	if err := validateDecodeObjectConfig(config.Objects); err != nil {
		return err
	}
	if err := validateDecodeOwnershipConfig(config.Ownership); err != nil {
		return err
	}

	return nil
}

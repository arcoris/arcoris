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

// EncodeConfig groups JSON encode policy by document concern.
type EncodeConfig struct {
	// Output controls layout and document framing.
	Output EncodeOutputConfig

	// Ordering controls member ordering for supported document parts.
	Ordering EncodeOrderingConfig

	// Strings controls JSON string escaping.
	Strings EncodeStringConfig

	// Numbers controls JSON number formatting.
	Numbers EncodeNumberConfig

	// Values controls descriptor-dependent generic value kinds.
	Values EncodeValueConfig

	// Object controls value-backed object envelope shape.
	Object EncodeObjectConfig

	// Ownership controls object ownership document shape.
	Ownership EncodeOwnershipConfig

	// Limits bounds output shape and size.
	Limits EncodeLimitsConfig
}

// defaultEncodeConfig returns the package's default JSON writer policy.
func defaultEncodeConfig() EncodeConfig {
	return EncodeConfig{
		Output:    defaultEncodeOutputConfig(),
		Ordering:  defaultEncodeOrderingConfig(),
		Strings:   defaultEncodeStringConfig(),
		Numbers:   defaultEncodeNumberConfig(),
		Values:    defaultEncodeValueConfig(),
		Object:    defaultEncodeObjectConfig(),
		Ownership: defaultEncodeOwnershipConfig(),
		Limits:    defaultEncodeLimitsConfig(),
	}
}

// resolveEncodeConfig applies encode defaults in place.
func resolveEncodeConfig(config *EncodeConfig) {
	resolveEncodeOutputConfig(&config.Output)
	resolveEncodeOrderingConfig(&config.Ordering)
	resolveEncodeStringConfig(&config.Strings)
	resolveEncodeNumberConfig(&config.Numbers)
	resolveEncodeValueConfig(&config.Values)
	resolveEncodeObjectConfig(&config.Object)
	resolveEncodeOwnershipConfig(&config.Ownership)
	resolveEncodeLimitsConfig(&config.Limits)
}

// validateEncodeConfig checks every encode policy group after Resolve.
func validateEncodeConfig(config EncodeConfig) error {
	if err := validateEncodeOutputConfig(config.Output); err != nil {
		return err
	}
	if err := validateEncodeOrderingConfig(config.Ordering); err != nil {
		return err
	}
	if err := validateEncodeStringConfig(config.Strings); err != nil {
		return err
	}
	if err := validateEncodeNumberConfig(config.Numbers); err != nil {
		return err
	}
	if err := validateEncodeValueConfig(config.Values); err != nil {
		return err
	}
	if err := validateEncodeObjectConfig(config.Object); err != nil {
		return err
	}
	if err := validateEncodeOwnershipConfig(config.Ownership); err != nil {
		return err
	}
	if err := validateEncodeLimitsConfig(config.Limits); err != nil {
		return err
	}

	return nil
}

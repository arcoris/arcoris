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

// Config is the complete construction-time configuration for the JSON codec.
type Config struct {
	// Decode controls JSON parser and target-decoder policy.
	Decode DecodeConfig

	// Encode controls JSON writer policy.
	Encode EncodeConfig
}

// Default returns safe API-oriented JSON codec defaults.
func Default() Config {
	return Config{
		Decode: defaultDecodeConfig(),
		Encode: defaultEncodeConfig(),
	}
}

// Resolve applies defaults and validates the resulting concrete configuration.
func Resolve(config Config) (Config, error) {
	resolved := config
	resolveDecodeConfig(&resolved.Decode)
	resolveEncodeConfig(&resolved.Encode)
	if err := validateResolvedConfig(resolved); err != nil {
		return Config{}, err
	}

	return resolved, nil
}

// Validate checks whether config can resolve to a supported JSON codec policy.
func Validate(config Config) error {
	_, err := Resolve(config)
	return err
}

// validateResolvedConfig checks a fully resolved Config.
func validateResolvedConfig(config Config) error {
	if err := validateDecodeConfig(config.Decode); err != nil {
		return err
	}
	if err := validateEncodeConfig(config.Encode); err != nil {
		return err
	}

	return nil
}

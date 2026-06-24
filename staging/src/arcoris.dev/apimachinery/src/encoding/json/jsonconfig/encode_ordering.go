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

// EncodeOrderingConfig controls ordering-sensitive JSON output.
type EncodeOrderingConfig struct {
	// Mode controls whether existing order is preserved or made deterministic.
	Mode OrderingMode
}

// defaultEncodeOrderingConfig returns order-preserving output.
func defaultEncodeOrderingConfig() EncodeOrderingConfig {
	return EncodeOrderingConfig{Mode: OrderingPreserve}
}

// resolveEncodeOrderingConfig applies ordering defaults in place.
func resolveEncodeOrderingConfig(config *EncodeOrderingConfig) {
	if config.Mode == OrderingDefault {
		config.Mode = OrderingPreserve
	}
}

// validateEncodeOrderingConfig checks deterministic ordering policy.
func validateEncodeOrderingConfig(config EncodeOrderingConfig) error {
	if !isKnownOrderingMode(config.Mode) {
		return invalidConfig("encode.ordering.mode", "unknown ordering mode %d", config.Mode)
	}

	return nil
}

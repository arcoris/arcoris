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

// EncodeOwnershipConfig controls object ownership document output.
type EncodeOwnershipConfig struct {
	// Normalize controls whether objectownership.Normalize is applied.
	Normalize OwnershipNormalizeMode

	// EmptyDesired controls the desired field when the desired surface is empty.
	EmptyDesired EmptyOwnershipSurfaceMode

	// EmptyEntries controls entries emission for empty ownership surfaces.
	EmptyEntries EmptyEntriesMode
}

// defaultEncodeOwnershipConfig returns stable ownership document output policy.
func defaultEncodeOwnershipConfig() EncodeOwnershipConfig {
	return EncodeOwnershipConfig{
		Normalize:    OwnershipNormalizeWhenDeterministic,
		EmptyDesired: EmptyOwnershipSurfaceEmit,
		EmptyEntries: EmptyEntriesEmit,
	}
}

// resolveEncodeOwnershipConfig applies ownership encode defaults in place.
func resolveEncodeOwnershipConfig(config *EncodeOwnershipConfig) {
	if config.Normalize == OwnershipNormalizeDefault {
		config.Normalize = OwnershipNormalizeWhenDeterministic
	}
	if config.EmptyDesired == EmptyOwnershipSurfaceDefault {
		config.EmptyDesired = EmptyOwnershipSurfaceEmit
	}
	if config.EmptyEntries == EmptyEntriesDefault {
		config.EmptyEntries = EmptyEntriesEmit
	}
}

// validateEncodeOwnershipConfig checks ownership document output policy.
func validateEncodeOwnershipConfig(config EncodeOwnershipConfig) error {
	switch {
	case !isKnownOwnershipNormalizeMode(config.Normalize):
		return invalidConfig("encode.ownership.normalize", "unknown normalize mode %d", config.Normalize)
	case !isKnownEmptyOwnershipSurfaceMode(config.EmptyDesired):
		return invalidConfig("encode.ownership.empty_desired", "unknown empty desired mode %d", config.EmptyDesired)
	case !isKnownEmptyEntriesMode(config.EmptyEntries):
		return invalidConfig("encode.ownership.empty_entries", "unknown empty entries mode %d", config.EmptyEntries)
	default:
		return nil
	}
}

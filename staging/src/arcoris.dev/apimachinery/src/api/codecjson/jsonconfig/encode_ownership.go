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

// EncodeOwnershipConfig controls object ownership state output.
type EncodeOwnershipConfig struct {
	// Normalize controls whether objectownership.Normalize is applied.
	Normalize OwnershipNormalizeMode

	// EmptySurfaces controls ownership surface emission when a surface is empty.
	EmptySurfaces EmptyOwnershipSurfaceMode

	// EmptyEntries controls entries emission for empty ownership surfaces.
	EmptyEntries EmptyEntriesMode
}

// defaultEncodeOwnershipConfig returns stable ownership state output policy.
func defaultEncodeOwnershipConfig() EncodeOwnershipConfig {
	return EncodeOwnershipConfig{
		Normalize:     OwnershipNormalizeWhenDeterministic,
		EmptySurfaces: EmptyOwnershipSurfaceEmit,
		EmptyEntries:  EmptyEntriesEmit,
	}
}

// resolveEncodeOwnershipConfig applies ownership encode defaults in place.
func resolveEncodeOwnershipConfig(config *EncodeOwnershipConfig) {
	if config.Normalize == OwnershipNormalizeDefault {
		config.Normalize = OwnershipNormalizeWhenDeterministic
	}
	if config.EmptySurfaces == EmptyOwnershipSurfaceDefault {
		config.EmptySurfaces = EmptyOwnershipSurfaceEmit
	}
	if config.EmptyEntries == EmptyEntriesDefault {
		config.EmptyEntries = EmptyEntriesEmit
	}
}

// validateEncodeOwnershipConfig checks ownership state output policy.
func validateEncodeOwnershipConfig(config EncodeOwnershipConfig) error {
	switch {
	case !isKnownOwnershipNormalizeMode(config.Normalize):
		return invalidConfig("encode.ownership.normalize", "unknown normalize mode %d", config.Normalize)
	case !isKnownEmptyOwnershipSurfaceMode(config.EmptySurfaces):
		return invalidConfig("encode.ownership.empty_surfaces", "unknown empty surfaces mode %d", config.EmptySurfaces)
	case !isKnownEmptyEntriesMode(config.EmptyEntries):
		return invalidConfig("encode.ownership.empty_entries", "unknown empty entries mode %d", config.EmptyEntries)
	default:
		return nil
	}
}

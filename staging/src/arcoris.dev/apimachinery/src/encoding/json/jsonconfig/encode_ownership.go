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
	// EmptySurfaces controls ownership surface emission when a surface is empty.
	EmptySurfaces EmptyOwnershipSurfaceMode
}

// defaultEncodeOwnershipConfig returns stable ownership state output policy.
func defaultEncodeOwnershipConfig() EncodeOwnershipConfig {
	return EncodeOwnershipConfig{
		EmptySurfaces: EmptyOwnershipSurfaceEmit,
	}
}

// resolveEncodeOwnershipConfig applies ownership encode defaults in place.
func resolveEncodeOwnershipConfig(config *EncodeOwnershipConfig) {
	if config.EmptySurfaces == EmptyOwnershipSurfaceDefault {
		config.EmptySurfaces = EmptyOwnershipSurfaceEmit
	}
}

// validateEncodeOwnershipConfig checks ownership state output policy.
func validateEncodeOwnershipConfig(config EncodeOwnershipConfig) error {
	if !isKnownEmptyOwnershipSurfaceMode(config.EmptySurfaces) {
		return invalidConfig("encode.ownership.empty_surfaces", "unknown empty surfaces mode %d", config.EmptySurfaces)
	}

	return nil
}

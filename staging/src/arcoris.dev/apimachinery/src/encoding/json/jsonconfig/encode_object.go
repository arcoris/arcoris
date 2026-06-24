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

// EncodeObjectConfig controls API object envelope output.
type EncodeObjectConfig struct {
	// TypeMeta controls apiVersion/kind emission.
	TypeMeta TypeMetaEncodeMode

	// Metadata controls metadata emission.
	Metadata MetadataEncodeMode

	// Observed controls absent observed payload emission.
	Observed ObservedEncodeMode
}

// defaultEncodeObjectConfig returns stable API object envelope policy.
func defaultEncodeObjectConfig() EncodeObjectConfig {
	return EncodeObjectConfig{
		TypeMeta: TypeMetaOmitZero,
		Metadata: MetadataOmitZero,
		Observed: ObservedOmitAbsent,
	}
}

// resolveEncodeObjectConfig applies object encode defaults in place.
func resolveEncodeObjectConfig(config *EncodeObjectConfig) {
	if config.TypeMeta == TypeMetaDefault {
		config.TypeMeta = TypeMetaOmitZero
	}
	if config.Metadata == MetadataDefault {
		config.Metadata = MetadataOmitZero
	}
	if config.Observed == ObservedDefault {
		config.Observed = ObservedOmitAbsent
	}
}

// validateEncodeObjectConfig checks API object envelope output policy.
func validateEncodeObjectConfig(config EncodeObjectConfig) error {
	switch {
	case !isKnownTypeMetaEncodeMode(config.TypeMeta):
		return invalidConfig("encode.object.type_meta", "unknown type meta mode %d", config.TypeMeta)
	case !isKnownMetadataEncodeMode(config.Metadata):
		return invalidConfig("encode.object.metadata", "unknown metadata mode %d", config.Metadata)
	case !isKnownObservedEncodeMode(config.Observed):
		return invalidConfig("encode.object.observed", "unknown observed mode %d", config.Observed)
	default:
		return nil
	}
}

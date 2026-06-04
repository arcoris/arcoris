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

// DecodeObjectConfig controls JSON object and object-envelope behavior.
type DecodeObjectConfig struct {
	// DuplicateKeys controls duplicate JSON object member names.
	DuplicateKeys DuplicateKeyMode

	// TrailingData controls tokens after the first JSON document.
	TrailingData TrailingDataMode

	// UnknownEnvelopeFields controls unknown value-backed object envelope fields.
	UnknownEnvelopeFields UnknownFieldMode
}

// defaultDecodeObjectConfig returns strict object decode policy.
func defaultDecodeObjectConfig() DecodeObjectConfig {
	return DecodeObjectConfig{
		DuplicateKeys:         DuplicateKeyReject,
		TrailingData:          TrailingDataReject,
		UnknownEnvelopeFields: UnknownFieldReject,
	}
}

// resolveDecodeObjectConfig applies object decode defaults in place.
func resolveDecodeObjectConfig(config *DecodeObjectConfig) {
	if config.DuplicateKeys == DuplicateKeyDefault {
		config.DuplicateKeys = DuplicateKeyReject
	}
	if config.TrailingData == TrailingDataDefault {
		config.TrailingData = TrailingDataReject
	}
	if config.UnknownEnvelopeFields == UnknownFieldDefault {
		config.UnknownEnvelopeFields = UnknownFieldReject
	}
}

// validateDecodeObjectConfig checks object parser and envelope field policy.
func validateDecodeObjectConfig(config DecodeObjectConfig) error {
	switch {
	case !isKnownDuplicateKeyMode(config.DuplicateKeys):
		return invalidConfig("decode.objects.duplicate_keys", "unknown duplicate key mode %d", config.DuplicateKeys)
	case config.DuplicateKeys != DuplicateKeyReject:
		return unsupportedConfig("decode.objects.duplicate_keys", "duplicate key mode %d is not implemented", config.DuplicateKeys)
	case !isKnownTrailingDataMode(config.TrailingData):
		return invalidConfig("decode.objects.trailing_data", "unknown trailing data mode %d", config.TrailingData)
	case config.TrailingData != TrailingDataReject:
		return unsupportedConfig("decode.objects.trailing_data", "trailing data mode %d is not implemented", config.TrailingData)
	case !isKnownUnknownFieldMode(config.UnknownEnvelopeFields):
		return invalidConfig("decode.objects.unknown_envelope_fields", "unknown field mode %d", config.UnknownEnvelopeFields)
	default:
		return nil
	}
}

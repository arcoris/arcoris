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

package codecjson

import "arcoris.dev/apimachinery/api/codecjson/jsonconfig"

// resolveConfig turns public jsonconfig policy into private runtime fields.
func resolveConfig(config jsonconfig.Config) (resolvedConfig, error) {
	public, err := jsonconfig.Resolve(config)
	if err != nil {
		return resolvedConfig{}, err
	}

	return resolvedConfig{
		decode: resolveDecodeConfig(public.Decode),
		encode: resolveEncodeConfig(public.Encode),
	}, nil
}

// resolveDecodeConfig copies resolved decode policy into hot-path fields.
func resolveDecodeConfig(config jsonconfig.DecodeConfig) resolvedDecodeConfig {
	return resolvedDecodeConfig{
		maxDepth:                     config.Limits.MaxDepth,
		maxDocumentBytes:             config.Limits.MaxDocumentBytes,
		maxStringBytes:               config.Limits.MaxStringBytes,
		maxNumberDigits:              config.Limits.MaxNumberDigits,
		rejectUnknownEnvelopeFields:  config.Objects.UnknownEnvelopeFields == jsonconfig.UnknownFieldReject,
		rejectUnknownOwnershipFields: config.Ownership.UnknownFields == jsonconfig.UnknownFieldReject,
		validateOwnershipState:       config.Ownership.Validation == jsonconfig.OwnershipValidationEnable,
	}
}

// resolveEncodeConfig copies resolved encode policy into hot-path fields.
func resolveEncodeConfig(config jsonconfig.EncodeConfig) resolvedEncodeConfig {
	return resolvedEncodeConfig{
		pretty:          config.Output.Layout == jsonconfig.LayoutPretty,
		indent:          config.Output.Indent,
		finalNewline:    config.Output.FinalNewline == jsonconfig.FinalNewlineAppend,
		deterministic:   config.Ordering.Mode == jsonconfig.OrderingDeterministic,
		escapeHTML:      config.Strings.EscapeHTML,
		maxDepth:        config.Limits.MaxDepth,
		maxOutputBytes:  config.Limits.MaxOutputBytes,
		maxNumberDigits: config.Numbers.MaxDigits,
		floatFormat:     config.Numbers.FloatFormat,
		negativeZero:    config.Numbers.NegativeZero,
		typeMeta:        config.Object.TypeMeta,
		metadata:        config.Object.Metadata,
		observed:        config.Object.Observed,
		emptySurfaces:   config.Ownership.EmptySurfaces,
	}
}

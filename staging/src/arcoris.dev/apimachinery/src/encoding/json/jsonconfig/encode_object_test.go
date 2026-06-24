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

import "testing"

func TestDefaultEncodeObjectConfig(t *testing.T) {
	t.Parallel()

	config := defaultEncodeObjectConfig()

	if config.TypeMeta != TypeMetaOmitZero {
		t.Fatalf("type meta = %d; want omit zero", config.TypeMeta)
	}
	if config.Metadata != MetadataOmitZero {
		t.Fatalf("metadata = %d; want omit zero", config.Metadata)
	}
	if config.Observed != ObservedOmitAbsent {
		t.Fatalf("observed = %d; want omit absent", config.Observed)
	}
}

func TestResolveEncodeObjectConfig(t *testing.T) {
	t.Parallel()

	config := EncodeObjectConfig{}
	resolveEncodeObjectConfig(&config)

	if config.TypeMeta == TypeMetaDefault {
		t.Fatalf("type meta still default")
	}
	if config.Metadata == MetadataDefault {
		t.Fatalf("metadata still default")
	}
	if config.Observed == ObservedDefault {
		t.Fatalf("observed still default")
	}
}

func TestValidateEncodeObjectConfigRejectsUnknownModes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config EncodeObjectConfig
		path   string
	}{
		"type meta": {
			config: EncodeObjectConfig{TypeMeta: TypeMetaEncodeMode(99), Metadata: MetadataOmitZero, Observed: ObservedOmitAbsent},
			path:   "encode.object.type_meta",
		},
		"metadata": {
			config: EncodeObjectConfig{TypeMeta: TypeMetaOmitZero, Metadata: MetadataEncodeMode(99), Observed: ObservedOmitAbsent},
			path:   "encode.object.metadata",
		},
		"observed": {
			config: EncodeObjectConfig{TypeMeta: TypeMetaOmitZero, Metadata: MetadataOmitZero, Observed: ObservedEncodeMode(99)},
			path:   "encode.object.observed",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateEncodeObjectConfig(testCase.config)
			requireConfigErrorIs(t, err, ErrInvalidConfig)
			requireErrorTextContains(t, err, testCase.path)
		})
	}
}

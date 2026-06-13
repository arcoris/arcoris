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

func TestDefaultEncodeConfig(t *testing.T) {
	t.Parallel()

	config := defaultEncodeConfig()

	if config.Output.Layout != LayoutCompact {
		t.Fatalf("layout = %d; want compact", config.Output.Layout)
	}
	if config.Ordering.Mode != OrderingPreserve {
		t.Fatalf("ordering = %d; want preserve", config.Ordering.Mode)
	}
	if config.Strings.InvalidUTF8 != InvalidUTF8Reject {
		t.Fatalf("invalid UTF-8 = %d; want reject", config.Strings.InvalidUTF8)
	}
	if config.Numbers.MaxDigits != DefaultEncodeMaxNumberDigits {
		t.Fatalf("max digits = %d; want %d", config.Numbers.MaxDigits, DefaultEncodeMaxNumberDigits)
	}
	if config.Values.Bytes != BytesEncodingReject {
		t.Fatalf("bytes = %d; want reject", config.Values.Bytes)
	}
	if config.Object.TypeMeta != TypeMetaOmitZero {
		t.Fatalf("type meta = %d; want omit zero", config.Object.TypeMeta)
	}
	if config.Ownership.EmptySurfaces != EmptyOwnershipSurfaceEmit {
		t.Fatalf("empty surfaces = %d; want emit", config.Ownership.EmptySurfaces)
	}
	if config.Limits.MaxDepth != DefaultMaxDepth {
		t.Fatalf("max depth = %d; want %d", config.Limits.MaxDepth, DefaultMaxDepth)
	}
}

func TestResolveEncodeConfig(t *testing.T) {
	t.Parallel()

	config := EncodeConfig{}
	resolveEncodeConfig(&config)

	if config.Output.Layout == LayoutDefault {
		t.Fatalf("layout still default")
	}
	if config.Ordering.Mode == OrderingDefault {
		t.Fatalf("ordering still default")
	}
	if config.Strings.InvalidUTF8 == InvalidUTF8Default {
		t.Fatalf("string invalid UTF-8 still default")
	}
	if config.Numbers.DecimalScale == DecimalScaleDefault {
		t.Fatalf("decimal scale still default")
	}
	if config.Values.Bytes == BytesEncodingDefault {
		t.Fatalf("bytes mode still default")
	}
	if config.Object.Metadata == MetadataDefault {
		t.Fatalf("metadata mode still default")
	}
	if config.Ownership.EmptySurfaces == EmptyOwnershipSurfaceDefault {
		t.Fatalf("empty surfaces mode still default")
	}
	if config.Limits.MaxDepth != DefaultMaxDepth {
		t.Fatalf("max depth = %d; want %d", config.Limits.MaxDepth, DefaultMaxDepth)
	}
}

func TestValidateEncodeConfig(t *testing.T) {
	t.Parallel()

	if err := validateEncodeConfig(defaultEncodeConfig()); err != nil {
		t.Fatalf("validateEncodeConfig() error = %v", err)
	}

	config := defaultEncodeConfig()
	config.Values.Bytes = BytesEncodingBase64Std

	err := validateEncodeConfig(config)
	requireConfigErrorIs(t, err, ErrUnsupportedConfig)
	requireErrorTextContains(t, err, "encode.values.bytes")
}

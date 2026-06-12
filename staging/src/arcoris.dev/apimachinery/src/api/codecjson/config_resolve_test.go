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

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/codecjson/jsonconfig"
)

func TestNewRejectsInvalidConfig(t *testing.T) {
	config := jsonconfig.Default()
	config.Decode.Limits.MaxDepth = -1

	_, err := New(config)
	if !errors.Is(err, jsonconfig.ErrInvalidConfig) {
		t.Fatalf("New() error = %v; want ErrInvalidConfig", err)
	}
}

func TestNewRejectsUnsupportedConfig(t *testing.T) {
	config := jsonconfig.Default()
	config.Encode.Values.Bytes = jsonconfig.BytesEncodingBase64Std

	_, err := New(config)
	if !errors.Is(err, jsonconfig.ErrUnsupportedConfig) {
		t.Fatalf("New() error = %v; want ErrUnsupportedConfig", err)
	}
}

func TestNewStoresResolvedConfig(t *testing.T) {
	config := jsonconfig.Config{}

	codec, err := New(config)
	requireNoError(t, err)

	if codec.decode.maxDepth != jsonconfig.DefaultMaxDepth {
		t.Fatalf("decode max depth = %d; want %d", codec.decode.maxDepth, jsonconfig.DefaultMaxDepth)
	}
	if codec.encode.maxNumberDigits != jsonconfig.DefaultEncodeMaxNumberDigits {
		t.Fatalf("encode max digits = %d; want %d", codec.encode.maxNumberDigits, jsonconfig.DefaultEncodeMaxNumberDigits)
	}
}

func TestResolveDecodeConfig(t *testing.T) {
	config := jsonconfig.Default()
	config.Decode.Limits.MaxDocumentBytes = 128
	config.Decode.Limits.MaxStringBytes = 16
	config.Decode.Objects.UnknownEnvelopeFields = jsonconfig.UnknownFieldIgnore
	config.Decode.Ownership.UnknownFields = jsonconfig.UnknownFieldIgnore
	config.Decode.Ownership.ValidateState = false

	resolved := resolveDecodeConfig(config.Decode)

	if resolved.maxDocumentBytes != 128 {
		t.Fatalf("max document bytes = %d; want 128", resolved.maxDocumentBytes)
	}
	if resolved.maxStringBytes != 16 {
		t.Fatalf("max string bytes = %d; want 16", resolved.maxStringBytes)
	}
	if resolved.rejectUnknownEnvelopeFields {
		t.Fatalf("reject unknown envelope fields = true; want false")
	}
	if resolved.rejectUnknownOwnershipFields {
		t.Fatalf("reject unknown ownership fields = true; want false")
	}
	if resolved.validateOwnershipState {
		t.Fatalf("validate ownership state = true; want false")
	}
}

func TestResolveEncodeConfig(t *testing.T) {
	config := jsonconfig.Default()
	config.Encode.Output.Layout = jsonconfig.LayoutPretty
	config.Encode.Output.FinalNewline = jsonconfig.FinalNewlineAppend
	config.Encode.Ordering.Mode = jsonconfig.OrderingDeterministic
	config.Encode.Strings.EscapeHTML = true
	config.Encode.Limits.MaxOutputBytes = 256

	resolved := resolveEncodeConfig(config.Encode)

	if !resolved.pretty {
		t.Fatalf("pretty = false; want true")
	}
	if !resolved.finalNewline {
		t.Fatalf("final newline = false; want true")
	}
	if !resolved.deterministic {
		t.Fatalf("deterministic = false; want true")
	}
	if !resolved.escapeHTML {
		t.Fatalf("escape HTML = false; want true")
	}
	if resolved.maxOutputBytes != 256 {
		t.Fatalf("max output bytes = %d; want 256", resolved.maxOutputBytes)
	}
}

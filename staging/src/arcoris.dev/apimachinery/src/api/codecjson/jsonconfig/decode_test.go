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

func TestDefaultDecodeConfig(t *testing.T) {
	t.Parallel()

	config := defaultDecodeConfig()

	if config.Limits.MaxDepth != DefaultMaxDepth {
		t.Fatalf("max depth = %d; want %d", config.Limits.MaxDepth, DefaultMaxDepth)
	}
	if config.Numbers.Mode != NumberModeExact {
		t.Fatalf("number mode = %d; want exact", config.Numbers.Mode)
	}
	if config.Strings.InvalidUTF8 != InvalidUTF8Reject {
		t.Fatalf("invalid UTF-8 mode = %d; want reject", config.Strings.InvalidUTF8)
	}
	if config.Objects.DuplicateKeys != DuplicateKeyReject {
		t.Fatalf("duplicate key mode = %d; want reject", config.Objects.DuplicateKeys)
	}
	if config.Ownership.UnknownFields != UnknownFieldReject {
		t.Fatalf("ownership unknown fields = %d; want reject", config.Ownership.UnknownFields)
	}
}

func TestResolveDecodeConfig(t *testing.T) {
	t.Parallel()

	config := DecodeConfig{}
	resolveDecodeConfig(&config)

	if config.Limits.MaxDepth != DefaultMaxDepth {
		t.Fatalf("max depth = %d; want %d", config.Limits.MaxDepth, DefaultMaxDepth)
	}
	if config.Numbers.Mode == NumberModeDefault {
		t.Fatalf("number mode still default")
	}
	if config.Strings.InvalidUTF8 == InvalidUTF8Default {
		t.Fatalf("string invalid UTF-8 mode still default")
	}
	if config.Objects.DuplicateKeys == DuplicateKeyDefault {
		t.Fatalf("duplicate key mode still default")
	}
	if config.Ownership.UnknownFields == UnknownFieldDefault {
		t.Fatalf("ownership unknown field mode still default")
	}
}

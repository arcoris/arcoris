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

func TestDefaultDecodeOwnershipConfig(t *testing.T) {
	t.Parallel()

	config := defaultDecodeOwnershipConfig()

	if config.UnknownFields != UnknownFieldReject {
		t.Fatalf("unknown fields = %d; want reject", config.UnknownFields)
	}
	if !config.ValidateDocument {
		t.Fatalf("validate document = false; want true")
	}
}

func TestResolveDecodeOwnershipConfig(t *testing.T) {
	t.Parallel()

	config := DecodeOwnershipConfig{}
	resolveDecodeOwnershipConfig(&config)

	if config.UnknownFields != UnknownFieldReject {
		t.Fatalf("unknown fields = %d; want reject", config.UnknownFields)
	}
	if !config.ValidateDocument {
		t.Fatalf("validate document = false; want true")
	}
}

func TestValidateDecodeOwnershipConfigRejectsUnknownMode(t *testing.T) {
	t.Parallel()

	err := validateDecodeOwnershipConfig(DecodeOwnershipConfig{UnknownFields: UnknownFieldMode(99)})
	requireConfigErrorIs(t, err, ErrInvalidConfig)
	requireErrorTextContains(t, err, "decode.ownership.unknown_fields")
}

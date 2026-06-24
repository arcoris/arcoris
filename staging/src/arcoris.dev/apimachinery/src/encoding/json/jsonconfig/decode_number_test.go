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

func TestResolveDecodeNumberConfig(t *testing.T) {
	t.Parallel()

	config := DecodeNumberConfig{}
	resolveDecodeNumberConfig(&config)

	if config.Mode != NumberModeExact {
		t.Fatalf("number mode = %d; want exact", config.Mode)
	}
}

func TestValidateDecodeNumberConfig(t *testing.T) {
	t.Parallel()

	if err := validateDecodeNumberConfig(DecodeNumberConfig{Mode: NumberModeExact}); err != nil {
		t.Fatalf("validateDecodeNumberConfig() error = %v", err)
	}

	err := validateDecodeNumberConfig(DecodeNumberConfig{Mode: NumberMode(99)})
	requireConfigErrorIs(t, err, ErrInvalidConfig)
	requireErrorTextContains(t, err, "decode.numbers.mode")
}

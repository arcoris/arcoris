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

func TestDefaultEncodeOrderingConfig(t *testing.T) {
	t.Parallel()

	config := defaultEncodeOrderingConfig()

	if config.Mode != OrderingPreserve {
		t.Fatalf("ordering mode = %d; want preserve", config.Mode)
	}
}

func TestResolveEncodeOrderingConfig(t *testing.T) {
	t.Parallel()

	config := EncodeOrderingConfig{}
	resolveEncodeOrderingConfig(&config)

	if config.Mode != OrderingPreserve {
		t.Fatalf("ordering mode = %d; want preserve", config.Mode)
	}
}

func TestValidateEncodeOrderingConfigRejectsUnknownMode(t *testing.T) {
	t.Parallel()

	err := validateEncodeOrderingConfig(EncodeOrderingConfig{Mode: OrderingMode(99)})
	requireConfigErrorIs(t, err, ErrInvalidConfig)
	requireErrorTextContains(t, err, "encode.ordering.mode")
}

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

func TestDefaultEncodeLimitsConfig(t *testing.T) {
	t.Parallel()

	limits := defaultEncodeLimitsConfig()

	if limits.MaxDepth != DefaultMaxDepth {
		t.Fatalf("max depth = %d; want %d", limits.MaxDepth, DefaultMaxDepth)
	}
	if limits.MaxOutputBytes != 0 {
		t.Fatalf("max output bytes = %d; want unlimited", limits.MaxOutputBytes)
	}
}

func TestResolveEncodeLimitsConfig(t *testing.T) {
	t.Parallel()

	limits := EncodeLimitsConfig{}
	resolveEncodeLimitsConfig(&limits)

	if limits.MaxDepth != DefaultMaxDepth {
		t.Fatalf("max depth = %d; want %d", limits.MaxDepth, DefaultMaxDepth)
	}
}

func TestValidateEncodeLimitsConfigRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		limits EncodeLimitsConfig
		path   string
	}{
		"max depth": {
			limits: EncodeLimitsConfig{MaxDepth: -1},
			path:   "encode.limits.max_depth",
		},
		"output bytes": {
			limits: EncodeLimitsConfig{MaxDepth: DefaultMaxDepth, MaxOutputBytes: -1},
			path:   "encode.limits.max_output_bytes",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateEncodeLimitsConfig(testCase.limits)
			requireConfigErrorIs(t, err, ErrInvalidConfig)
			requireErrorTextContains(t, err, testCase.path)
		})
	}
}

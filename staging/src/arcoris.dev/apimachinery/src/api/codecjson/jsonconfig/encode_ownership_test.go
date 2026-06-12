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

func TestDefaultEncodeOwnershipConfig(t *testing.T) {
	t.Parallel()

	config := defaultEncodeOwnershipConfig()

	if config.Normalize != OwnershipNormalizeWhenDeterministic {
		t.Fatalf("normalize = %d; want when deterministic", config.Normalize)
	}
	if config.EmptySurfaces != EmptyOwnershipSurfaceEmit {
		t.Fatalf("empty surfaces = %d; want emit", config.EmptySurfaces)
	}
	if config.EmptyEntries != EmptyEntriesEmit {
		t.Fatalf("empty entries = %d; want emit", config.EmptyEntries)
	}
}

func TestResolveEncodeOwnershipConfig(t *testing.T) {
	t.Parallel()

	config := EncodeOwnershipConfig{}
	resolveEncodeOwnershipConfig(&config)

	if config.Normalize == OwnershipNormalizeDefault {
		t.Fatalf("normalize still default")
	}
	if config.EmptySurfaces == EmptyOwnershipSurfaceDefault {
		t.Fatalf("empty surfaces still default")
	}
	if config.EmptyEntries == EmptyEntriesDefault {
		t.Fatalf("empty entries still default")
	}
}

func TestValidateEncodeOwnershipConfigRejectsUnknownModes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config EncodeOwnershipConfig
		path   string
	}{
		"normalize": {
			config: EncodeOwnershipConfig{Normalize: OwnershipNormalizeMode(99), EmptySurfaces: EmptyOwnershipSurfaceEmit, EmptyEntries: EmptyEntriesEmit},
			path:   "encode.ownership.normalize",
		},
		"empty surfaces": {
			config: EncodeOwnershipConfig{Normalize: OwnershipNormalizeWhenDeterministic, EmptySurfaces: EmptyOwnershipSurfaceMode(99), EmptyEntries: EmptyEntriesEmit},
			path:   "encode.ownership.empty_surfaces",
		},
		"empty entries": {
			config: EncodeOwnershipConfig{Normalize: OwnershipNormalizeWhenDeterministic, EmptySurfaces: EmptyOwnershipSurfaceEmit, EmptyEntries: EmptyEntriesMode(99)},
			path:   "encode.ownership.empty_entries",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateEncodeOwnershipConfig(testCase.config)
			requireConfigErrorIs(t, err, ErrInvalidConfig)
			requireErrorTextContains(t, err, testCase.path)
		})
	}
}

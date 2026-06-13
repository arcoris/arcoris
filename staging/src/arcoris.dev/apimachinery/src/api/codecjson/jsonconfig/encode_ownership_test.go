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

	if config.EmptySurfaces != EmptyOwnershipSurfaceEmit {
		t.Fatalf("empty surfaces = %d; want emit", config.EmptySurfaces)
	}
}

func TestResolveEncodeOwnershipConfig(t *testing.T) {
	t.Parallel()

	config := EncodeOwnershipConfig{}
	resolveEncodeOwnershipConfig(&config)

	if config.EmptySurfaces == EmptyOwnershipSurfaceDefault {
		t.Fatalf("empty surfaces still default")
	}
}

func TestValidateEncodeOwnershipConfigRejectsUnknownMode(t *testing.T) {
	t.Parallel()

	err := validateEncodeOwnershipConfig(EncodeOwnershipConfig{EmptySurfaces: EmptyOwnershipSurfaceMode(99)})
	requireConfigErrorIs(t, err, ErrInvalidConfig)
	requireErrorTextContains(t, err, "encode.ownership.empty_surfaces")
}

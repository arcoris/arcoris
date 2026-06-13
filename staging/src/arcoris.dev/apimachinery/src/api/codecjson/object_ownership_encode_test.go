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
	"testing"

	"arcoris.dev/apimachinery/api/codecjson/jsonconfig"
	"arcoris.dev/apimachinery/api/objectownership"
)

func TestEncodeObjectOwnershipState(t *testing.T) {
	state := objectownership.NewState(ownershipSurface(
		ownershipEntry("user-cli", "$.image", "$.replicas"),
	))

	data, err := newTestCodec(t).EncodeObjectOwnership(state)
	requireNoError(t, err)

	want := `{"desired":{"entries":[{"owner":"user-cli","fields":["$.image","$.replicas"]}]},"observed":{"entries":[]},"metadata":{"labels":{"entries":[]},"annotations":{"entries":[]}}}`
	if string(data) != want {
		t.Fatalf("encoded = %s; want %s", data, want)
	}
}

func TestEncodeObjectOwnershipUsesCanonicalStateOrder(t *testing.T) {
	state := objectownership.NewState(ownershipSurface(
		ownershipEntry("user-b", "$.b"),
		ownershipEntry("user-a", "$.b", "$.a"),
	))

	data, err := newTestCodec(t).EncodeObjectOwnership(state)
	requireNoError(t, err)

	want := `{"desired":{"entries":[{"owner":"user-a","fields":["$.a","$.b"]},{"owner":"user-b","fields":["$.b"]}]},"observed":{"entries":[]},"metadata":{"labels":{"entries":[]},"annotations":{"entries":[]}}}`
	if string(data) != want {
		t.Fatalf("encoded = %s; want %s", data, want)
	}
}

func TestEncodeObjectOwnershipPretty(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Output.Layout = jsonconfig.LayoutPretty
	})

	data, err := c.EncodeObjectOwnership(objectownership.EmptyState())
	requireNoError(t, err)

	want := "{\n  \"desired\": {\n    \"entries\": []\n  },\n  \"observed\": {\n    \"entries\": []\n  },\n  \"metadata\": {\n    \"labels\": {\n      \"entries\": []\n    },\n    \"annotations\": {\n      \"entries\": []\n    }\n  }\n}"
	if string(data) != want {
		t.Fatalf("encoded = %q; want %q", data, want)
	}
}

func TestEncodeObjectOwnershipOmitsEmptySurfacesWhenConfigured(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Ownership.EmptySurfaces = jsonconfig.EmptyOwnershipSurfaceOmit
	})

	data, err := c.EncodeObjectOwnership(objectownership.EmptyState())
	requireNoError(t, err)

	if string(data) != `{}` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeObjectOwnershipOmitEmptySurfacesKeepsNonEmptyDesired(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Ownership.EmptySurfaces = jsonconfig.EmptyOwnershipSurfaceOmit
	})
	state := objectownership.NewState(ownershipSurface(ownershipEntry("user-cli", "$.image")))

	data, err := c.EncodeObjectOwnership(state)
	requireNoError(t, err)

	if string(data) != `{"desired":{"entries":[{"owner":"user-cli","fields":["$.image"]}]}}` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeObjectOwnershipAllSurfaces(t *testing.T) {
	state := ownershipState(
		ownershipSurface(ownershipEntry("user-cli", "$.image")),
		ownershipSurface(ownershipEntry("controller", "$.ready")),
		objectownership.NewMetadataState(
			ownershipSurface(ownershipEntry("labeler", `$["scheduler.arcoris.dev/mode"]`)),
			ownershipSurface(ownershipEntry("annotator", `$["with.dots"]`)),
		),
	)

	data, err := newTestCodec(t).EncodeObjectOwnership(state)
	requireNoError(t, err)

	want := `{"desired":{"entries":[{"owner":"user-cli","fields":["$.image"]}]},"observed":{"entries":[{"owner":"controller","fields":["$.ready"]}]},"metadata":{"labels":{"entries":[{"owner":"labeler","fields":["$[\"scheduler.arcoris.dev/mode\"]"]}]},"annotations":{"entries":[{"owner":"annotator","fields":["$[\"with.dots\"]"]}]}}}`
	if string(data) != want {
		t.Fatalf("encoded = %s; want %s", data, want)
	}
}

func TestEncodeObjectOwnershipOmitEmptySurfacesWithLabelsOnly(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Ownership.EmptySurfaces = jsonconfig.EmptyOwnershipSurfaceOmit
	})
	state := ownershipState(
		ownershipSurface(),
		ownershipSurface(),
		objectownership.NewMetadataState(
			ownershipSurface(ownershipEntry("labeler", `$["app"]`)),
			ownershipSurface(),
		),
	)

	data, err := c.EncodeObjectOwnership(state)
	requireNoError(t, err)

	want := `{"metadata":{"labels":{"entries":[{"owner":"labeler","fields":["$[\"app\"]"]}]}}}`
	if string(data) != want {
		t.Fatalf("encoded = %s; want %s", data, want)
	}
}

func TestEncodeObjectOwnershipOmitEmptySurfacesWithAnnotationsOnly(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Ownership.EmptySurfaces = jsonconfig.EmptyOwnershipSurfaceOmit
	})
	state := ownershipState(
		ownershipSurface(),
		ownershipSurface(),
		objectownership.NewMetadataState(
			ownershipSurface(),
			ownershipSurface(ownershipEntry("annotator", `$["note"]`)),
		),
	)

	data, err := c.EncodeObjectOwnership(state)
	requireNoError(t, err)

	want := `{"metadata":{"annotations":{"entries":[{"owner":"annotator","fields":["$[\"note\"]"]}]}}}`
	if string(data) != want {
		t.Fatalf("encoded = %s; want %s", data, want)
	}
}

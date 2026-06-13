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
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestDecodeObjectOwnershipState(t *testing.T) {
	data := []byte(`{"desired":{"entries":[{"owner":"user-cli","fields":["$.image","$.replicas"]}]},"observed":{"entries":[{"owner":"controller","fields":["$.ready"]}]},"metadata":{"labels":{"entries":[{"owner":"labeler","fields":["$[\"scheduler.arcoris.dev/mode\"]"]}]},"annotations":{"entries":[{"owner":"annotator","fields":["$[\"with.dots\"]"]}]}}}`)

	got, err := newTestCodec(t).DecodeObjectOwnership(data)
	requireNoError(t, err)

	requireOwnershipFields(t, got.Desired(), "user-cli", "$.image", "$.replicas")
	requireOwnershipFields(t, got.Observed(), "controller", "$.ready")
	requireOwnershipFields(t, got.Metadata().Labels(), "labeler", `$["scheduler.arcoris.dev/mode"]`)
	requireOwnershipFields(t, got.Metadata().Annotations(), "annotator", `$["with.dots"]`)
}

func TestDecodeObjectOwnershipCollapsesDuplicateEntries(t *testing.T) {
	got, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{"entries":[{"owner":"user-cli","fields":["$.image"]},{"owner":"user-cli","fields":["$.image","$.replicas"]}]}}`))
	requireNoError(t, err)

	requireOwnershipFields(t, got.Desired(), "user-cli", "$.image", "$.replicas")
}

func TestDecodeObjectOwnershipEmptySurfaces(t *testing.T) {
	got, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{"entries":[]}}`))
	requireNoError(t, err)

	if !got.Desired().IsEmpty() {
		t.Fatalf("desired = %#v; want empty", got.Desired())
	}
}

func TestDecodeObjectOwnershipEmptyState(t *testing.T) {
	got, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{}`))
	requireNoError(t, err)

	if !got.IsEmpty() {
		t.Fatalf("state = %#v; want empty", got)
	}
}

func TestDecodeObjectOwnershipRejectsUnknownMetadataField(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"metadata":{"finalizers":{}}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$.metadata.finalizers", ErrorReasonInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsInvalidSurfaceShapes(t *testing.T) {
	testCases := map[string]struct {
		json string
		path string
	}{
		"desired entries non-array": {
			json: `{"desired":{"entries":{}}}`,
			path: "$.desired.entries",
		},
		"observed entries non-array": {
			json: `{"observed":{"entries":{}}}`,
			path: "$.observed.entries",
		},
		"metadata labels entries non-array": {
			json: `{"metadata":{"labels":{"entries":{}}}}`,
			path: "$.metadata.labels.entries",
		},
		"metadata annotations entries non-array": {
			json: `{"metadata":{"annotations":{"entries":{}}}}`,
			path: "$.metadata.annotations.entries",
		},
		"desired entry non-object": {
			json: `{"desired":{"entries":[1]}}`,
			path: "$.desired.entries[0]",
		},
		"observed missing owner": {
			json: `{"observed":{"entries":[{"fields":[]}]}}`,
			path: "$.observed.entries[0].owner",
		},
		"metadata labels owner non-string": {
			json: `{"metadata":{"labels":{"entries":[{"owner":1}]}}}`,
			path: "$.metadata.labels.entries[0].owner",
		},
		"metadata annotations field non-string": {
			json: `{"metadata":{"annotations":{"entries":[{"owner":"annotator","fields":[1]}]}}}`,
			path: "$.metadata.annotations.entries[0].fields[0]",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := newTestCodec(t).DecodeObjectOwnership([]byte(testCase.json))

			requireErrorIs(t, err, ErrInvalidEnvelope)
			requireCodecJSONError(t, err, testCase.path, ErrorReasonInvalidEnvelope)
		})
	}
}

func TestDecodeObjectOwnershipRejectsUnknownField(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"unknown":true}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$.unknown", ErrorReasonInvalidEnvelope)
}

func TestDecodeObjectOwnershipIgnoresUnknownFieldsWhenConfigured(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Decode.Ownership.UnknownFields = jsonconfig.UnknownFieldIgnore
	})

	state, err := c.DecodeObjectOwnership([]byte(`{"extra":1,"metadata":{"extra":true},"desired":{"entries":[{"owner":"user-cli","fields":["$.image"],"extra":1}]}}`))
	requireNoError(t, err)

	requireOwnershipFields(t, state.Desired(), "user-cli", "$.image")
}

func TestDecodeObjectOwnershipRejectsInvalidOwner(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{"entries":[{"owner":"","fields":[]}]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireErrorIs(t, err, fieldownership.ErrInvalidOwner)
}

func TestDecodeObjectOwnershipRejectsInvalidPath(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{"entries":[{"owner":"user-cli","fields":["not-a-path"]}]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireErrorIs(t, err, fieldpath.ErrInvalidPath)
}

func TestDecodeObjectOwnershipRejectsDuplicateKey(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{},"desired":{}}`))

	requireErrorIs(t, err, ErrDuplicateKey)
}

func requireOwnershipFields(t *testing.T, state fieldownership.State, ownerName string, paths ...string) {
	t.Helper()

	got := state.FieldsFor(fieldownership.MustOwner(ownerName))
	want := ownershipFields(paths...)
	if !got.Equal(want) {
		t.Fatalf("%s fields = %s; want %s", ownerName, got.String(), want.String())
	}
}

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

func TestDecodeObjectOwnershipNormalizesDuplicateEntries(t *testing.T) {
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

func TestDecodeObjectOwnershipRejectsEntriesNonArray(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{"entries":{}}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsEntryNonObject(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{"entries":[1]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsOwnerMissing(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{"entries":[{"fields":[]}]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$.desired.entries[0].owner", ErrorReasonInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsOwnerNonString(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{"entries":[{"owner":1}]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsFieldNonString(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{"entries":[{"owner":"user-cli","fields":[1]}]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
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

	state, err := c.DecodeObjectOwnership([]byte(`{"extra":1}`))
	requireNoError(t, err)

	if !state.IsEmpty() {
		t.Fatalf("state = %#v; want empty", state)
	}
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

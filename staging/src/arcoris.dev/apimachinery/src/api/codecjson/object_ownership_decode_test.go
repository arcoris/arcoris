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
	"arcoris.dev/apimachinery/api/objectownership"
)

func TestDecodeObjectOwnershipDocument(t *testing.T) {
	data := []byte(`{"version":"v1","desired":{"entries":[{"owner":"user-cli","fields":["$.image","$.replicas"]}]}}`)

	got, err := newTestCodec(t).DecodeObjectOwnership(data)
	requireNoError(t, err)

	if got.Version != objectownership.VersionV1 {
		t.Fatalf("version = %q", got.Version)
	}
	if len(got.Desired.Entries) != 1 {
		t.Fatalf("entries = %#v", got.Desired.Entries)
	}
	if got.Desired.Entries[0].Owner != fieldownership.Owner("user-cli") {
		t.Fatalf("owner = %q", got.Desired.Entries[0].Owner)
	}
}

func TestDecodeObjectOwnershipEmptyDesired(t *testing.T) {
	got, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","desired":{"entries":[]}}`))
	requireNoError(t, err)

	if !got.Desired.IsEmpty() {
		t.Fatalf("desired = %#v; want empty", got.Desired)
	}
}

func TestDecodeObjectOwnershipMissingDesired(t *testing.T) {
	got, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1"}`))
	requireNoError(t, err)

	if !got.Desired.IsEmpty() {
		t.Fatalf("desired = %#v; want empty", got.Desired)
	}
}

func TestDecodeObjectOwnershipRejectsMissingVersion(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"desired":{}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$.version", ErrorReasonInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsVersionNonString(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":1}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsEntriesNonArray(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","desired":{"entries":{}}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsEntryNonObject(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","desired":{"entries":[1]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsOwnerMissing(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","desired":{"entries":[{"fields":[]}]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$.desired.entries[0].owner", ErrorReasonInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsOwnerNonString(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","desired":{"entries":[{"owner":1}]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsFieldNonString(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","desired":{"entries":[{"owner":"user-cli","fields":[1]}]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
}

func TestDecodeObjectOwnershipRejectsUnknownField(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","unknown":true}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$.unknown", ErrorReasonInvalidEnvelope)
}

func TestDecodeObjectOwnershipIgnoresUnknownFieldsWhenConfigured(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Decode.Ownership.UnknownFields = jsonconfig.UnknownFieldIgnore
	})

	doc, err := c.DecodeObjectOwnership([]byte(`{"version":"v1","extra":1}`))
	requireNoError(t, err)

	if doc.Version != objectownership.VersionV1 {
		t.Fatalf("version = %q; want v1", doc.Version)
	}
}

func TestDecodeObjectOwnershipRejectsInvalidOwner(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","desired":{"entries":[{"owner":"","fields":[]}]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireErrorIs(t, err, objectownership.ErrInvalidEntry)
}

func TestDecodeObjectOwnershipRejectsInvalidPath(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","desired":{"entries":[{"owner":"user-cli","fields":["not-a-path"]}]}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireErrorIs(t, err, objectownership.ErrInvalidPath)
}

func TestDecodeObjectOwnershipRejectsDuplicateKey(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","version":"v1"}`))

	requireErrorIs(t, err, ErrDuplicateKey)
}

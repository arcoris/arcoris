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

func TestEncodeObjectOwnershipDocument(t *testing.T) {
	doc := objectownership.Document{
		Version: objectownership.VersionV1,
		Desired: objectownership.Surface{Entries: []objectownership.Entry{
			{
				Owner:  fieldownership.Owner("user-cli"),
				Fields: []objectownership.Path{"$.image", "$.replicas"},
			},
		}},
	}

	data, err := newTestCodec(t).EncodeObjectOwnership(doc)
	requireNoError(t, err)

	want := `{"version":"v1","desired":{"entries":[{"owner":"user-cli","fields":["$.image","$.replicas"]}]}}`
	if string(data) != want {
		t.Fatalf("encoded = %s; want %s", data, want)
	}
}

func TestEncodeObjectOwnershipValidatesBeforeEncode(t *testing.T) {
	doc := objectownership.Document{Version: objectownership.VersionV1}
	doc.Desired.Entries = []objectownership.Entry{{Owner: fieldownership.Owner(""), Fields: []objectownership.Path{"$.image"}}}

	_, err := newTestCodec(t).EncodeObjectOwnership(doc)

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireErrorIs(t, err, objectownership.ErrInvalidEntry)
}

func TestEncodeObjectOwnershipDeterministicNormalizes(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Ordering.Mode = jsonconfig.OrderingDeterministic
	})
	doc := objectownership.Document{
		Version: objectownership.VersionV1,
		Desired: objectownership.Surface{Entries: []objectownership.Entry{
			{Owner: fieldownership.Owner("user-b"), Fields: []objectownership.Path{"$.b"}},
			{Owner: fieldownership.Owner("user-a"), Fields: []objectownership.Path{"$.b", "$.a"}},
		}},
	}

	data, err := c.EncodeObjectOwnership(doc)
	requireNoError(t, err)

	want := `{"version":"v1","desired":{"entries":[{"owner":"user-a","fields":["$.a","$.b"]},{"owner":"user-b","fields":["$.b"]}]}}`
	if string(data) != want {
		t.Fatalf("encoded = %s; want %s", data, want)
	}
}

func TestEncodeObjectOwnershipPreservesOrderByDefault(t *testing.T) {
	doc := objectownership.Document{
		Version: objectownership.VersionV1,
		Desired: objectownership.Surface{Entries: []objectownership.Entry{
			{Owner: fieldownership.Owner("user-b"), Fields: []objectownership.Path{"$.b"}},
			{Owner: fieldownership.Owner("user-a"), Fields: []objectownership.Path{"$.b", "$.a"}},
		}},
	}

	data, err := newTestCodec(t).EncodeObjectOwnership(doc)
	requireNoError(t, err)

	want := `{"version":"v1","desired":{"entries":[{"owner":"user-b","fields":["$.b"]},{"owner":"user-a","fields":["$.b","$.a"]}]}}`
	if string(data) != want {
		t.Fatalf("encoded = %s; want %s", data, want)
	}
}

func TestEncodeObjectOwnershipPretty(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Output.Layout = jsonconfig.LayoutPretty
	})
	doc := objectownership.Document{Version: objectownership.VersionV1}

	data, err := c.EncodeObjectOwnership(doc)
	requireNoError(t, err)

	want := "{\n  \"version\": \"v1\",\n  \"desired\": {\n    \"entries\": []\n  }\n}"
	if string(data) != want {
		t.Fatalf("encoded = %q; want %q", data, want)
	}
}

func TestEncodeObjectOwnershipOmitsEmptyDesiredWhenConfigured(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Ownership.EmptyDesired = jsonconfig.EmptyOwnershipSurfaceOmit
	})

	data, err := c.EncodeObjectOwnership(objectownership.Document{Version: objectownership.VersionV1})
	requireNoError(t, err)

	if string(data) != `{"version":"v1"}` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeObjectOwnershipOmitsEmptyEntriesWhenConfigured(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Ownership.EmptyEntries = jsonconfig.EmptyEntriesOmit
	})

	data, err := c.EncodeObjectOwnership(objectownership.Document{Version: objectownership.VersionV1})
	requireNoError(t, err)

	if string(data) != `{"version":"v1","desired":{}}` {
		t.Fatalf("encoded = %s", data)
	}
}

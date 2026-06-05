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

package codecregistry

import (
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

func TestNewAllowsRepeatedMediaTypeWithDifferentEntryIDs(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.storage", codec.FormatJSON, codec.MediaTypeJSON),
	)

	entries := registry.EntriesByMediaType(codec.MediaTypeJSON)
	if len(entries) != 2 {
		t.Fatalf("EntriesByMediaType(json) length = %d; want 2", len(entries))
	}
}

func TestNewAllowsDuplicateFormat(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.storage", codec.FormatJSON, "application/vnd.arcoris.storage+json"),
	)

	entries := registry.EntriesByFormat(codec.FormatJSON)
	if len(entries) != 2 {
		t.Fatalf("EntriesByFormat(json) length = %d; want 2", len(entries))
	}
}

func TestNewRejectsDuplicateEntryID(t *testing.T) {
	_, err := New(
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.public", codec.FormatJSON, "application/vnd.arcoris.storage+json"),
	)

	requireErrorIs(t, err, ErrDuplicateEntryID)
	requireRegistryError(t, err, "registrations[1].id", ErrorReasonDuplicateEntryID)
	requireRegistryDetailContains(t, err, "registrations[0]")
}

func TestRepeatedMediaTypeNoLongerErrors(t *testing.T) {
	_, err := New(
		testObjectByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testObjectByteRegistration("json.storage", codec.FormatJSON, codec.MediaTypeJSON),
	)

	requireNoError(t, err)
}

func TestDuplicateCodecInstanceAllowedWithDifferentEntryIDs(t *testing.T) {
	c := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)

	registry := testRegistry(
		t,
		testRegistration("json.public", c),
		testRegistration("json.storage", c),
	)

	if registry.Len() != 2 {
		t.Fatalf("Len() = %d; want 2", registry.Len())
	}
}

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

package codecselection

import (
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

func TestNormalizeDecodeBindingAt(t *testing.T) {
	registry := testRegistry(
		t,
		testObjectByteRegistration("json.public", codec.MediaTypeJSON),
	)

	record, err := normalizeDecodeBindingAt(
		"codecselection.decodeBindings[0]",
		Config{Registry: registry},
		DecodeBinding{
			ContentType: testContentType(" APPLICATION/JSON "),
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     "json.public",
		},
	)
	requireNoError(t, err)

	if record.contentType.MediaType() != codec.MediaTypeJSON {
		t.Fatalf("media type = %q; want %q", record.contentType.MediaType(), codec.MediaTypeJSON)
	}
	if record.entryID != codecregistry.MustEntryID("json.public") {
		t.Fatalf("EntryID = %q; want json.public", record.entryID)
	}
	if record.entry.IsZero() {
		t.Fatalf("entry is zero; want registry entry")
	}
}

func TestNormalizeEncodeBindingAt(t *testing.T) {
	registry := testRegistry(
		t,
		testObjectByteRegistration("json.public", codec.MediaTypeJSON),
	)

	record, err := normalizeEncodeBindingAt(
		"codecselection.encodeBindings[0]",
		Config{Registry: registry},
		EncodeBinding{
			ContentType: testContentType(codec.MediaTypeJSON),
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.public"),
		},
	)
	requireNoError(t, err)

	if record.target != codec.TargetObject {
		t.Fatalf("target = %q; want %q", record.target, codec.TargetObject)
	}
	if record.transport != TransportBytes {
		t.Fatalf("transport = %q; want %q", record.transport, TransportBytes)
	}
}

func TestNormalizeBindingAtWrapsInvalidContentTypeAsInvalidBinding(t *testing.T) {
	_, err := normalizeDecodeBindingAt(
		"codecselection.decodeBindings[0]",
		Config{},
		DecodeBinding{
			ContentType: ContentType{mediaType: "application/json; charset=utf-8"},
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.public"),
		},
	)

	requireErrorIs(t, err, ErrInvalidBinding)
	requireErrorIs(t, err, ErrInvalidContentType)
	requireSelectionError(t, err, "codecselection.decodeBindings[0].contentType", ErrorReasonInvalidBinding)
}

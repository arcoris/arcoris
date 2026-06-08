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

func TestBuildEncodeBindings(t *testing.T) {
	registry := testRegistry(
		t,
		testObjectByteRegistration("json.public", codec.MediaTypeJSON),
	)
	contentType := testContentType(codec.MediaTypeJSON)

	records, bindings, err := buildEncodeBindings(Config{
		Registry: registry,
		EncodeBindings: []EncodeBinding{{
			ContentType: contentType,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.public"),
		}},
	})
	requireNoError(t, err)

	if len(records) != 1 {
		t.Fatalf("records length = %d; want 1", len(records))
	}
	if len(bindings) != 1 {
		t.Fatalf("bindings length = %d; want 1", len(bindings))
	}
	if bindings[0].EntryID != codecregistry.MustEntryID("json.public") {
		t.Fatalf("binding EntryID = %q; want json.public", bindings[0].EntryID)
	}
}

func TestBuildEncodeBindingsRejectsDuplicateKey(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.MediaTypeJSON),
		testFullByteRegistration("json.storage", codec.MediaTypeJSON),
	)

	_, _, err := buildEncodeBindings(Config{
		Registry: registry,
		EncodeBindings: []EncodeBinding{
			{
				ContentType: testContentType(codec.MediaTypeJSON),
				Target:      codec.TargetObject,
				Transport:   TransportBytes,
				EntryID:     codecregistry.MustEntryID("json.public"),
			},
			{
				ContentType: testContentType(" APPLICATION/JSON "),
				Target:      codec.TargetObject,
				Transport:   TransportBytes,
				EntryID:     codecregistry.MustEntryID("json.storage"),
			},
		},
	})

	requireErrorIs(t, err, ErrDuplicateEncodeBinding)
	requireSelectionError(t, err, "codecselection.encodeBindings[1]", ErrorReasonDuplicateEncodeBinding)
}

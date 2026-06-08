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

func TestBuildDecodeBindings(t *testing.T) {
	registry := testRegistry(
		t,
		testObjectByteRegistration("json.public", codec.MediaTypeJSON),
	)
	contentType := testContentType(codec.MediaTypeJSON)

	records, bindings, err := buildDecodeBindings(Config{
		Registry: registry,
		DecodeBindings: []DecodeBinding{{
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
	if bindings[0].ContentType.String() != contentType.String() {
		t.Fatalf("binding content type = %q; want %q", bindings[0].ContentType, contentType)
	}
}

func TestBuildDecodeBindingsRejectsDuplicateKey(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.MediaTypeJSON),
		testFullByteRegistration("json.storage", codec.MediaTypeJSON),
	)

	_, _, err := buildDecodeBindings(Config{
		Registry: registry,
		DecodeBindings: []DecodeBinding{
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

	requireErrorIs(t, err, ErrDuplicateDecodeBinding)
	requireSelectionError(t, err, "codecselection.decodeBindings[1]", ErrorReasonDuplicateDecodeBinding)
}

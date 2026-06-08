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

func TestApplicationJSONAloneIsNotRegistryIdentity(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.MediaTypeJSON),
		testFullByteRegistration("json.storage", codec.MediaTypeJSON),
		testFullByteRegistration("json.canonical", codec.MediaTypeJSON),
	)
	contentType := testContentType(codec.MediaTypeJSON)
	selector := testSelector(t, Config{
		Registry: registry,
		DecodeBindings: []DecodeBinding{{
			ContentType: contentType,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.storage"),
		}},
	})

	selection, _, err := selector.SelectObjectDecoder(contentType)
	requireNoError(t, err)

	if selection.EntryID != codecregistry.MustEntryID("json.storage") {
		t.Fatalf("EntryID = %q; want explicit json.storage binding", selection.EntryID)
	}
}

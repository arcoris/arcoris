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

func TestSelectionIsZero(t *testing.T) {
	var selection Selection

	if !selection.IsZero() {
		t.Fatalf("zero selection IsZero() = false; want true")
	}
}

func TestSelectionIsZeroFalseForSelectedEntry(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.MediaTypeJSON),
	)
	entry, ok := registry.LookupID(codecregistry.MustEntryID("json.public"))
	if !ok {
		t.Fatalf("json.public entry missing")
	}

	selection := Selection{
		Direction:   DirectionDecode,
		Transport:   TransportBytes,
		Target:      codec.TargetObject,
		ContentType: testContentType(codec.MediaTypeJSON),
		EntryID:     codecregistry.MustEntryID("json.public"),
		Entry:       entry,
	}

	if selection.IsZero() {
		t.Fatalf("selected selection IsZero() = true; want false")
	}
}

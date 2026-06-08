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

func TestDecodeBindingFromRecord(t *testing.T) {
	record := bindingRecord{
		contentType: testContentType(codec.MediaTypeJSON),
		target:      codec.TargetObject,
		transport:   TransportBytes,
		entryID:     codecregistry.MustEntryID("json.public"),
	}

	binding := decodeBindingFromRecord(record)

	if binding.ContentType.String() != record.contentType.String() ||
		binding.Target != record.target ||
		binding.Transport != record.transport ||
		binding.EntryID != record.entryID {
		t.Fatalf("binding = %#v; want record projection", binding)
	}
}

func TestEncodeBindingFromRecord(t *testing.T) {
	record := bindingRecord{
		contentType: testContentType(codec.MediaTypeJSON),
		target:      codec.TargetObject,
		transport:   TransportBytes,
		entryID:     codecregistry.MustEntryID("json.public"),
	}

	binding := encodeBindingFromRecord(record)

	if binding.ContentType.String() != record.contentType.String() ||
		binding.Target != record.target ||
		binding.Transport != record.transport ||
		binding.EntryID != record.entryID {
		t.Fatalf("binding = %#v; want record projection", binding)
	}
}

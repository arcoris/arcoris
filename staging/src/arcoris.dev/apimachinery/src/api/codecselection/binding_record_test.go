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

func TestBindingRecordKey(t *testing.T) {
	record := bindingRecord{
		contentType: testContentType(codec.MediaTypeJSON, MustParameter("profile", "public")),
		target:      codec.TargetObject,
		transport:   TransportStream,
		entryID:     codecregistry.MustEntryID("json.public"),
	}

	key := record.key()

	if key.contentType != record.contentType.key() {
		t.Fatalf("contentType key = %q; want %q", key.contentType, record.contentType.key())
	}
	if key.target != codec.TargetObject {
		t.Fatalf("target = %q; want %q", key.target, codec.TargetObject)
	}
	if key.transport != TransportStream {
		t.Fatalf("transport = %q; want %q", key.transport, TransportStream)
	}
}

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

func TestEntrySupportsCapability(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.bytes", codec.MediaTypeJSON),
		testValueStreamRegistration("json.stream", codec.MediaTypeCBOR),
	)
	byteEntry, ok := registry.LookupID(codecregistry.MustEntryID("json.bytes"))
	if !ok {
		t.Fatalf("json.bytes entry missing")
	}
	streamEntry, ok := registry.LookupID(codecregistry.MustEntryID("json.stream"))
	if !ok {
		t.Fatalf("json.stream entry missing")
	}

	if !entrySupportsCapability(byteEntry, codec.TargetValue, TransportBytes) {
		t.Fatalf("byte entry does not support byte value capability")
	}
	if entrySupportsCapability(byteEntry, codec.TargetValue, TransportStream) {
		t.Fatalf("byte entry supports stream value capability")
	}
	if !entrySupportsCapability(streamEntry, codec.TargetValue, TransportStream) {
		t.Fatalf("stream entry does not support stream value capability")
	}
	if entrySupportsCapability(streamEntry, codec.TargetValue, TransportBytes) {
		t.Fatalf("stream entry supports byte value capability")
	}
}

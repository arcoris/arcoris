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
		testFullByteRegistration("json.bytes", codec.MediaTypeJSON),
		testFullStreamRegistration("json.stream", codec.MediaTypeJSON),
		testValueByteRegistration("json.value", codec.MediaTypeJSON),
	)
	byteEntry, ok := registry.LookupID(codecregistry.MustEntryID("json.bytes"))
	if !ok {
		t.Fatalf("json.bytes entry missing")
	}
	streamEntry, ok := registry.LookupID(codecregistry.MustEntryID("json.stream"))
	if !ok {
		t.Fatalf("json.stream entry missing")
	}
	valueEntry, ok := registry.LookupID(codecregistry.MustEntryID("json.value"))
	if !ok {
		t.Fatalf("json.value entry missing")
	}

	tests := []struct {
		name      string
		entry     codecregistry.Entry
		target    codec.Target
		transport Transport
		want      bool
	}{
		{
			name:      "value bytes",
			entry:     byteEntry,
			target:    codec.TargetValue,
			transport: TransportBytes,
			want:      true,
		},
		{
			name:      "object bytes",
			entry:     byteEntry,
			target:    codec.TargetObject,
			transport: TransportBytes,
			want:      true,
		},
		{
			name:      "ownership bytes",
			entry:     byteEntry,
			target:    codec.TargetObjectOwnership,
			transport: TransportBytes,
			want:      true,
		},
		{
			name:      "byte entry does not support stream",
			entry:     byteEntry,
			target:    codec.TargetValue,
			transport: TransportStream,
			want:      false,
		},
		{
			name:      "value stream",
			entry:     streamEntry,
			target:    codec.TargetValue,
			transport: TransportStream,
			want:      true,
		},
		{
			name:      "object stream",
			entry:     streamEntry,
			target:    codec.TargetObject,
			transport: TransportStream,
			want:      true,
		},
		{
			name:      "ownership stream",
			entry:     streamEntry,
			target:    codec.TargetObjectOwnership,
			transport: TransportStream,
			want:      true,
		},
		{
			name:      "stream entry does not support bytes",
			entry:     streamEntry,
			target:    codec.TargetValue,
			transport: TransportBytes,
			want:      false,
		},
		{
			name:      "target-specific entry rejects other target",
			entry:     valueEntry,
			target:    codec.TargetObject,
			transport: TransportBytes,
			want:      false,
		},
		{
			name:      "unknown target",
			entry:     byteEntry,
			target:    codec.Target("unknown"),
			transport: TransportBytes,
			want:      false,
		},
		{
			name:      "unknown transport",
			entry:     byteEntry,
			target:    codec.TargetValue,
			transport: Transport(99),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entrySupportsCapability(tt.entry, tt.target, tt.transport)
			if got != tt.want {
				t.Fatalf("entrySupportsCapability(...) = %v; want %v", got, tt.want)
			}
		})
	}
}

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
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

func TestSelectorBindingAccessorsReturnDetachedSlices(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.MediaTypeJSON),
	)
	contentType := testContentType(codec.MediaTypeJSON)
	selector := testSelector(t, Config{
		Registry: registry,
		DecodeBindings: []DecodeBinding{{
			ContentType: contentType,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.public"),
		}},
		EncodeBindings: []EncodeBinding{{
			ContentType: contentType,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.public"),
		}},
	})

	decodeBindings := selector.DecodeBindings()
	encodeBindings := selector.EncodeBindings()
	decodeBindings[0].EntryID = codecregistry.MustEntryID("json.mutated")
	encodeBindings[0].EntryID = codecregistry.MustEntryID("json.mutated")

	if got := selector.DecodeBindings()[0].EntryID; got != codecregistry.MustEntryID("json.public") {
		t.Fatalf("DecodeBindings returned mutable source, got %q", got)
	}
	if got := selector.EncodeBindings()[0].EntryID; got != codecregistry.MustEntryID("json.public") {
		t.Fatalf("EncodeBindings returned mutable source, got %q", got)
	}
}

func TestSelectorConcurrentRuntimeSelection(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.MediaTypeJSON),
	)
	contentType := testContentType(codec.MediaTypeJSON)
	selector := testSelector(t, Config{
		Registry: registry,
		DecodeBindings: []DecodeBinding{{
			ContentType: contentType,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.public"),
		}},
		EncodeBindings: []EncodeBinding{{
			ContentType: contentType,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.public"),
		}},
	})
	preferences := testPreferenceSet(testPreference(contentType, int(WeightDefault)))

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 128; j++ {
				decodeSelection, _, err := selector.SelectObjectDecoder(contentType)
				requireNoError(t, err)
				if decodeSelection.EntryID != codecregistry.MustEntryID("json.public") {
					t.Errorf("decode EntryID = %q; want json.public", decodeSelection.EntryID)
				}

				encodeSelection, _, err := selector.SelectObjectEncoder(preferences)
				requireNoError(t, err)
				if encodeSelection.EntryID != codecregistry.MustEntryID("json.public") {
					t.Errorf("encode EntryID = %q; want json.public", encodeSelection.EntryID)
				}
			}
		}()
	}
	wg.Wait()
}

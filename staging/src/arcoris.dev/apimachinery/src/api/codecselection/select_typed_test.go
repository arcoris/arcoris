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

func TestSelectTypedPublicMethods(t *testing.T) {
	const (
		byteEntryID   = "json.bytes"
		streamEntryID = "json.stream"
	)

	contentType := testContentType(codec.MediaTypeJSON)
	byteID := codecregistry.MustEntryID(byteEntryID)
	streamID := codecregistry.MustEntryID(streamEntryID)
	selector := testSelector(t, Config{
		Registry: testRegistry(
			t,
			testFullByteRegistration(byteEntryID, codec.MediaTypeJSON),
			testFullStreamRegistration(streamEntryID, codec.MediaTypeJSON),
		),
		DecodeBindings: []DecodeBinding{
			{ContentType: contentType, Target: codec.TargetValue, Transport: TransportBytes, EntryID: byteID},
			{ContentType: contentType, Target: codec.TargetValue, Transport: TransportStream, EntryID: streamID},
			{ContentType: contentType, Target: codec.TargetObject, Transport: TransportBytes, EntryID: byteID},
			{ContentType: contentType, Target: codec.TargetObject, Transport: TransportStream, EntryID: streamID},
			{ContentType: contentType, Target: codec.TargetObjectOwnership, Transport: TransportBytes, EntryID: byteID},
			{ContentType: contentType, Target: codec.TargetObjectOwnership, Transport: TransportStream, EntryID: streamID},
		},
		EncodeBindings: []EncodeBinding{
			{ContentType: contentType, Target: codec.TargetValue, Transport: TransportBytes, EntryID: byteID},
			{ContentType: contentType, Target: codec.TargetValue, Transport: TransportStream, EntryID: streamID},
			{ContentType: contentType, Target: codec.TargetObject, Transport: TransportBytes, EntryID: byteID},
			{ContentType: contentType, Target: codec.TargetObject, Transport: TransportStream, EntryID: streamID},
			{ContentType: contentType, Target: codec.TargetObjectOwnership, Transport: TransportBytes, EntryID: byteID},
			{ContentType: contentType, Target: codec.TargetObjectOwnership, Transport: TransportStream, EntryID: streamID},
		},
	})
	preferences := testPreferenceSet(testPreference(contentType, int(WeightDefault)))

	tests := []struct {
		name      string
		direction Direction
		target    codec.Target
		transport Transport
		entryID   codecregistry.EntryID
		run       func() (Selection, any, error)
	}{
		{
			name:      "value decoder",
			direction: DirectionDecode,
			target:    codec.TargetValue,
			transport: TransportBytes,
			entryID:   byteID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectValueDecoder(contentType)
				return selection, selected, err
			},
		},
		{
			name:      "value encoder",
			direction: DirectionEncode,
			target:    codec.TargetValue,
			transport: TransportBytes,
			entryID:   byteID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectValueEncoder(preferences)
				return selection, selected, err
			},
		},
		{
			name:      "value stream decoder",
			direction: DirectionDecode,
			target:    codec.TargetValue,
			transport: TransportStream,
			entryID:   streamID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectValueStreamDecoder(contentType)
				return selection, selected, err
			},
		},
		{
			name:      "value stream encoder",
			direction: DirectionEncode,
			target:    codec.TargetValue,
			transport: TransportStream,
			entryID:   streamID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectValueStreamEncoder(preferences)
				return selection, selected, err
			},
		},
		{
			name:      "object decoder",
			direction: DirectionDecode,
			target:    codec.TargetObject,
			transport: TransportBytes,
			entryID:   byteID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectObjectDecoder(contentType)
				return selection, selected, err
			},
		},
		{
			name:      "object encoder",
			direction: DirectionEncode,
			target:    codec.TargetObject,
			transport: TransportBytes,
			entryID:   byteID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectObjectEncoder(preferences)
				return selection, selected, err
			},
		},
		{
			name:      "object stream decoder",
			direction: DirectionDecode,
			target:    codec.TargetObject,
			transport: TransportStream,
			entryID:   streamID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectObjectStreamDecoder(contentType)
				return selection, selected, err
			},
		},
		{
			name:      "object stream encoder",
			direction: DirectionEncode,
			target:    codec.TargetObject,
			transport: TransportStream,
			entryID:   streamID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectObjectStreamEncoder(preferences)
				return selection, selected, err
			},
		},
		{
			name:      "ownership decoder",
			direction: DirectionDecode,
			target:    codec.TargetObjectOwnership,
			transport: TransportBytes,
			entryID:   byteID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectObjectOwnershipDecoder(contentType)
				return selection, selected, err
			},
		},
		{
			name:      "ownership encoder",
			direction: DirectionEncode,
			target:    codec.TargetObjectOwnership,
			transport: TransportBytes,
			entryID:   byteID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectObjectOwnershipEncoder(preferences)
				return selection, selected, err
			},
		},
		{
			name:      "ownership stream decoder",
			direction: DirectionDecode,
			target:    codec.TargetObjectOwnership,
			transport: TransportStream,
			entryID:   streamID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectObjectOwnershipStreamDecoder(contentType)
				return selection, selected, err
			},
		},
		{
			name:      "ownership stream encoder",
			direction: DirectionEncode,
			target:    codec.TargetObjectOwnership,
			transport: TransportStream,
			entryID:   streamID,
			run: func() (Selection, any, error) {
				selection, selected, err := selector.SelectObjectOwnershipStreamEncoder(preferences)
				return selection, selected, err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selection, selected, err := tt.run()
			requireNoError(t, err)

			if selected == nil {
				t.Fatalf("selected codec = nil")
			}
			if selection.Direction != tt.direction {
				t.Fatalf("Direction = %q; want %q", selection.Direction, tt.direction)
			}
			if selection.Transport != tt.transport {
				t.Fatalf("Transport = %q; want %q", selection.Transport, tt.transport)
			}
			if selection.Target != tt.target {
				t.Fatalf("Target = %q; want %q", selection.Target, tt.target)
			}
			if !selection.ContentType.Equal(contentType) {
				t.Fatalf("ContentType = %q; want %q", selection.ContentType, contentType)
			}
			if selection.EntryID != tt.entryID {
				t.Fatalf("EntryID = %q; want %q", selection.EntryID, tt.entryID)
			}
		})
	}
}

func TestTypedSelectionRejectsMissingRuntimeCapability(t *testing.T) {
	selection := Selection{
		Direction: DirectionDecode,
		Transport: TransportBytes,
		Target:    codec.TargetValue,
		EntryID:   codecregistry.MustEntryID("json.public"),
	}

	_, _, err := typedSelection[codec.ValueCodec](selection, "codecselection.decode.value")

	requireErrorIs(t, err, ErrEntryCapabilityMismatch)
	requireSelectionError(t, err, "codecselection.decode.value", ErrorReasonEntryCapabilityMismatch)
}

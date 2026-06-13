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

package codecjson

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/objectownership"
)

func TestDecodeObjectOwnershipFromMatchesDecodeObjectOwnership(t *testing.T) {
	c := newTestCodec(t)
	testCases := map[string]string{
		"empty":        `{}`,
		"all surfaces": `{"desired":{"entries":[{"owner":"user-cli","fields":["$.image"]}]},"observed":{"entries":[{"owner":"controller","fields":["$.ready"]}]},"metadata":{"labels":{"entries":[{"owner":"labeler","fields":["$[\"app\"]"]}]},"annotations":{"entries":[{"owner":"annotator","fields":["$[\"note\"]"]}]}}}`,
	}

	for name, data := range testCases {
		t.Run(name, func(t *testing.T) {
			fromBytes, err := c.DecodeObjectOwnership([]byte(data))
			requireNoError(t, err)
			fromStream, err := c.DecodeObjectOwnershipFrom(strings.NewReader(data))
			requireNoError(t, err)

			if !reflect.DeepEqual(fromStream, fromBytes) {
				t.Fatalf("stream state = %#v; bytes state = %#v", fromStream, fromBytes)
			}
		})
	}
}

func TestEncodeObjectOwnershipToMatchesEncodeObjectOwnership(t *testing.T) {
	c := newTestCodec(t)
	testCases := map[string]objectownership.State{
		"empty": objectownership.EmptyState(),
		"all surfaces": ownershipState(
			ownershipSurface(ownershipEntry("user-cli", "$.image")),
			ownershipSurface(ownershipEntry("controller", "$.ready")),
			objectownership.NewMetadataState(
				ownershipSurface(ownershipEntry("labeler", `$["app"]`)),
				ownershipSurface(ownershipEntry("annotator", `$["note"]`)),
			),
		),
	}

	for name, state := range testCases {
		t.Run(name, func(t *testing.T) {
			fromBytes, err := c.EncodeObjectOwnership(state)
			requireNoError(t, err)
			var buffer bytes.Buffer
			requireNoError(t, c.EncodeObjectOwnershipTo(&buffer, state))

			if buffer.String() != string(fromBytes) {
				t.Fatalf("stream = %s; bytes = %s", buffer.String(), fromBytes)
			}
		})
	}
}

func TestDecodeObjectOwnershipFromErrorSentinelsMatchBytes(t *testing.T) {
	c := newTestCodec(t)
	testCases := map[string]struct {
		data    string
		targets []error
	}{
		"unknown field": {
			data:    `{"unknown":true}`,
			targets: []error{ErrInvalidEnvelope},
		},
		"invalid path": {
			data:    `{"desired":{"entries":[{"owner":"user-cli","fields":["not-a-path"]}]}}`,
			targets: []error{ErrInvalidEnvelope, fieldpath.ErrInvalidPath},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			_, bytesErr := c.DecodeObjectOwnership([]byte(testCase.data))
			_, streamErr := c.DecodeObjectOwnershipFrom(strings.NewReader(testCase.data))

			for _, target := range testCase.targets {
				requireErrorIs(t, bytesErr, target)
				requireErrorIs(t, streamErr, target)
			}
		})
	}
}

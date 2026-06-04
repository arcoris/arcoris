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

package codecregistry

import (
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

func TestNewAllowsDuplicateFormatWithDifferentMediaTypes(t *testing.T) {
	registry, err := New(
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeYAML),
	)
	requireNoError(t, err)

	entries := registry.EntriesByFormat(codec.FormatJSON)
	if len(entries) != 2 {
		t.Fatalf("EntriesByFormat(json) length = %d; want 2", len(entries))
	}
}

func TestNewRejectsDuplicateMediaType(t *testing.T) {
	_, err := New(
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
		newValueByteCodec(codec.FormatYAML, codec.MediaTypeJSON),
	)

	requireErrorIs(t, err, ErrDuplicateMediaType)
	requireRegistryError(t, err, "codecs[1].info.mediaTypes[0]", ErrorReasonDuplicateMediaType)
}

func TestDuplicateMediaTypeErrorUsesNormalizedMediaTypeIndex(t *testing.T) {
	c := &fakeValueByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: codec.Info{
			Format: codec.FormatYAML,
			MediaTypes: []codec.MediaType{
				"application/aaa",
				codec.MediaTypeJSON,
			},
			Targets: []codec.Target{codec.TargetValue},
		},
	}}

	_, err := New(
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
		c,
	)

	requireErrorIs(t, err, ErrDuplicateMediaType)
	requireRegistryError(t, err, "codecs[1].info.mediaTypes[1]", ErrorReasonDuplicateMediaType)
}

func TestDuplicateMediaTypeError(t *testing.T) {
	_, err := New(
		newObjectByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
		newObjectByteCodec(codec.FormatYAML, codec.MediaTypeJSON),
	)

	requireErrorIs(t, err, ErrDuplicateMediaType)
}

func TestDuplicateCodecInstanceRejected(t *testing.T) {
	c := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)

	_, err := New(c, c)

	requireErrorIs(t, err, ErrDuplicateMediaType)
}

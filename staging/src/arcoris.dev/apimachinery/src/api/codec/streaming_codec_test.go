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

package codec

import (
	"io"

	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/value"
)

// StreamingCodec is the full streaming aggregate and intentionally does not
// require byte-slice Codec support.
var _ StreamingCodec = fakeStreamingCodec{}

type fakeStreamingCodec struct {
	fakeBaseCodec
}

func (fakeStreamingCodec) DecodeValueFrom(io.Reader, DecodeOptions) (value.Value, error) {
	return value.NullValue(), nil
}

func (fakeStreamingCodec) EncodeValueTo(io.Writer, value.Value, EncodeOptions) error {
	return nil
}

func (fakeStreamingCodec) DecodeObjectFrom(io.Reader, DecodeOptions) (Object, error) {
	return Object{}, nil
}

func (fakeStreamingCodec) EncodeObjectTo(io.Writer, Object, EncodeOptions) error {
	return nil
}

func (fakeStreamingCodec) DecodeObjectOwnershipFrom(io.Reader, DecodeOptions) (objectownership.Document, error) {
	return objectownership.Document{}, nil
}

func (fakeStreamingCodec) EncodeObjectOwnershipTo(io.Writer, objectownership.Document, EncodeOptions) error {
	return nil
}

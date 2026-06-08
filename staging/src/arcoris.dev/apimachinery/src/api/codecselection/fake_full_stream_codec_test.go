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
	"io"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/value"
)

func testFullStreamRegistration(id string, mediaTypes ...codec.MediaType) codecregistry.Registration {
	return testRegistration(id, newFullStreamCodec(codec.FormatJSON, mediaTypes...))
}

type fakeFullStreamCodec struct {
	fakeBaseCodec
}

func newFullStreamCodec(format codec.Format, mediaTypes ...codec.MediaType) *fakeFullStreamCodec {
	return &fakeFullStreamCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(
			format,
			mediaTypes,
			codec.TargetValue,
			codec.TargetObject,
			codec.TargetObjectOwnership,
		),
	}}
}

func (f fakeFullStreamCodec) DecodeValueFrom(r io.Reader) (value.Value, error) {
	return fakeValueStreamCodec{}.DecodeValueFrom(r)
}

func (f fakeFullStreamCodec) EncodeValueTo(w io.Writer, v value.Value) error {
	return fakeValueStreamCodec{}.EncodeValueTo(w, v)
}

func (f fakeFullStreamCodec) DecodeObjectFrom(r io.Reader) (codec.Object, error) {
	return fakeObjectStreamCodec{}.DecodeObjectFrom(r)
}

func (f fakeFullStreamCodec) EncodeObjectTo(w io.Writer, obj codec.Object) error {
	return fakeObjectStreamCodec{}.EncodeObjectTo(w, obj)
}

func (f fakeFullStreamCodec) DecodeObjectOwnershipFrom(
	r io.Reader,
) (objectownership.Document, error) {
	return fakeOwnershipStreamCodec{}.DecodeObjectOwnershipFrom(r)
}

func (f fakeFullStreamCodec) EncodeObjectOwnershipTo(
	w io.Writer,
	doc objectownership.Document,
) error {
	return fakeOwnershipStreamCodec{}.EncodeObjectOwnershipTo(w, doc)
}

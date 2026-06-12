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
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/value"
)

func testFullByteRegistration(id string, mediaTypes ...codec.MediaType) codecregistry.Registration {
	return testRegistration(id, newFullByteCodec(codec.FormatJSON, mediaTypes...))
}

type fakeFullByteCodec struct {
	fakeBaseCodec
}

func newFullByteCodec(format codec.Format, mediaTypes ...codec.MediaType) *fakeFullByteCodec {
	return &fakeFullByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(
			format,
			mediaTypes,
			codec.TargetValue,
			codec.TargetObject,
			codec.TargetObjectOwnership,
		),
	}}
}

func (f fakeFullByteCodec) DecodeValue(data []byte) (value.Value, error) {
	return fakeValueByteCodec{}.DecodeValue(data)
}

func (f fakeFullByteCodec) EncodeValue(v value.Value) ([]byte, error) {
	return fakeValueByteCodec{}.EncodeValue(v)
}

func (f fakeFullByteCodec) DecodeObject(data []byte) (codec.Object, error) {
	return fakeObjectByteCodec{}.DecodeObject(data)
}

func (f fakeFullByteCodec) EncodeObject(obj codec.Object) ([]byte, error) {
	return fakeObjectByteCodec{}.EncodeObject(obj)
}

func (f fakeFullByteCodec) DecodeObjectOwnership(
	data []byte,
) (objectownership.State, error) {
	return fakeOwnershipByteCodec{}.DecodeObjectOwnership(data)
}

func (f fakeFullByteCodec) EncodeObjectOwnership(
	state objectownership.State,
) ([]byte, error) {
	return fakeOwnershipByteCodec{}.EncodeObjectOwnership(state)
}

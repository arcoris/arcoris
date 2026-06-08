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
)

func testObjectByteRegistration(id string, mediaTypes ...codec.MediaType) codecregistry.Registration {
	return testRegistration(id, newObjectByteCodec(codec.FormatJSON, mediaTypes...))
}

type fakeObjectByteCodec struct {
	fakeBaseCodec
}

func newObjectByteCodec(format codec.Format, mediaTypes ...codec.MediaType) *fakeObjectByteCodec {
	return &fakeObjectByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaTypes, codec.TargetObject),
	}}
}

func (f fakeObjectByteCodec) DecodeObject([]byte) (codec.Object, error) {
	return codec.Object{}, nil
}

func (f fakeObjectByteCodec) EncodeObject(codec.Object) ([]byte, error) {
	return nil, nil
}

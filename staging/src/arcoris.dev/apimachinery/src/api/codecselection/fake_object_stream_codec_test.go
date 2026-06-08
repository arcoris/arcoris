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
)

func testObjectStreamRegistration(id string, mediaTypes ...codec.MediaType) codecregistry.Registration {
	return testRegistration(id, newObjectStreamCodec(codec.FormatJSON, mediaTypes...))
}

type fakeObjectStreamCodec struct {
	fakeBaseCodec
}

func newObjectStreamCodec(format codec.Format, mediaTypes ...codec.MediaType) *fakeObjectStreamCodec {
	return &fakeObjectStreamCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaTypes, codec.TargetObject),
	}}
}

func (f fakeObjectStreamCodec) DecodeObjectFrom(io.Reader) (codec.Object, error) {
	return codec.Object{}, nil
}

func (f fakeObjectStreamCodec) EncodeObjectTo(io.Writer, codec.Object) error {
	return nil
}

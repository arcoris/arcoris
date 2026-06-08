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
	"arcoris.dev/apimachinery/api/value"
)

func testValueStreamRegistration(id string, mediaTypes ...codec.MediaType) codecregistry.Registration {
	return testRegistration(id, newValueStreamCodec(codec.FormatJSON, mediaTypes...))
}

type fakeValueStreamCodec struct {
	fakeBaseCodec
}

func newValueStreamCodec(format codec.Format, mediaTypes ...codec.MediaType) *fakeValueStreamCodec {
	return &fakeValueStreamCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaTypes, codec.TargetValue),
	}}
}

func (f fakeValueStreamCodec) DecodeValueFrom(io.Reader) (value.Value, error) {
	return value.NullValue(), nil
}

func (f fakeValueStreamCodec) EncodeValueTo(io.Writer, value.Value) error {
	return nil
}

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
	"arcoris.dev/apimachinery/api/value"
)

func testValueByteRegistration(id string, mediaTypes ...codec.MediaType) codecregistry.Registration {
	return testRegistration(id, newValueByteCodec(codec.FormatJSON, mediaTypes...))
}

type fakeValueByteCodec struct {
	fakeBaseCodec
}

func newValueByteCodec(format codec.Format, mediaTypes ...codec.MediaType) *fakeValueByteCodec {
	return &fakeValueByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaTypes, codec.TargetValue),
	}}
}

func (f fakeValueByteCodec) DecodeValue([]byte) (value.Value, error) {
	return value.NullValue(), nil
}

func (f fakeValueByteCodec) EncodeValue(value.Value) ([]byte, error) {
	return nil, nil
}

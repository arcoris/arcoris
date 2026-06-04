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

import "arcoris.dev/apimachinery/api/value"

var _ ValueDecoder = fakeValueCodec{}
var _ ValueEncoder = fakeValueCodec{}

// ValueCodec is a partial byte capability for the value target.
var _ ValueCodec = fakeValueCodec{}
var _ ValueCodec = fakeByteCodec{}

type fakeValueCodec struct {
	fakeBaseCodec
}

func (fakeValueCodec) DecodeValue([]byte) (value.Value, error) {
	return value.NullValue(), nil
}

func (fakeValueCodec) EncodeValue(value.Value) ([]byte, error) {
	return nil, nil
}

func (fakeByteCodec) DecodeValue(data []byte) (value.Value, error) {
	return fakeValueCodec{}.DecodeValue(data)
}

func (fakeByteCodec) EncodeValue(v value.Value) ([]byte, error) {
	return fakeValueCodec{}.EncodeValue(v)
}

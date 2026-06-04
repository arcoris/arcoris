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

	"arcoris.dev/apimachinery/api/value"
)

var _ ValueStreamDecoder = fakeValueStreamCodec{}
var _ ValueStreamEncoder = fakeValueStreamCodec{}

// ValueStreamCodec is a partial streaming capability for the value target.
var _ ValueStreamCodec = fakeValueStreamCodec{}

type fakeValueStreamCodec struct {
	fakeBaseCodec
}

func (fakeValueStreamCodec) DecodeValueFrom(io.Reader) (value.Value, error) {
	return value.NullValue(), nil
}

func (fakeValueStreamCodec) EncodeValueTo(io.Writer, value.Value) error {
	return nil
}

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
	"testing"

	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/value"
)

var _ ObjectDecoder = fakeObjectCodec{}
var _ ObjectEncoder = fakeObjectCodec{}
var _ ObjectCodec = fakeObjectCodec{}
var _ ObjectCodec = fakeFullCodec{}

type fakeObjectCodec struct {
	fakeBaseCodec
}

func (fakeObjectCodec) DecodeObject([]byte, DecodeOptions) (Object, error) {
	return Object{}, nil
}

func (fakeObjectCodec) EncodeObject(Object, EncodeOptions) ([]byte, error) {
	return nil, nil
}

func (fakeFullCodec) DecodeObject(data []byte, opts DecodeOptions) (Object, error) {
	return fakeObjectCodec{}.DecodeObject(data, opts)
}

func (fakeFullCodec) EncodeObject(obj Object, opts EncodeOptions) ([]byte, error) {
	return fakeObjectCodec{}.EncodeObject(obj, opts)
}

func TestObjectAlias(t *testing.T) {
	var got Object
	var _ object.Object[value.Value, value.Value] = got
}

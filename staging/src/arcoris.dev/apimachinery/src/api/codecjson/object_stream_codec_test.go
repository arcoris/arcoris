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

package codecjson

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/value"
)

func TestDecodeObjectFromMatchesDecodeObject(t *testing.T) {
	c := newTestCodec(t)
	data := `{"desired":{"a":1}}`

	fromBytes, err := c.DecodeObject([]byte(data))
	requireNoError(t, err)
	fromStream, err := c.DecodeObjectFrom(strings.NewReader(data))
	requireNoError(t, err)

	if !reflect.DeepEqual(fromStream, fromBytes) {
		t.Fatalf("stream object = %#v; bytes object = %#v", fromStream, fromBytes)
	}
}

func TestEncodeObjectToMatchesEncodeObject(t *testing.T) {
	c := newTestCodec(t)
	obj := codec.Object{Desired: value.NullValue()}

	fromBytes, err := c.EncodeObject(obj)
	requireNoError(t, err)
	var buffer bytes.Buffer
	requireNoError(t, c.EncodeObjectTo(&buffer, obj))

	if buffer.String() != string(fromBytes) {
		t.Fatalf("stream = %s; bytes = %s", buffer.String(), fromBytes)
	}
}

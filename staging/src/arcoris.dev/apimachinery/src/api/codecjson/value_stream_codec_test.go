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
	"arcoris.dev/apimachinery/api/value"
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeValueFromMatchesDecodeValue(t *testing.T) {
	c := newTestCodec(t)

	fromBytes, err := c.DecodeValue([]byte(`{"a":[1]}`))
	requireNoError(t, err)
	fromStream, err := c.DecodeValueFrom(strings.NewReader(`{"a":[1]}`))
	requireNoError(t, err)

	if !reflect.DeepEqual(fromStream, fromBytes) {
		t.Fatalf("stream value = %#v; bytes value = %#v", fromStream, fromBytes)
	}
}

func TestEncodeValueToMatchesEncodeValue(t *testing.T) {
	c := newTestCodec(t)
	v := value.MustRecordValue(value.MustRecordMember("a", value.Int64Value(1)))

	fromBytes, err := c.EncodeValue(v)
	requireNoError(t, err)
	var buffer bytes.Buffer
	requireNoError(t, c.EncodeValueTo(&buffer, v))

	if buffer.String() != string(fromBytes) {
		t.Fatalf("stream = %s; bytes = %s", buffer.String(), fromBytes)
	}
}

func TestDecodeValueFromRejectsTrailingData(t *testing.T) {
	_, err := newTestCodec(t).DecodeValueFrom(strings.NewReader(`null true`))

	requireErrorIs(t, err, ErrTrailingData)
}

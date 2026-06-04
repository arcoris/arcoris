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
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecjson/jsonconfig"
	"arcoris.dev/apimachinery/api/value"
)

func TestDecodeValueNull(t *testing.T) {
	got, err := newTestCodec(t).DecodeValue([]byte(`null`))
	requireNoError(t, err)

	requireKind(t, got, value.KindNull)
}

func TestDecodeValueBool(t *testing.T) {
	got, err := newTestCodec(t).DecodeValue([]byte(`true`))
	requireNoError(t, err)

	requireKind(t, got, value.KindBool)
	payload, _ := got.Bool()
	if !payload {
		t.Fatalf("bool = false; want true")
	}
}

func TestDecodeValueString(t *testing.T) {
	got, err := newTestCodec(t).DecodeValue([]byte(`"hello"`))
	requireNoError(t, err)

	requireStringValue(t, got, "hello")
}

func TestDecodeValueIntegerInt64(t *testing.T) {
	got, err := newTestCodec(t).DecodeValue([]byte(`-42`))
	requireNoError(t, err)

	requireIntegerText(t, got, "-42")
}

func TestDecodeValueIntegerUint64(t *testing.T) {
	got, err := newTestCodec(t).DecodeValue([]byte(`9223372036854775808`))
	requireNoError(t, err)

	requireIntegerText(t, got, "9223372036854775808")
}

func TestDecodeValueRejectsIntegerBeyondUint64(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte(`18446744073709551616`))

	requireErrorIs(t, err, ErrInvalidNumber)
	requireErrorIs(t, err, codec.ErrInvalidNumber)
}

func TestDecodeValueDecimal(t *testing.T) {
	got, err := newTestCodec(t).DecodeValue([]byte(`1.250`))
	requireNoError(t, err)

	requireDecimalText(t, got, "1.250")
}

func TestDecodeValueExponentDecimal(t *testing.T) {
	got, err := newTestCodec(t).DecodeValue([]byte(`1e3`))
	requireNoError(t, err)

	requireKind(t, got, value.KindDecimal)
	requireDecimalText(t, got, "1000")
}

func TestDecodeValueArrayPreservesOrder(t *testing.T) {
	got, err := newTestCodec(t).DecodeValue([]byte(`[1,"two",null]`))
	requireNoError(t, err)

	requireKind(t, got, value.KindList)
	list, _ := got.List()
	items := list.Items()
	requireIntegerText(t, items[0], "1")
	requireStringValue(t, items[1], "two")
	requireKind(t, items[2], value.KindNull)
}

func TestDecodeValueObjectPreservesOrder(t *testing.T) {
	got, err := newTestCodec(t).DecodeValue([]byte(`{"b":1,"a":2}`))
	requireNoError(t, err)

	requireObjectMemberNames(t, got, "b", "a")
}

func TestDecodeValueRejectsDuplicateKey(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte(`{"a":1,"a":2}`))

	requireErrorIs(t, err, ErrDuplicateKey)
}

func TestDecodeValueRejectsTrailingData(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte(`{} []`))

	requireErrorIs(t, err, ErrTrailingData)
}

func TestDecodeValueRejectsMalformedJSON(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte(`[`))

	requireErrorIs(t, err, ErrInvalidJSON)
	requireErrorIs(t, err, codec.ErrDecodeFailed)
}

func TestDecodeValueRejectsInvalidUTF8(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte{'"', 0xff, '"'})

	requireErrorIs(t, err, ErrInvalidJSON)
	requireErrorIs(t, err, codec.ErrDecodeFailed)
}

func TestDecodeValueRejectsMaxDepthExceeded(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Decode.Limits.MaxDepth = 2
	})

	_, err := c.DecodeValue([]byte(`[[null]]`))

	requireErrorIs(t, err, codec.ErrDepthExceeded)
}

func TestDecodeValueReportsJSONPath(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Decode.Limits.MaxDepth = 3
	})

	_, err := c.DecodeValue([]byte(`{"items":[[null]]}`))

	requireCodecJSONError(t, err, "$.items[0][0]", ErrorReasonMaxDepthExceeded)
}

func TestDecodeValueUsesConfiguredMaxDocumentBytes(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Decode.Limits.MaxDocumentBytes = 3
	})

	_, err := c.DecodeValue([]byte(`true`))

	requireErrorIs(t, err, ErrInvalidJSON)
}

func TestDecodeValueUsesConfiguredMaxStringBytes(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Decode.Limits.MaxStringBytes = 2
	})

	_, err := c.DecodeValue([]byte(`"abc"`))

	requireErrorIs(t, err, ErrInvalidJSON)
}

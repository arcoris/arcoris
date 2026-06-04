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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecjson/jsonconfig"
	"arcoris.dev/apimachinery/api/value"
)

func TestInvalidJSONErrorWrapping(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte(`[`))

	requireErrorIs(t, err, ErrInvalidJSON)
	requireErrorIs(t, err, codec.ErrDecodeFailed)
}

func TestDuplicateKeyErrorWrapping(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte(`{"a":1,"a":2}`))

	requireErrorIs(t, err, ErrDuplicateKey)
	requireErrorIs(t, err, codec.ErrInvalidDocument)
}

func TestTrailingDataErrorWrapping(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte(`null null`))

	requireErrorIs(t, err, ErrTrailingData)
	requireErrorIs(t, err, codec.ErrDecodeFailed)
}

func TestInvalidNumberErrorWrapping(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte(`18446744073709551616`))

	requireErrorIs(t, err, ErrInvalidNumber)
	requireErrorIs(t, err, codec.ErrInvalidNumber)
}

func TestUnsupportedValueErrorWrapping(t *testing.T) {
	_, err := newTestCodec(t).EncodeValue(value.BytesValue([]byte("abc")))

	requireErrorIs(t, err, ErrUnsupportedValue)
	requireErrorIs(t, err, codec.ErrEncodeFailed)
	requireErrorIs(t, err, codec.ErrUnsupportedFeature)
}

func TestMaxDepthErrorWrapping(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Decode.Limits.MaxDepth = 2
	})

	_, err := c.DecodeValue([]byte(`[[null]]`))

	requireErrorIs(t, err, ErrInvalidJSON)
	requireErrorIs(t, err, codec.ErrDepthExceeded)
}

func TestErrorAsCodecJSONErrorIfLocalTypeExists(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte(`[`))

	var codecJSONErr *Error
	if !errors.As(err, &codecJSONErr) {
		t.Fatalf("errors.As(*Error) = false")
	}
}

func TestErrorDiagnosticPath(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Decode.Limits.MaxDepth = 3
	})

	_, err := c.DecodeValue([]byte(`{"items":[[null]]}`))

	requireCodecJSONError(t, err, "$.items[0][0]", ErrorReasonMaxDepthExceeded)
}

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
	"math"
	"testing"
	"time"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecjson/jsonconfig"
	"arcoris.dev/apimachinery/api/value"
)

func TestEncodeValueNull(t *testing.T) {
	data, err := newTestCodec(t).EncodeValue(value.NullValue())
	requireNoError(t, err)

	if string(data) != `null` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeValueBool(t *testing.T) {
	data, err := newTestCodec(t).EncodeValue(value.BoolValue(true))
	requireNoError(t, err)

	if string(data) != `true` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeValueString(t *testing.T) {
	data, err := newTestCodec(t).EncodeValue(value.StringValue("hello"))
	requireNoError(t, err)

	if string(data) != `"hello"` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeValueInteger(t *testing.T) {
	data, err := newTestCodec(t).EncodeValue(value.Uint64Value(9223372036854775808))
	requireNoError(t, err)

	if string(data) != `9223372036854775808` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeValueDecimal(t *testing.T) {
	data, err := newTestCodec(t).EncodeValue(mustDecimalValue(t, "1.250"))
	requireNoError(t, err)

	if string(data) != `1.250` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeValueFloatFinite(t *testing.T) {
	data, err := newTestCodec(t).EncodeValue(mustFloatValue(t, 1.5))
	requireNoError(t, err)

	if string(data) != `1.5` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeValueList(t *testing.T) {
	list, err := value.ListValue(value.Int64Value(1), value.StringValue("two"))
	requireNoError(t, err)

	data, err := newTestCodec(t).EncodeValue(list)
	requireNoError(t, err)

	if string(data) != `[1,"two"]` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeValueObjectPreservesOrderByDefault(t *testing.T) {
	object := value.MustRecordValue(
		value.MustRecordMember("b", value.Int64Value(1)),
		value.MustRecordMember("a", value.Int64Value(2)),
	)

	data, err := newTestCodec(t).EncodeValue(object)
	requireNoError(t, err)

	if string(data) != `{"b":1,"a":2}` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeValueObjectSortsWhenDeterministic(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Ordering.Mode = jsonconfig.OrderingDeterministic
	})
	object := value.MustRecordValue(
		value.MustRecordMember("b", value.Int64Value(1)),
		value.MustRecordMember("a", value.Int64Value(2)),
	)

	data, err := c.EncodeValue(object)
	requireNoError(t, err)

	if string(data) != `{"a":2,"b":1}` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeValuePretty(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Output.Layout = jsonconfig.LayoutPretty
	})
	object := value.MustRecordValue(value.MustRecordMember("a", value.NullValue()))

	data, err := c.EncodeValue(object)
	requireNoError(t, err)

	if string(data) != "{\n  \"a\": null\n}" {
		t.Fatalf("encoded = %q", data)
	}
}

func TestEncodeValueUsesFinalNewline(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Output.FinalNewline = jsonconfig.FinalNewlineAppend
	})

	data, err := c.EncodeValue(value.BoolValue(true))
	requireNoError(t, err)

	if string(data) != "true\n" {
		t.Fatalf("encoded = %q; want final newline", data)
	}
}

func TestEncodeValueEnforcesMaxOutputBytes(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Limits.MaxOutputBytes = 4
	})

	_, err := c.EncodeValue(value.StringValue("abcd"))

	requireErrorIs(t, err, codec.ErrEncodeFailed)
}

func TestEncodeValueEnforcesMaxDepth(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Limits.MaxDepth = 1
	})

	_, err := c.EncodeValue(value.MustListValue(value.NullValue()))

	requireErrorIs(t, err, codec.ErrDepthExceeded)
}

func TestEncodeValueRejectsFloatWhenConfigured(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Numbers.FloatFormat = jsonconfig.FloatFormatReject
	})

	_, err := c.EncodeValue(mustFloatValue(t, 1.25))

	requireErrorIs(t, err, ErrUnsupportedValue)
}

func TestEncodeValueNormalizesNegativeZeroBeforeEncoding(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Numbers.NegativeZero = jsonconfig.NegativeZeroReject
	})

	data, err := c.EncodeValue(mustFloatValue(t, math.Copysign(0, -1)))
	requireNoError(t, err)

	if string(data) != `0` {
		t.Fatalf("encoded = %s; want 0", data)
	}
}

func TestEncodeValueRejectsInvalidZeroValue(t *testing.T) {
	_, err := newTestCodec(t).EncodeValue(value.Value{})

	requireErrorIs(t, err, ErrUnsupportedValue)
	requireErrorIs(t, err, codec.ErrUnsupportedFeature)
}

func TestEncodeValueRejectsBytes(t *testing.T) {
	_, err := newTestCodec(t).EncodeValue(value.BytesValue([]byte("abc")))

	requireErrorIs(t, err, ErrUnsupportedValue)
}

func TestEncodeValueRejectsTimestamp(t *testing.T) {
	_, err := newTestCodec(t).EncodeValue(value.TimestampValue(time.Now()))

	requireErrorIs(t, err, ErrUnsupportedValue)
}

func TestEncodeValueRejectsDate(t *testing.T) {
	_, err := newTestCodec(t).EncodeValue(mustDateValue(t))

	requireErrorIs(t, err, ErrUnsupportedValue)
}

func TestEncodeValueRejectsTimeOfDay(t *testing.T) {
	_, err := newTestCodec(t).EncodeValue(mustTimeOfDayValue(t))

	requireErrorIs(t, err, ErrUnsupportedValue)
}

func TestEncodeValueRejectsDuration(t *testing.T) {
	_, err := newTestCodec(t).EncodeValue(value.DurationValue(time.Second))

	requireErrorIs(t, err, ErrUnsupportedValue)
}

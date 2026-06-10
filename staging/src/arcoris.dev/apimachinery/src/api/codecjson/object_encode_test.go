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

func TestEncodeObjectEnvelope(t *testing.T) {
	desired := value.MustRecordValue(value.MustRecordMember("replicas", value.Int64Value(3)))
	observed := value.MustRecordValue(value.MustRecordMember("ready", value.BoolValue(true)))
	obj := codec.Object{
		TypeMeta:   testTypeMeta(t),
		ObjectMeta: testObjectMeta(),
		Desired:    desired,
		Observed:   &observed,
	}

	data, err := newTestCodec(t).EncodeObject(obj)
	requireNoError(t, err)

	want := `{"apiVersion":"control.arcoris.dev/v1","kind":"Worker","metadata":{"name":"example","namespace":"default"},"desired":{"replicas":3},"observed":{"ready":true}}`
	if string(data) != want {
		t.Fatalf("encoded = %s; want %s", data, want)
	}
}

func TestEncodeObjectOmitsObservedWhenAbsent(t *testing.T) {
	obj := codec.Object{Desired: value.NullValue()}

	data, err := newTestCodec(t).EncodeObject(obj)
	requireNoError(t, err)

	if string(data) != `{"desired":null}` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeObjectOmitsZeroMetadata(t *testing.T) {
	obj := codec.Object{
		TypeMeta: testTypeMeta(t),
		Desired:  value.NullValue(),
	}

	data, err := newTestCodec(t).EncodeObject(obj)
	requireNoError(t, err)

	if string(data) != `{"apiVersion":"control.arcoris.dev/v1","kind":"Worker","desired":null}` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeObjectIncludesDesired(t *testing.T) {
	obj := codec.Object{Desired: value.BoolValue(false)}

	data, err := newTestCodec(t).EncodeObject(obj)
	requireNoError(t, err)

	if string(data) != `{"desired":false}` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeObjectDeterministicDesired(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Ordering.Mode = jsonconfig.OrderingDeterministic
	})
	obj := codec.Object{Desired: value.MustRecordValue(
		value.MustRecordMember("b", value.Int64Value(1)),
		value.MustRecordMember("a", value.Int64Value(2)),
	)}

	data, err := c.EncodeObject(obj)
	requireNoError(t, err)

	if string(data) != `{"desired":{"a":2,"b":1}}` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeObjectPretty(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Output.Layout = jsonconfig.LayoutPretty
	})
	obj := codec.Object{Desired: value.NullValue()}

	data, err := c.EncodeObject(obj)
	requireNoError(t, err)

	if string(data) != "{\n  \"desired\": null\n}" {
		t.Fatalf("encoded = %q", data)
	}
}

func TestEncodeObjectRequiresTypeMetaWhenConfigured(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Object.TypeMeta = jsonconfig.TypeMetaRequire
	})

	_, err := c.EncodeObject(codec.Object{Desired: value.NullValue()})

	requireErrorIs(t, err, ErrInvalidEnvelope)
}

func TestEncodeObjectEmitsEmptyMetadataWhenConfigured(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Object.Metadata = jsonconfig.MetadataEmitEmpty
	})

	data, err := c.EncodeObject(codec.Object{Desired: value.BoolValue(false)})
	requireNoError(t, err)

	if string(data) != `{"metadata":{},"desired":false}` {
		t.Fatalf("encoded = %s", data)
	}
}

func TestEncodeObjectEmitsNullObservedWhenConfigured(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Encode.Object.Observed = jsonconfig.ObservedEmitNullWhenAbsent
	})

	data, err := c.EncodeObject(codec.Object{Desired: value.NullValue()})
	requireNoError(t, err)

	if string(data) != `{"desired":null,"observed":null}` {
		t.Fatalf("encoded = %s", data)
	}
}

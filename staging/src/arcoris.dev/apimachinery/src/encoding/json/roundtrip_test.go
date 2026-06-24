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

package jsoncodec

import (
	"bytes"
	"reflect"
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/value"
)

func TestValueByteAndStreamRoundTrip(t *testing.T) {
	c := newTestCodec(t)
	original := value.MustRecordValue(
		value.MustRecordMember("image", value.StringValue("api:v1")),
		value.MustRecordMember("ports", value.MustListValue(value.Int64Value(80), value.Int64Value(443))),
	)

	bytesEncoded, err := c.EncodeValue(original)
	requireNoError(t, err)
	bytesDecoded, err := c.DecodeValue(bytesEncoded)
	requireNoError(t, err)

	var stream bytes.Buffer
	requireNoError(t, c.EncodeValueTo(&stream, original))
	streamDecoded, err := c.DecodeValueFrom(bytes.NewReader(stream.Bytes()))
	requireNoError(t, err)

	if !reflect.DeepEqual(bytesDecoded, original) {
		t.Fatalf("byte roundtrip = %#v; want %#v", bytesDecoded, original)
	}
	if !reflect.DeepEqual(streamDecoded, original) {
		t.Fatalf("stream roundtrip = %#v; want %#v", streamDecoded, original)
	}
	if stream.String() != string(bytesEncoded) {
		t.Fatalf("stream bytes = %s; byte bytes = %s", stream.String(), bytesEncoded)
	}
}

func TestObjectByteAndStreamRoundTrip(t *testing.T) {
	c := newTestCodec(t)
	observed := value.MustRecordValue(value.MustRecordMember("ready", value.BoolValue(true)))
	objectMeta := testObjectMeta()
	objectMeta.Labels = labels.Set{"app": "worker"}
	objectMeta.Annotations = annotations.Set{"scheduler.arcoris.dev/mode": "balanced"}

	original := codec.Object{
		TypeMeta:   testTypeMeta(t),
		ObjectMeta: objectMeta,
		Desired:    value.MustRecordValue(value.MustRecordMember("image", value.StringValue("api:v1"))),
		Observed:   &observed,
	}

	bytesEncoded, err := c.EncodeObject(original)
	requireNoError(t, err)
	bytesDecoded, err := c.DecodeObject(bytesEncoded)
	requireNoError(t, err)

	var stream bytes.Buffer
	requireNoError(t, c.EncodeObjectTo(&stream, original))
	streamDecoded, err := c.DecodeObjectFrom(bytes.NewReader(stream.Bytes()))
	requireNoError(t, err)

	if !reflect.DeepEqual(bytesDecoded, original) {
		t.Fatalf("byte roundtrip = %#v; want %#v", bytesDecoded, original)
	}
	if !reflect.DeepEqual(streamDecoded, original) {
		t.Fatalf("stream roundtrip = %#v; want %#v", streamDecoded, original)
	}
	if stream.String() != string(bytesEncoded) {
		t.Fatalf("stream bytes = %s; byte bytes = %s", stream.String(), bytesEncoded)
	}
}

func TestObjectOwnershipByteAndStreamRoundTrip(t *testing.T) {
	c := newTestCodec(t)
	original := ownershipState(
		ownershipSurface(ownershipEntry("user-cli", "$.image")),
		ownershipSurface(ownershipEntry("controller", "$.ready")),
		objectownership.NewMetadataState(
			ownershipSurface(ownershipEntry("labeler", `$["app"]`)),
			ownershipSurface(ownershipEntry("annotator", `$["scheduler.arcoris.dev/mode"]`)),
		),
	)

	bytesEncoded, err := c.EncodeObjectOwnership(original)
	requireNoError(t, err)
	bytesDecoded, err := c.DecodeObjectOwnership(bytesEncoded)
	requireNoError(t, err)

	var stream bytes.Buffer
	requireNoError(t, c.EncodeObjectOwnershipTo(&stream, original))
	streamDecoded, err := c.DecodeObjectOwnershipFrom(bytes.NewReader(stream.Bytes()))
	requireNoError(t, err)

	if !reflect.DeepEqual(bytesDecoded, original) {
		t.Fatalf("byte roundtrip = %#v; want %#v", bytesDecoded, original)
	}
	if !reflect.DeepEqual(streamDecoded, original) {
		t.Fatalf("stream roundtrip = %#v; want %#v", streamDecoded, original)
	}
	if stream.String() != string(bytesEncoded) {
		t.Fatalf("stream bytes = %s; byte bytes = %s", stream.String(), bytesEncoded)
	}
}

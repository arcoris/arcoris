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

	"arcoris.dev/apimachinery/api/codecjson/jsonconfig"
	"arcoris.dev/apimachinery/api/value"
)

func TestDecodeObjectEnvelope(t *testing.T) {
	data := []byte(`{"apiVersion":"control.arcoris.dev/v1","kind":"Worker","metadata":{"name":"example","namespace":"default"},"desired":{"replicas":3},"observed":{"ready":true}}`)

	got, err := newTestCodec(t).DecodeObject(data)
	requireNoError(t, err)

	if got.TypeMeta.String() != "control.arcoris.dev/v1#Worker" {
		t.Fatalf("type meta = %q", got.TypeMeta.String())
	}
	if got.ObjectMeta.Name.String() != "example" || got.ObjectMeta.Namespace.String() != "default" {
		t.Fatalf("object meta = %#v", got.ObjectMeta)
	}
	requireObjectMemberNames(t, got.Desired, "replicas")
	if got.Observed == nil {
		t.Fatalf("observed is nil")
	}
	requireObjectMemberNames(t, *got.Observed, "ready")
}

func TestDecodeObjectWithoutObserved(t *testing.T) {
	got, err := newTestCodec(t).DecodeObject([]byte(`{"desired":{}}`))
	requireNoError(t, err)

	if got.Observed != nil {
		t.Fatalf("observed = %#v; want nil", got.Observed)
	}
}

func TestDecodeObjectWithObserved(t *testing.T) {
	got, err := newTestCodec(t).DecodeObject([]byte(`{"desired":{},"observed":{"ready":true}}`))
	requireNoError(t, err)

	if got.Observed == nil {
		t.Fatalf("observed is nil")
	}
}

func TestDecodeObjectWithoutMetadata(t *testing.T) {
	got, err := newTestCodec(t).DecodeObject([]byte(`{"apiVersion":"v1","kind":"Worker","desired":{}}`))
	requireNoError(t, err)

	if !got.ObjectMeta.IsZero() {
		t.Fatalf("metadata = %#v; want zero", got.ObjectMeta)
	}
}

func TestDecodeObjectRejectsRootNonObject(t *testing.T) {
	_, err := newTestCodec(t).DecodeObject([]byte(`[]`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$", ErrorReasonInvalidEnvelope)
}

func TestDecodeObjectRejectsMissingDesired(t *testing.T) {
	_, err := newTestCodec(t).DecodeObject([]byte(`{"metadata":{"name":"example"}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$.desired", ErrorReasonInvalidEnvelope)
}

func TestDecodeObjectRejectsUnknownEnvelopeField(t *testing.T) {
	_, err := newTestCodec(t).DecodeObject([]byte(`{"desired":{},"extra":true}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$.extra", ErrorReasonInvalidEnvelope)
}

func TestDecodeObjectIgnoresUnknownEnvelopeFieldsWhenConfigured(t *testing.T) {
	c := newTestCodecWith(t, func(config *jsonconfig.Config) {
		config.Decode.Objects.UnknownEnvelopeFields = jsonconfig.UnknownFieldIgnore
	})

	obj, err := c.DecodeObject([]byte(`{"desired":true,"extra":1}`))
	requireNoError(t, err)

	if obj.Desired.Kind() != value.KindBool {
		t.Fatalf("desired kind = %s; want bool", obj.Desired.Kind())
	}
}

func TestDecodeObjectRejectsDuplicateEnvelopeField(t *testing.T) {
	_, err := newTestCodec(t).DecodeObject([]byte(`{"desired":{},"desired":null}`))

	requireErrorIs(t, err, ErrDuplicateKey)
}

func TestDecodeObjectRejectsInvalidAPIVersionKindTypes(t *testing.T) {
	testCases := []string{
		`{"apiVersion":1,"kind":"Worker","desired":{}}`,
		`{"apiVersion":"v1","kind":1,"desired":{}}`,
		`{"apiVersion":"not valid","kind":"Worker","desired":{}}`,
		`{"apiVersion":"v1","kind":"worker","desired":{}}`,
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			_, err := newTestCodec(t).DecodeObject([]byte(input))
			requireErrorIs(t, err, ErrInvalidEnvelope)
		})
	}
}

func TestDecodeObjectRejectsMetadataNonObject(t *testing.T) {
	_, err := newTestCodec(t).DecodeObject([]byte(`{"metadata":[],"desired":{}}`))

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$.metadata", ErrorReasonInvalidEnvelope)
}

func TestDecodeObjectPreservesDesiredNull(t *testing.T) {
	got, err := newTestCodec(t).DecodeObject([]byte(`{"desired":null}`))
	requireNoError(t, err)

	requireKind(t, got.Desired, value.KindNull)
}

func TestDecodeObjectPreservesObservedNullWhenPresent(t *testing.T) {
	got, err := newTestCodec(t).DecodeObject([]byte(`{"desired":{},"observed":null}`))
	requireNoError(t, err)

	if got.Observed == nil {
		t.Fatalf("observed is nil")
	}
	requireKind(t, *got.Observed, value.KindNull)
}

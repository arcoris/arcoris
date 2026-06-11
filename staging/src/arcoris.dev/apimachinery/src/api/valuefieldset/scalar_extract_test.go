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

package valuefieldset

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestExtractOwnershipFieldsScalarIncludesCurrentPath(t *testing.T) {
	path := rootField("spec", "replicas")

	got, err := ExtractOwnershipFieldsAt(
		path,
		value.Int64Value(3),
		types.Int64().Descriptor(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractOwnershipFieldsScalarMismatchReturnsError(t *testing.T) {
	path := rootField("spec", "replicas")

	_, err := ExtractOwnershipFieldsAt(
		path,
		value.StringValue("three"),
		types.Int64().Descriptor(),
		Options{},
	)

	requireErrorIs(t, err, ErrKindMismatch)
	requireErrorReason(t, err, ErrorReasonKindMismatch)
	requireErrorPath(t, err, "$.spec.replicas")
}

func TestScalarDescriptorKind(t *testing.T) {
	tests := []struct {
		name string
		code types.DescriptorKind
		want value.Kind
		ok   bool
	}{
		{name: "bool", code: types.DescriptorBool, want: value.KindBool, ok: true},
		{name: "string", code: types.DescriptorString, want: value.KindString, ok: true},
		{name: "bytes", code: types.DescriptorBytes, want: value.KindBytes, ok: true},
		{name: "signed integer", code: types.DescriptorInt64, want: value.KindInteger, ok: true},
		{name: "unsigned integer", code: types.DescriptorUint64, want: value.KindInteger, ok: true},
		{name: "float", code: types.DescriptorFloat64, want: value.KindFloat, ok: true},
		{name: "decimal", code: types.DescriptorDecimal, want: value.KindDecimal, ok: true},
		{name: "timestamp", code: types.DescriptorTimestamp, want: value.KindTimestamp, ok: true},
		{name: "date", code: types.DescriptorDate, want: value.KindDate, ok: true},
		{name: "time", code: types.DescriptorTime, want: value.KindTimeOfDay, ok: true},
		{name: "duration", code: types.DescriptorDuration, want: value.KindDuration, ok: true},
		{name: "object", code: types.DescriptorObject, want: value.KindInvalid, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := scalarKind(tt.code)

			if got != tt.want || ok != tt.ok {
				t.Fatalf("scalarKind(%s) = %s, %v; want %s, %v", tt.code, got, ok, tt.want, tt.ok)
			}
		})
	}
}

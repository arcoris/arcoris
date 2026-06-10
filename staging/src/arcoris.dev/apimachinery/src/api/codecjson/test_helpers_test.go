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
	"reflect"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/codecjson/jsonconfig"
	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/value"
)

// requireNoError fails the current test when err is non-nil.
func requireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// newTestCodec constructs a codec with safe default configuration.
func newTestCodec(t *testing.T) Codec {
	t.Helper()

	codec, err := New(jsonconfig.Default())
	requireNoError(t, err)

	return codec
}

// newTestCodecWith constructs a codec from a targeted config override.
func newTestCodecWith(t *testing.T, mutate func(*jsonconfig.Config)) Codec {
	t.Helper()

	config := jsonconfig.Default()
	mutate(&config)
	codec, err := New(config)
	requireNoError(t, err)

	return codec
}

// requireErrorIs asserts errors.Is compatibility for sentinel contracts.
func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want errors.Is(..., %v)", err, target)
	}
}

// requireCodecJSONError asserts the structured codecjson diagnostic surface.
func requireCodecJSONError(t *testing.T, err error, path string, reason ErrorReason) {
	t.Helper()

	var codecErr *Error
	if !errors.As(err, &codecErr) {
		t.Fatalf("error = %T; want *Error", err)
	}
	if codecErr.Path != path {
		t.Fatalf("path = %q; want %q", codecErr.Path, path)
	}
	if codecErr.Reason != reason {
		t.Fatalf("reason = %q; want %q", codecErr.Reason, reason)
	}
	if codecErr.Detail == "" {
		t.Fatalf("detail is empty")
	}
}

// requireDetailContains asserts useful human detail without coupling tests to
// an entire diagnostic sentence.
func requireDetailContains(t *testing.T, err error, want string) {
	t.Helper()

	var codecErr *Error
	if !errors.As(err, &codecErr) {
		t.Fatalf("error = %T; want *Error", err)
	}
	if !strings.Contains(codecErr.Detail, want) {
		t.Fatalf("detail = %q; want substring %q", codecErr.Detail, want)
	}
}

// requireKind asserts the decoded or constructed value kind.
func requireKind(t *testing.T, got value.Value, want value.Kind) {
	t.Helper()

	if got.Kind() != want {
		t.Fatalf("kind = %s; want %s", got.Kind(), want)
	}
}

// requireStringValue asserts a string value payload.
func requireStringValue(t *testing.T, got value.Value, want string) {
	t.Helper()

	requireKind(t, got, value.KindString)
	text, ok := got.AsString()
	if !ok || text != want {
		t.Fatalf("string = %q ok=%v; want %q", text, ok, want)
	}
}

// requireIntegerText asserts integer value text instead of machine integer
// width, which keeps signed and unsigned boundary tests precise.
func requireIntegerText(t *testing.T, got value.Value, want string) {
	t.Helper()

	requireKind(t, got, value.KindInteger)
	integer, ok := got.AsInteger()
	if !ok || integer.String() != want {
		t.Fatalf("integer = %q ok=%v; want %q", integer.String(), ok, want)
	}
}

// requireDecimalText asserts decimal value text after exact JSON expansion.
func requireDecimalText(t *testing.T, got value.Value, want string) {
	t.Helper()

	requireKind(t, got, value.KindDecimal)
	decimal, ok := got.AsDecimal()
	if !ok || decimal.String() != want {
		t.Fatalf("decimal = %q ok=%v; want %q", decimal.String(), ok, want)
	}
}

// requireObjectMemberNames asserts object member order.
//
// Order is part of codecjson behavior because value decoding preserves JSON
// object member order and deterministic encoding opts in to sorting explicitly.
func requireObjectMemberNames(t *testing.T, got value.Value, want ...string) {
	t.Helper()

	requireKind(t, got, value.KindRecord)
	object, ok := got.AsRecord()
	if !ok {
		t.Fatalf("object view unavailable")
	}
	wantNames := make([]value.MemberName, 0, len(want))
	for _, name := range want {
		wantNames = append(wantNames, value.MemberName(name))
	}
	if names := object.Names(); !reflect.DeepEqual(names, wantNames) {
		t.Fatalf("object names = %#v; want %#v", names, want)
	}
}

// mustDecimalValue constructs a decimal value for encode tests.
func mustDecimalValue(t *testing.T, text string) value.Value {
	t.Helper()

	decimal, err := value.ParseDecimal(text)
	requireNoError(t, err)

	return value.DecimalValue(decimal)
}

// mustFloatValue constructs a finite float value for encode tests.
func mustFloatValue(t *testing.T, f float64) value.Value {
	t.Helper()

	v, err := value.FloatValue(f)
	requireNoError(t, err)

	return v
}

// mustDateValue constructs a descriptor-dependent value kind that generic JSON
// must reject during encode.
func mustDateValue(t *testing.T) value.Value {
	t.Helper()

	date, err := value.NewDate(2026, 6, 4)
	requireNoError(t, err)
	v, err := value.DateValue(date)
	requireNoError(t, err)

	return v
}

// mustTimeOfDayValue constructs a descriptor-dependent value kind that generic
// JSON must reject during encode.
func mustTimeOfDayValue(t *testing.T) value.Value {
	t.Helper()

	timeOfDay, err := value.NewTimeOfDay(12, 30, 0, 0)
	requireNoError(t, err)
	v, err := value.TimeOfDayValue(timeOfDay)
	requireNoError(t, err)

	return v
}

// mustGroupVersion parses a group/version literal used by object envelope tests.
func mustGroupVersion(t *testing.T, text string) apiidentity.GroupVersion {
	t.Helper()

	groupVersion, err := apiidentity.ParseGroupVersion(text)
	requireNoError(t, err)

	return groupVersion
}

// testTypeMeta returns stable type metadata for object envelope fixtures.
func testTypeMeta(t *testing.T) meta.TypeMeta {
	t.Helper()

	return meta.TypeMeta{
		APIVersion: mustGroupVersion(t, "control.arcoris.dev/v1"),
		Kind:       apiidentity.Kind("Worker"),
	}
}

// testObjectMeta returns stable object metadata for object envelope fixtures.
func testObjectMeta() meta.ObjectMeta {
	return meta.ObjectMeta{
		Name:      metaidentity.Name("example"),
		Namespace: metaidentity.Namespace("default"),
	}
}

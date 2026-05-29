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

package identity

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// identityMarshaler is the shared scalar-encoding surface implemented by every
// identity value. Tests use it to verify the package-wide scalar contract
// without repeating the same marshaling checks in every domain file.
type identityMarshaler interface {
	String() string
	Validate() error
	MarshalText() ([]byte, error)
	MarshalJSON() ([]byte, error)
}

// identityUnmarshaler is the pointer-side scalar-decoding surface implemented
// by every identity value. It keeps encoding tests generic while still forcing
// each concrete type to exercise its own UnmarshalText and UnmarshalJSON
// methods.
type identityUnmarshaler interface {
	String() string
	Validate() error
	UnmarshalText([]byte) error
	UnmarshalJSON([]byte) error
}

// comparableIdentity is the value-level contract used by parser roundtrip and
// fuzz tests. Identity structs are intentionally small comparable values, so a
// successful parse can be compared directly with its re-parsed canonical form.
type comparableIdentity interface {
	identityMarshaler
	comparable
}

// identityIdentifier is implemented by composite identities that expose a
// stable diagnostic/map-key spelling in addition to String.
type identityIdentifier interface {
	String() string
	Identifier() string
}

// identityPointer is the pointer-side companion for a concrete identity value.
//
// Each public identity type has value-receiver marshal methods and
// pointer-receiver unmarshal methods. The helper keeps that shape explicit in
// tests without hiding the fact that every concrete type still has its own
// encoding file.
type identityPointer[T any] interface {
	*T
	identityUnmarshaler
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func requireIdentityError(t *testing.T, err error, name string, reason ErrorReason) *Error {
	t.Helper()
	requireErrorIs(t, err, ErrInvalidIdentifier)
	var identityErr *Error
	if !errors.As(err, &identityErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if identityErr.Name != name {
		t.Fatalf("Error.Name = %q, want %q", identityErr.Name, name)
	}
	if identityErr.Reason != reason {
		t.Fatalf("Error.Reason = %q, want %q", identityErr.Reason, reason)
	}
	if identityErr.Detail == "" {
		t.Fatalf("Error.Detail is empty")
	}
	return identityErr
}

func requireDetailContains(t *testing.T, err error, text string) {
	t.Helper()
	var identityErr *Error
	if !errors.As(err, &identityErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if !strings.Contains(identityErr.Detail, text) {
		t.Fatalf("Error.Detail = %q, want to contain %q", identityErr.Detail, text)
	}
}

func requireString(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func requireEqual[T comparable](t *testing.T, label string, got T, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("%s = %#v, want %#v", label, got, want)
	}
}

func requireIdentifier(t *testing.T, value identityIdentifier, want string) {
	t.Helper()
	if got := value.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
	if got := value.Identifier(); got != want {
		t.Fatalf("Identifier() = %q, want %q", got, want)
	}
}

func requireValidIdentity(t *testing.T, value identityMarshaler) {
	t.Helper()
	requireNoError(t, value.Validate())
}

func requireInvalidIdentity(
	t *testing.T,
	value identityMarshaler,
	name string,
	reason ErrorReason,
	detail string,
) {
	t.Helper()
	err := value.Validate()
	requireIdentityError(t, err, name, reason)
	requireDetailContains(t, err, detail)
}

func requireIdentityEncoding[T identityMarshaler, PT identityPointer[T]](
	t *testing.T,
	input string,
	value T,
	invalid T,
) {
	t.Helper()

	text, err := value.MarshalText()
	requireNoError(t, err)
	if string(text) != input {
		t.Fatalf("MarshalText() = %q, want %q", text, input)
	}

	parsedText := PT(new(T))
	requireNoError(t, parsedText.UnmarshalText(text))
	requireNoError(t, parsedText.Validate())
	if parsedText.String() != input {
		t.Fatalf("UnmarshalText String() = %q, want %q", parsedText.String(), input)
	}

	jsonData, err := value.MarshalJSON()
	requireNoError(t, err)
	var jsonString string
	requireNoError(t, json.Unmarshal(jsonData, &jsonString))
	if jsonString != input {
		t.Fatalf("MarshalJSON() scalar = %q, want %q", jsonString, input)
	}

	parsedJSON := PT(new(T))
	requireNoError(t, parsedJSON.UnmarshalJSON(jsonData))
	requireNoError(t, parsedJSON.Validate())
	if parsedJSON.String() != input {
		t.Fatalf("UnmarshalJSON String() = %q, want %q", parsedJSON.String(), input)
	}

	_, err = invalid.MarshalText()
	requireErrorIs(t, err, ErrInvalidIdentifier)

	_, err = invalid.MarshalJSON()
	requireErrorIs(t, err, ErrInvalidIdentifier)

	requireRejectsNonStringJSON[T, PT](t)
	requireRejectsNilReceiver[T, PT](t, input)
}

func requireRejectsNonStringJSON[T identityMarshaler, PT identityPointer[T]](t *testing.T) {
	t.Helper()

	for _, payload := range []string{"null", `{}`, `[]`, `1`, `true`} {
		target := PT(new(T))
		err := target.UnmarshalJSON([]byte(payload))
		requireErrorIs(t, err, ErrInvalidJSON)

		var identityErr *Error
		if !errors.As(err, &identityErr) {
			t.Fatalf("expected *Error, got %T", err)
		}
		if identityErr.Reason != ErrorReasonInvalidJSON {
			t.Fatalf("Error.Reason = %q, want %q", identityErr.Reason, ErrorReasonInvalidJSON)
		}
	}
}

func requireRejectsNilReceiver[T identityMarshaler, PT identityPointer[T]](
	t *testing.T,
	input string,
) {
	t.Helper()

	var nilTarget PT
	requireErrorIs(t, nilTarget.UnmarshalText([]byte(input)), ErrNilReceiver)

	err := nilTarget.UnmarshalJSON([]byte(`"` + input + `"`))
	requireErrorIs(t, err, ErrNilReceiver)

	var identityErr *Error
	if !errors.As(err, &identityErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if identityErr.Reason != ErrorReasonNilReceiver {
		t.Fatalf("Error.Reason = %q, want %q", identityErr.Reason, ErrorReasonNilReceiver)
	}
}

func requireParseOK[T comparableIdentity](t *testing.T, input string, parse func(string) (T, error)) T {
	t.Helper()
	value, err := parse(input)
	requireNoError(t, err)
	requireNoError(t, value.Validate())
	if value.String() != input {
		t.Fatalf("String() = %q, want %q", value.String(), input)
	}
	return value
}

func requireParseReject[T comparableIdentity](t *testing.T, input string, parse func(string) (T, error)) {
	t.Helper()
	_, err := parse(input)
	requireErrorIs(t, err, ErrInvalidIdentifier)
}

func requireParseRoundTrip[T comparableIdentity](t *testing.T, input string, parse func(string) (T, error)) {
	t.Helper()
	value := requireParseOK(t, input, parse)
	again, err := parse(value.String())
	requireNoError(t, err)
	if again != value {
		t.Fatalf("parse(String()) = %q, want %q", again.String(), value.String())
	}

	text, err := value.MarshalText()
	requireNoError(t, err)
	if string(text) != input {
		t.Fatalf("MarshalText() = %q, want %q", text, input)
	}
	fromText, err := parse(string(text))
	requireNoError(t, err)
	if fromText != value {
		t.Fatalf("text roundtrip = %q, want %q", fromText.String(), value.String())
	}

	jsonData, err := value.MarshalJSON()
	requireNoError(t, err)
	var jsonText string
	requireNoError(t, json.Unmarshal(jsonData, &jsonText))
	if jsonText != input {
		t.Fatalf("MarshalJSON() scalar = %q, want %q", jsonText, input)
	}
	fromJSON, err := parse(jsonText)
	requireNoError(t, err)
	if fromJSON != value {
		t.Fatalf("JSON roundtrip = %q, want %q", fromJSON.String(), value.String())
	}
}

func requireCanonicalSet[T comparableIdentity](t *testing.T, inputs []string, parse func(string) (T, error)) {
	t.Helper()
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			requireParseRoundTrip(t, input, parse)
		})
	}
}

func requireRejectedSet[T comparableIdentity](t *testing.T, inputs []string, parse func(string) (T, error)) {
	t.Helper()
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			requireParseReject(t, input, parse)
		})
	}
}

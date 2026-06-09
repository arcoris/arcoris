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

package resource

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/types"
)

type fakeResolver map[types.TypeName]types.Definition

func (f fakeResolver) Resolve(name types.TypeName) (types.Definition, bool) {
	def, ok := f[name]
	return def, ok
}

func objectType() types.Descriptor { return types.Object().Descriptor() }

func stringType() types.Descriptor { return types.String().Descriptor() }

func refType(name string) types.Descriptor { return types.Ref(name).Descriptor() }

func validVersion(options ...VersionOption) VersionDefinition {
	opts := append([]VersionOption{Observed(objectType()), Exposed(), Canonical()}, options...)
	return NewVersion(identity.Version("v1"), objectType(), opts...)
}

func validDefinition() Definition {
	return NewDefinition(
		identity.Group("control.arcoris.dev"),
		identity.Kind("Worker"),
		identity.Resource("workers"),
		ScopeNamespaced,
		NewVersion(
			identity.Version("v1"),
			objectType(),
			Observed(objectType()),
			Exposed(),
			Canonical(),
		),
	)
}

// requireEqual compares small descriptor values in tests without repeating
// boilerplate failure messages.
func requireEqual[T comparable](t *testing.T, got T, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
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

func requireResourceError(t *testing.T, err error, target error, path string, reason ErrorReason) *Error {
	t.Helper()
	requireErrorIs(t, err, target)
	var resourceErr *Error
	if !errors.As(err, &resourceErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if resourceErr.Path != path {
		t.Fatalf("Error.Path = %q, want %q", resourceErr.Path, path)
	}
	if resourceErr.Reason != reason {
		t.Fatalf("Error.Reason = %q, want %q", resourceErr.Reason, reason)
	}
	if resourceErr.Detail == "" {
		t.Fatalf("Error.Detail is empty")
	}
	return resourceErr
}

func requireDetailContains(t *testing.T, err error, text string) {
	t.Helper()
	var resourceErr *Error
	if !errors.As(err, &resourceErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if !strings.Contains(resourceErr.Detail, text) {
		t.Fatalf("Error.Detail = %q, want to contain %q", resourceErr.Detail, text)
	}
}

func requireScopeJSONRejects(t *testing.T, payload string) {
	t.Helper()
	var scope Scope
	err := scope.UnmarshalJSON([]byte(payload))
	requireResourceError(t, err, ErrInvalidJSON, pathScope, ErrorReasonInvalidJSON)
}

func requireScopeJSONRoundTrip(t *testing.T, scope Scope, text string) {
	t.Helper()
	data, err := scope.MarshalJSON()
	requireNoError(t, err)
	var scalar string
	requireNoError(t, json.Unmarshal(data, &scalar))
	if scalar != text {
		t.Fatalf("MarshalJSON scalar = %q, want %q", scalar, text)
	}
	var parsed Scope
	requireNoError(t, parsed.UnmarshalJSON(data))
	if parsed != scope {
		t.Fatalf("UnmarshalJSON = %v, want %v", parsed, scope)
	}
}

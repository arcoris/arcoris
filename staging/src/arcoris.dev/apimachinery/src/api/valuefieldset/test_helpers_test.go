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
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
)

type testResolver map[types.TypeName]types.TypeDefinition

func (r testResolver) ResolveType(name types.TypeName) (types.TypeDefinition, bool) {
	definition, ok := r[name]
	return definition, ok
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func requireFieldSet(t *testing.T, got fieldpath.Set, want ...fieldpath.Path) {
	t.Helper()

	expected := fieldpath.MustSet(want...)
	if !got.Equal(expected) {
		t.Fatalf("set = %s, want %s", got, expected)
	}
}

func requireErrorPath(t *testing.T, err error, want string) {
	t.Helper()

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(*Error) = false for %v", err)
	}
	if got.Path != want {
		t.Fatalf("error path = %q, want %q", got.Path, want)
	}
}

func requireErrorReason(t *testing.T, err error, want ErrorReason) {
	t.Helper()

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(*Error) = false for %v", err)
	}
	if got.Reason != want {
		t.Fatalf("error reason = %q, want %q", got.Reason, want)
	}
}

func requireErrorDetailContains(t *testing.T, err error, want string) {
	t.Helper()

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(*Error) = false for %v", err)
	}
	if !strings.Contains(got.Detail, want) {
		t.Fatalf("error detail = %q, want substring %q", got.Detail, want)
	}
}

func rootField(names ...string) fieldpath.Path {
	path := fieldpath.RootPath()
	for _, name := range names {
		path = path.Field(name)
	}

	return path
}

func readySelector() fieldpath.Selector {
	return fieldpath.MustSelector(
		fieldpath.NewSelectorEntry("type", fieldpath.StringLiteral("Ready")),
	)
}

func routeSelector() fieldpath.Selector {
	return fieldpath.MustSelector(
		fieldpath.NewSelectorEntry("port", fieldpath.Uint64Literal(443)),
		fieldpath.NewSelectorEntry("host", fieldpath.StringLiteral("api.example.com")),
	)
}

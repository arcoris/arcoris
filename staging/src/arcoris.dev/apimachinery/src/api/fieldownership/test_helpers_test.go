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

package fieldownership

import (
	"errors"
	"reflect"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func path(text string) fieldpath.Path {
	p, err := fieldpath.ParseCanonical(text)
	if err != nil {
		panic(err)
	}

	return p
}

func set(paths ...fieldpath.Path) fieldpath.Set {
	fields, err := fieldpath.NewSet(paths...)
	if err != nil {
		panic(err)
	}

	return fields
}

func owner(text string) Owner {
	return MustOwner(text)
}

func invalidOwner(text string) Owner {
	return Owner{text: text}
}

func entry(ownerText string, paths ...fieldpath.Path) Entry {
	return MustEntry(owner(ownerText), set(paths...))
}

func emptyEntry(ownerText string) Entry {
	return MustEntry(owner(ownerText), fieldpath.EmptySet())
}

func owners(values ...string) []Owner {
	out := make([]Owner, 0, len(values))
	for _, value := range values {
		out = append(out, owner(value))
	}

	return out
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

func requireFieldOwnershipError(t *testing.T, err error, path string, reason ErrorReason) {
	t.Helper()

	var fieldError *Error
	if !errors.As(err, &fieldError) {
		t.Fatalf("errors.As(%v, *fieldownership.Error) = false", err)
	}
	if fieldError.Path != path {
		t.Fatalf("error path = %q; want %q", fieldError.Path, path)
	}
	if fieldError.Reason != reason {
		t.Fatalf("error reason = %q; want %q", fieldError.Reason, reason)
	}
}

func requireEqual[T comparable](t *testing.T, got T, want T) {
	t.Helper()

	if got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func requireOwners(t *testing.T, got []Owner, want ...string) {
	t.Helper()

	expected := owners(want...)
	if len(got) == 0 && len(expected) == 0 {
		return
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("owners = %#v; want %#v", got, expected)
	}
}

func requirePathStrings(t *testing.T, got []fieldpath.Path, want ...string) {
	t.Helper()

	strings := make([]string, 0, len(got))
	for _, p := range got {
		strings = append(strings, p.String())
	}

	if len(strings) == 0 && len(want) == 0 {
		return
	}

	if !reflect.DeepEqual(strings, want) {
		t.Fatalf("paths = %#v; want %#v", strings, want)
	}
}

func requireSet(t *testing.T, got fieldpath.Set, want ...string) {
	t.Helper()

	requirePathStrings(t, got.Paths(), want...)
}

func requirePanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic")
		}
	}()

	fn()
}

func specPath() fieldpath.Path {
	return path("$.spec")
}

func imagePath() fieldpath.Path {
	return path("$.spec.image")
}

func replicasPath() fieldpath.Path {
	return path("$.spec.replicas")
}

func metadataPath() fieldpath.Path {
	return path("$.metadata")
}

func namePath() fieldpath.Path {
	return path("$.metadata.name")
}

func labelPath() fieldpath.Path {
	return path(`$.metadata.labels["app"]`)
}

func argsPath() fieldpath.Path {
	return path("$.args")
}

func argsIndexPath() fieldpath.Path {
	return path("$.args[0]")
}

func readyPath() fieldpath.Path {
	return path(`$.conditions[{"type":"Ready"}]`)
}

func readyStatusPath() fieldpath.Path {
	return path(`$.conditions[{"type":"Ready"}].status`)
}

func baseState() State {
	return MustState(
		entry("user-cli", imagePath(), replicasPath()),
		entry("autoscaler", replicasPath()),
		entry("health-controller", readyStatusPath()),
	)
}

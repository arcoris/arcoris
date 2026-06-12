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

package objectownership

import (
	"errors"
	"reflect"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
)

// owner constructs a validated test owner and fails fast on fixture mistakes.
func owner(name string) fieldownership.Owner {
	return fieldownership.MustOwner(name)
}

// path parses a canonical field path used by test fixtures.
func path(text string) fieldpath.Path {
	p, err := fieldpath.ParseCanonical(text)
	if err != nil {
		panic(err)
	}

	return p
}

// fields builds a canonical field set from textual test paths.
func fields(paths ...string) fieldpath.Set {
	parsed := make([]fieldpath.Path, 0, len(paths))
	for _, text := range paths {
		parsed = append(parsed, path(text))
	}

	return fieldpath.MustSet(parsed...)
}

// ownershipEntry builds one normalized ownership entry for tests.
func ownershipEntry(name string, paths ...string) fieldownership.Entry {
	return fieldownership.MustEntry(owner(name), fields(paths...))
}

// ownershipState builds a normalized fieldownership.State for tests.
func ownershipState(entries ...fieldownership.Entry) fieldownership.State {
	return fieldownership.MustState(entries...)
}

// requireNoError fails the test when err is non-nil.
func requireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// requireErrorIs asserts that err unwraps to target.
func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

// requireObjectOwnershipError checks the structured objectownership diagnostic.
func requireObjectOwnershipError(t *testing.T, err error, path string, reason ErrorReason) {
	t.Helper()

	var ownershipErr *Error
	if !errors.As(err, &ownershipErr) {
		t.Fatalf("error type = %T; want *Error", err)
	}
	if ownershipErr.Path != path {
		t.Fatalf("Error.Path = %q; want %q", ownershipErr.Path, path)
	}
	if ownershipErr.Reason != reason {
		t.Fatalf("Error.Reason = %q; want %q", ownershipErr.Reason, reason)
	}
	if ownershipErr.Detail == "" {
		t.Fatalf("Error.Detail is empty")
	}
}

// requireOwners compares owners in deterministic fieldownership order.
func requireOwners(t *testing.T, got []fieldownership.Owner, want ...string) {
	t.Helper()

	gotStrings := make([]string, 0, len(got))
	for _, owner := range got {
		gotStrings = append(gotStrings, owner.String())
	}
	if !reflect.DeepEqual(gotStrings, want) {
		t.Fatalf("owners = %#v; want %#v", gotStrings, want)
	}
}

// requireOwnersOf checks the owners reported for one field path.
func requireOwnersOf(t *testing.T, state fieldownership.State, path fieldpath.Path, want ...string) {
	t.Helper()

	got, err := state.OwnersOf(path)
	requireNoError(t, err)
	requireOwners(t, got, want...)
}

// requirePaths compares a fieldpath.Set against canonical path strings.
func requirePaths(t *testing.T, got fieldpath.Set, want ...string) {
	t.Helper()

	gotStrings := make([]string, 0, len(got.Paths()))
	for _, path := range got.Paths() {
		gotStrings = append(gotStrings, path.String())
	}
	if !reflect.DeepEqual(gotStrings, want) {
		t.Fatalf("paths = %#v; want %#v", gotStrings, want)
	}
}

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

package objectquery

import (
	"errors"
	"reflect"
	"testing"
)

func TestErrorSentinelsAreDistinct(t *testing.T) {
	sentinels := []error{
		ErrInvalidQuery,
		ErrInvalidSelector,
		ErrInvalidRequirement,
		ErrInvalidOperator,
	}

	seen := map[string]struct{}{}
	for _, sentinel := range sentinels {
		if sentinel == nil {
			t.Fatal("sentinel is nil")
		}
		text := sentinel.Error()
		if _, ok := seen[text]; ok {
			t.Fatalf("duplicate sentinel text %q", text)
		}
		seen[text] = struct{}{}
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

func requireStrings(t *testing.T, got []string, want ...string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("strings = %#v; want %#v", got, want)
	}
}

func requireRequirement(t *testing.T, got metadataRequirement, key string, op Operator, values ...string) {
	t.Helper()
	if got.key != key {
		t.Fatalf("key = %q; want %q", got.key, key)
	}
	if got.op != op {
		t.Fatalf("operator = %s; want %s", got.op, op)
	}
	requireStrings(t, got.values, values...)
}

func requireQueryError(t *testing.T, err error, path string, reason ErrorReason) {
	t.Helper()
	var queryErr *Error
	if !errors.As(err, &queryErr) {
		t.Fatalf("errors.As(%v, *Error) = false", err)
	}
	if queryErr.Path != path {
		t.Fatalf("Path = %q; want %q", queryErr.Path, path)
	}
	if queryErr.Reason != reason {
		t.Fatalf("Reason = %q; want %q", queryErr.Reason, reason)
	}
}

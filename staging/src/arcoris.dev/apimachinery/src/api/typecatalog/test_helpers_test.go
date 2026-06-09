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

package typecatalog

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error matching %v, got %v", target, err)
	}
}

func requireErrorNotIs(t *testing.T, err, target error) {
	t.Helper()
	if errors.Is(err, target) {
		t.Fatalf("expected error not matching %v, got %v", target, err)
	}
}

func requireEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func requireDefinition(t *testing.T, c *Catalog, name types.TypeName) types.Definition {
	t.Helper()
	def, ok := c.Resolve(name)
	if !ok {
		t.Fatalf("missing definition %s", name)
	}
	return def
}

func requireNames(t *testing.T, c *Catalog, want ...types.TypeName) {
	t.Helper()
	got := c.Names()
	requireEqual(t, len(got), len(want))
	for i := range want {
		requireEqual(t, got[i], want[i])
	}
}

func requireDefinitions(t *testing.T, c *Catalog, want ...types.TypeName) {
	t.Helper()
	got := c.Definitions()
	requireEqual(t, len(got), len(want))
	for i := range want {
		requireEqual(t, got[i].Name(), want[i])
	}
}

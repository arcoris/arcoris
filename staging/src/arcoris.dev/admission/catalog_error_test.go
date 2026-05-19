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

package admission

import (
	"errors"
	"testing"
)

func TestNewCatalogRejectsNilRegistries(t *testing.T) {
	t.Parallel()

	kinds := NewBuiltinKindRegistry()
	reasons := NewBuiltinReasonRegistry()
	components := NewBuiltinComponentRegistry(kinds)

	tests := []struct {
		name       string
		reasons    *ReasonRegistry
		kinds      *KindRegistry
		components *ComponentRegistry
		want       error
	}{
		{
			name:       "nil reason registry",
			kinds:      kinds,
			components: components,
			want:       ErrNilReasonRegistry,
		},
		{
			name:       "nil kind registry",
			reasons:    reasons,
			components: components,
			want:       ErrNilKindRegistry,
		},
		{
			name:    "nil component registry",
			reasons: reasons,
			kinds:   kinds,
			want:    ErrNilComponentRegistry,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			catalog, err := NewCatalog(
				test.reasons,
				test.kinds,
				test.components,
			)
			if catalog != nil {
				t.Fatal("catalog should be nil on construction error")
			}
			if !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestNewCatalogRejectsMismatchedKindRegistry(t *testing.T) {
	t.Parallel()

	reasons := NewBuiltinReasonRegistry()
	kindsA := NewBuiltinKindRegistry()
	kindsB := NewBuiltinKindRegistry()
	components := NewBuiltinComponentRegistry(kindsA)

	catalog, err := NewCatalog(reasons, kindsB, components)
	if catalog != nil {
		t.Fatal("catalog should be nil on mismatched kind registry")
	}
	if !errors.Is(err, ErrMismatchedKindRegistry) {
		t.Fatalf("error = %v, want ErrMismatchedKindRegistry", err)
	}
}

func TestCatalogErrorsSupportIs(t *testing.T) {
	t.Parallel()

	if !errors.Is(ErrNilReasonRegistry, ErrNilReasonRegistry) {
		t.Fatal("nil reason registry sentinel should match itself")
	}
	if !errors.Is(ErrNilComponentRegistry, ErrNilComponentRegistry) {
		t.Fatal("nil component registry sentinel should match itself")
	}
	if !errors.Is(ErrMismatchedKindRegistry, ErrMismatchedKindRegistry) {
		t.Fatal("mismatched kind registry sentinel should match itself")
	}
}

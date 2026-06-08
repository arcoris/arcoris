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

package admissioncatalog

import (
	"errors"
	"testing"

	"arcoris.dev/admission"
)

func TestMergeEmptyList(t *testing.T) {
	catalog, err := Merge()
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}
	if catalog.LenReasons() != 0 || catalog.LenKinds() != 0 || catalog.LenComponents() != 0 {
		t.Fatal("empty merge returned non-empty catalog")
	}
}

func TestMergeCatalogs(t *testing.T) {
	first := mustCatalog(t, Input{
		Reasons: []ReasonDescriptor{reasonDescriptor(testReason)},
		Kinds:   []ComponentKindDescriptor{kindDescriptor(testKind)},
	})
	second := mustCatalog(t, Input{
		Reasons:    []ReasonDescriptor{reasonDescriptor(testOtherReason)},
		Kinds:      []ComponentKindDescriptor{kindDescriptor(testOtherKind)},
		Components: []ComponentDescriptor{componentDescriptor(testOtherComponent, testOtherKind)},
	})

	merged, err := Merge(first, second)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}
	if !merged.HasReason(testReason) || !merged.HasReason(testOtherReason) {
		t.Fatal("merged catalog is missing reasons")
	}
	if !merged.HasKind(testKind) || !merged.HasKind(testOtherKind) {
		t.Fatal("merged catalog is missing kinds")
	}
	if !merged.HasComponent(testOtherComponent) {
		t.Fatal("merged catalog is missing component")
	}
}

func TestMergeRejectsNilCatalog(t *testing.T) {
	catalog, err := Merge((*Catalog)(nil))
	if err == nil {
		t.Fatal("Merge returned nil error")
	}
	if catalog != nil {
		t.Fatal("Merge returned partial catalog")
	}
	typed := requireErrorIs[NilCatalogError](t, err, ErrNilCatalog)
	if typed.Index != 0 {
		t.Fatalf("Index = %d, want 0", typed.Index)
	}
	if typed.Operation != "merge" {
		t.Fatalf("Operation = %q, want merge", typed.Operation)
	}
	if got, want := typed.Path, "catalogs[0]"; got != want {
		t.Fatalf("Path = %q, want %q", got, want)
	}
}

func TestMergeRejectsDuplicates(t *testing.T) {
	tests := []struct {
		name   string
		first  *Catalog
		second *Catalog
		want   error
	}{
		{
			name:   "reason",
			first:  mustCatalog(t, Input{Reasons: []ReasonDescriptor{reasonDescriptor(testReason)}}),
			second: mustCatalog(t, Input{Reasons: []ReasonDescriptor{reasonDescriptor(testReason)}}),
			want:   ErrDuplicateReasonDeclaration,
		},
		{
			name:   "kind",
			first:  mustCatalog(t, Input{Kinds: []ComponentKindDescriptor{kindDescriptor(testKind)}}),
			second: mustCatalog(t, Input{Kinds: []ComponentKindDescriptor{kindDescriptor(testKind)}}),
			want:   ErrDuplicateComponentKindDeclaration,
		},
		{
			name: "component",
			first: mustCatalog(t, Input{
				Kinds:      []ComponentKindDescriptor{kindDescriptor(testKind)},
				Components: []ComponentDescriptor{componentDescriptor(testComponent, testKind)},
			}),
			second: mustCatalog(t, Input{
				Kinds:      []ComponentKindDescriptor{kindDescriptor(testOtherKind)},
				Components: []ComponentDescriptor{componentDescriptor(testComponent, testOtherKind)},
			}),
			want: ErrDuplicateComponentDeclaration,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Merge(tt.first, tt.second)
			if !errors.Is(err, tt.want) {
				t.Fatalf("Merge error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestMergeReportsIndexedDuplicatePaths(t *testing.T) {
	tests := []struct {
		name   string
		first  *Catalog
		second *Catalog
		want   error
		path   string
	}{
		{
			name:   "reason",
			first:  mustCatalog(t, Input{Reasons: []ReasonDescriptor{reasonDescriptor(testReason)}}),
			second: mustCatalog(t, Input{Reasons: []ReasonDescriptor{reasonDescriptor(testReason)}}),
			want:   ErrDuplicateReasonDeclaration,
			path:   "catalogs[1].reasons[0]",
		},
		{
			name:   "kind",
			first:  mustCatalog(t, Input{Kinds: []ComponentKindDescriptor{kindDescriptor(testKind)}}),
			second: mustCatalog(t, Input{Kinds: []ComponentKindDescriptor{kindDescriptor(testKind)}}),
			want:   ErrDuplicateComponentKindDeclaration,
			path:   "catalogs[1].kinds[0]",
		},
		{
			name: "component",
			first: mustCatalog(t, Input{
				Kinds:      []ComponentKindDescriptor{kindDescriptor(testKind)},
				Components: []ComponentDescriptor{componentDescriptor(testComponent, testKind)},
			}),
			second: mustCatalog(t, Input{
				Kinds:      []ComponentKindDescriptor{kindDescriptor(testOtherKind)},
				Components: []ComponentDescriptor{componentDescriptor(testComponent, testOtherKind)},
			}),
			want: ErrDuplicateComponentDeclaration,
			path: "catalogs[1].components[0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Merge(tt.first, tt.second)
			if !errors.Is(err, tt.want) {
				t.Fatalf("Merge error = %v, want %v", err, tt.want)
			}
			switch tt.want {
			case ErrDuplicateReasonDeclaration:
				typed := requireErrorIs[DuplicateReasonDeclarationError](t, err, tt.want)
				if typed.Path != tt.path {
					t.Fatalf("Path = %q, want %q", typed.Path, tt.path)
				}
			case ErrDuplicateComponentKindDeclaration:
				typed := requireErrorIs[DuplicateComponentKindDeclarationError](t, err, tt.want)
				if typed.Path != tt.path {
					t.Fatalf("Path = %q, want %q", typed.Path, tt.path)
				}
			case ErrDuplicateComponentDeclaration:
				typed := requireErrorIs[DuplicateComponentDeclarationError](t, err, tt.want)
				if typed.Path != tt.path {
					t.Fatalf("Path = %q, want %q", typed.Path, tt.path)
				}
			}
		})
	}
}

func TestMergeAllowsComponentsToReferenceKindsFromAnotherCatalog(t *testing.T) {
	kindsOnly := mustCatalog(t, Input{Kinds: []ComponentKindDescriptor{kindDescriptor(testKind)}})

	componentOnly := &Catalog{
		components: newComponentStore(),
	}
	componentOnly.components.declare(componentDescriptor(testComponent, testKind))

	merged, err := Merge(componentOnly, kindsOnly)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}
	if !merged.HasComponent(testComponent) {
		t.Fatal("merged catalog is missing cross-catalog component")
	}
}

func TestMergeRejectsUnknownComponentKindWithIndexedPath(t *testing.T) {
	componentOnly := &Catalog{
		components: newComponentStore(),
	}
	componentOnly.components.declare(componentDescriptor(testComponent, testKind))

	catalog, err := Merge(componentOnly)
	if err == nil {
		t.Fatal("Merge returned nil error")
	}
	if catalog != nil {
		t.Fatal("Merge returned partial catalog")
	}
	typed := requireErrorIs[UnknownComponentKindError](t, err, ErrUnknownComponentKind)
	if got, want := typed.Path, "catalogs[0].components[0]"; got != want {
		t.Fatalf("Path = %q, want %q", got, want)
	}
}

func TestMergeRejectsInvalidDescriptorWithIndexedPath(t *testing.T) {
	invalid := &Catalog{
		reasons: newReasonStore(),
	}
	invalid.reasons.declare(ReasonDescriptor{Reason: admission.Reason("bad-reason")})

	catalog, err := Merge(invalid)
	if err == nil {
		t.Fatal("Merge returned nil error")
	}
	if catalog != nil {
		t.Fatal("Merge returned partial catalog")
	}
	typed := requireErrorIs[InvalidReasonDescriptorError](t, err, ErrInvalidReasonDescriptor)
	if got, want := typed.Path, "catalogs[0].reasons[0]"; got != want {
		t.Fatalf("Path = %q, want %q", got, want)
	}
}

func TestMergeReturnsDetachedCatalog(t *testing.T) {
	first := mustCatalog(t, Input{Reasons: []ReasonDescriptor{reasonDescriptor(testReason)}})
	merged, err := Merge(first)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	mutated := first.Reasons()
	mutated[0] = reasonDescriptor(admission.Reason("mutated_reason"))
	if merged.HasReason(admission.Reason("mutated_reason")) {
		t.Fatal("merged catalog retained caller-mutated source slice")
	}
}

func TestMustMergePanicsOnInvalidInput(t *testing.T) {
	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("MustMerge did not panic")
		}
		err, ok := got.(error)
		if !ok {
			t.Fatalf("panic = %T, want error", got)
		}
		if !errors.Is(err, ErrNilCatalog) {
			t.Fatalf("panic error = %v, want ErrNilCatalog", err)
		}
	}()
	_ = MustMerge(nil)
}

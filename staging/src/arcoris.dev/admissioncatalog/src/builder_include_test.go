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

func TestBuilderInclude(t *testing.T) {
	included := mustCatalog(t, validInput())

	var builder Builder
	if err := builder.Include(included); err != nil {
		t.Fatalf("Include returned error: %v", err)
	}
	catalog, err := builder.Build()
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if !catalog.HasReason(testReason) || !catalog.HasKind(testKind) || !catalog.HasComponent(testComponent) {
		t.Fatal("included catalog declarations were not copied")
	}
}

func TestBuilderIncludeRejectsNilAndDuplicates(t *testing.T) {
	var builder Builder
	err := builder.Include(nil)
	if !errors.Is(err, ErrNilCatalog) {
		t.Fatalf("Include(nil) error = %v", err)
	}
	typed := requireErrorIs[NilCatalogError](t, err, ErrNilCatalog)
	if typed.Operation != "include" || typed.Index != -1 {
		t.Fatalf("nil catalog error = %+v, want include/-1", typed)
	}
	if typed.Path != "include" {
		t.Fatalf("Path = %q, want include", typed.Path)
	}

	catalog := mustCatalog(t, validInput())
	if err := builder.Include(catalog); err != nil {
		t.Fatalf("first Include returned error: %v", err)
	}
	err = builder.Include(catalog)
	if !errors.Is(err, ErrDuplicateReasonDeclaration) {
		t.Fatalf("second Include error = %v", err)
	}
	duplicate := requireErrorIs[DuplicateReasonDeclarationError](t, err, ErrDuplicateReasonDeclaration)
	if got, want := duplicate.Path, "include.reasons[0]"; got != want {
		t.Fatalf("Path = %q, want %q", got, want)
	}
}

func TestBuilderIncludeIsAllOrNothing(t *testing.T) {
	var builder Builder
	duplicateReason := admission.Reason("zz_reason")
	earlierReason := admission.Reason("aa_reason")
	if err := builder.DeclareReason(reasonDescriptor(duplicateReason)); err != nil {
		t.Fatalf("DeclareReason returned error: %v", err)
	}

	catalog := mustCatalog(t, Input{
		Reasons: []ReasonDescriptor{
			reasonDescriptor(earlierReason),
			reasonDescriptor(duplicateReason),
		},
	})
	if err := builder.Include(catalog); !errors.Is(err, ErrDuplicateReasonDeclaration) {
		t.Fatalf("Include error = %v, want ErrDuplicateReasonDeclaration", err)
	}

	built, err := builder.Build()
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if built.HasReason(earlierReason) {
		t.Fatal("failed Include partially mutated builder")
	}
	if !built.HasReason(duplicateReason) {
		t.Fatal("failed Include removed existing builder declaration")
	}
}

func TestBuilderIncludeEmptyCatalogIsNoOp(t *testing.T) {
	var builder Builder
	if err := builder.DeclareReason(reasonDescriptor(testReason)); err != nil {
		t.Fatalf("DeclareReason returned error: %v", err)
	}

	if err := builder.Include(&Catalog{}); err != nil {
		t.Fatalf("Include(empty) returned error: %v", err)
	}
	catalog, err := builder.Build()
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if !catalog.HasReason(testReason) || catalog.LenReasons() != 1 {
		t.Fatal("empty Include changed builder declarations")
	}
}

func TestBuilderIncludeAllowsLaterDeclarations(t *testing.T) {
	var builder Builder
	if err := builder.Include(mustCatalog(t, Input{Reasons: []ReasonDescriptor{reasonDescriptor(testReason)}})); err != nil {
		t.Fatalf("Include returned error: %v", err)
	}
	if err := builder.DeclareKind(kindDescriptor(testKind)); err != nil {
		t.Fatalf("DeclareKind after Include returned error: %v", err)
	}
	if err := builder.DeclareComponent(componentDescriptor(testComponent, testKind)); err != nil {
		t.Fatalf("DeclareComponent after Include returned error: %v", err)
	}
}

func TestBuilderIncludeReportsIndexedInvalidDescriptorPath(t *testing.T) {
	invalid := &Catalog{
		reasons: newReasonStore(),
	}
	invalid.reasons.declare(ReasonDescriptor{Reason: admission.Reason("bad-reason")})

	var builder Builder
	err := builder.Include(invalid)
	if !errors.Is(err, ErrInvalidReasonDescriptor) {
		t.Fatalf("Include error = %v, want ErrInvalidReasonDescriptor", err)
	}
	typed := requireErrorIs[InvalidReasonDescriptorError](t, err, ErrInvalidReasonDescriptor)
	if got, want := typed.Path, "include.reasons[0]"; got != want {
		t.Fatalf("Path = %q, want %q", got, want)
	}
}

func TestBuilderIncludeCopiesCatalogValues(t *testing.T) {
	source := mustCatalog(t, Input{Reasons: []ReasonDescriptor{reasonDescriptor(testReason)}})
	var builder Builder
	if err := builder.Include(source); err != nil {
		t.Fatalf("Include returned error: %v", err)
	}

	mutated := source.Reasons()
	mutated[0] = reasonDescriptor(admission.Reason("mutated_reason"))

	catalog, err := builder.Build()
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if catalog.HasReason(admission.Reason("mutated_reason")) {
		t.Fatal("builder retained caller-mutated descriptor slice")
	}
}

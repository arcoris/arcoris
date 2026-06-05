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
	if err := builder.Include(nil); !errors.Is(err, ErrNilCatalog) {
		t.Fatalf("Include(nil) error = %v", err)
	}

	catalog := mustCatalog(t, validInput())
	if err := builder.Include(catalog); err != nil {
		t.Fatalf("first Include returned error: %v", err)
	}
	if err := builder.Include(catalog); !errors.Is(err, ErrDuplicateReasonDeclaration) {
		t.Fatalf("second Include error = %v", err)
	}
}

func TestBuilderIncludeIsAllOrNothing(t *testing.T) {
	var builder Builder
	if err := builder.DeclareReason(reasonDescriptor(testReason)); err != nil {
		t.Fatalf("DeclareReason returned error: %v", err)
	}

	catalog := mustCatalog(t, Input{
		Reasons: []ReasonDescriptor{
			reasonDescriptor(testOtherReason),
			reasonDescriptor(testReason),
		},
	})
	if err := builder.Include(catalog); !errors.Is(err, ErrDuplicateReasonDeclaration) {
		t.Fatalf("Include error = %v, want ErrDuplicateReasonDeclaration", err)
	}

	built, err := builder.Build()
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if built.HasReason(testOtherReason) {
		t.Fatal("failed Include partially mutated builder")
	}
}

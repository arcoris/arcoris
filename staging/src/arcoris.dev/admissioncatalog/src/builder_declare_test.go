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

func TestBuilderDeclareRejectsDuplicateDeclarations(t *testing.T) {
	var builder Builder
	if err := builder.DeclareReason(reasonDescriptor(testReason)); err != nil {
		t.Fatalf("DeclareReason returned error: %v", err)
	}
	if err := builder.DeclareReason(reasonDescriptor(testReason)); !errors.Is(err, ErrDuplicateReasonDeclaration) {
		t.Fatalf("duplicate reason error = %v", err)
	}

	if err := builder.DeclareKind(kindDescriptor(testKind)); err != nil {
		t.Fatalf("DeclareKind returned error: %v", err)
	}
	if err := builder.DeclareKind(kindDescriptor(testKind)); !errors.Is(err, ErrDuplicateComponentKindDeclaration) {
		t.Fatalf("duplicate kind error = %v", err)
	}

	if err := builder.DeclareComponent(componentDescriptor(testComponent, testKind)); err != nil {
		t.Fatalf("DeclareComponent returned error: %v", err)
	}
	if err := builder.DeclareComponent(componentDescriptor(testComponent, testKind)); !errors.Is(err, ErrDuplicateComponentDeclaration) {
		t.Fatalf("duplicate component error = %v", err)
	}
}

func TestBuilderDeclareRejectsInvalidDescriptors(t *testing.T) {
	var builder Builder
	if err := builder.DeclareReason(ReasonDescriptor{Reason: admission.Reason("bad-reason")}); !errors.Is(err, ErrInvalidReasonDescriptor) {
		t.Fatalf("invalid reason error = %v", err)
	}
	if err := builder.DeclareKind(ComponentKindDescriptor{Kind: admission.ComponentKind("bad-kind")}); !errors.Is(err, ErrInvalidComponentKindDescriptor) {
		t.Fatalf("invalid kind error = %v", err)
	}
	if err := builder.DeclareComponent(ComponentDescriptor{ID: admission.ComponentID("bad id"), Kind: testKind}); !errors.Is(err, ErrInvalidComponentDescriptor) {
		t.Fatalf("invalid component error = %v", err)
	}
}

func TestBuilderDeclareRequiresKindBeforeComponent(t *testing.T) {
	var builder Builder
	err := builder.DeclareComponent(componentDescriptor(testComponent, testKind))
	if !errors.Is(err, ErrUnknownComponentKind) {
		t.Fatalf("DeclareComponent error = %v, want ErrUnknownComponentKind", err)
	}
	if err := builder.DeclareKind(kindDescriptor(testKind)); err != nil {
		t.Fatalf("DeclareKind returned error: %v", err)
	}
	if err := builder.DeclareComponent(componentDescriptor(testComponent, testKind)); err != nil {
		t.Fatalf("DeclareComponent returned error after kind declaration: %v", err)
	}
}

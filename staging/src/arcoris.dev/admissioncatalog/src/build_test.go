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

func TestBuildValidInput(t *testing.T) {
	catalog, err := Build(validInput())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if !catalog.HasReason(testReason) {
		t.Fatal("catalog does not contain reason")
	}
	if !catalog.HasKind(testKind) {
		t.Fatal("catalog does not contain kind")
	}
	if !catalog.HasComponent(testComponent) {
		t.Fatal("catalog does not contain component")
	}
}

func TestBuildRejectsInvalidDescriptors(t *testing.T) {
	tests := []struct {
		name  string
		input Input
		want  error
	}{
		{
			name:  "reason",
			input: Input{Reasons: []ReasonDescriptor{{Reason: admission.Reason("bad-reason")}}},
			want:  ErrInvalidReasonDescriptor,
		},
		{
			name:  "kind",
			input: Input{Kinds: []ComponentKindDescriptor{{Kind: admission.ComponentKind("bad-kind")}}},
			want:  ErrInvalidComponentKindDescriptor,
		},
		{
			name:  "component",
			input: Input{Components: []ComponentDescriptor{{ID: admission.ComponentID("bad id"), Kind: testKind}}},
			want:  ErrInvalidComponentDescriptor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catalog, err := Build(tt.input)
			if err == nil {
				t.Fatal("Build returned nil error")
			}
			if catalog != nil {
				t.Fatal("Build returned partial catalog")
			}
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestBuildRejectsDuplicateDeclarations(t *testing.T) {
	tests := []struct {
		name  string
		input Input
		want  error
	}{
		{
			name: "reason",
			input: Input{Reasons: []ReasonDescriptor{
				reasonDescriptor(testReason),
				reasonDescriptor(testReason),
			}},
			want: ErrDuplicateReasonDeclaration,
		},
		{
			name: "kind",
			input: Input{Kinds: []ComponentKindDescriptor{
				kindDescriptor(testKind),
				kindDescriptor(testKind),
			}},
			want: ErrDuplicateComponentKindDeclaration,
		},
		{
			name: "component",
			input: Input{
				Kinds: []ComponentKindDescriptor{kindDescriptor(testKind)},
				Components: []ComponentDescriptor{
					componentDescriptor(testComponent, testKind),
					componentDescriptor(testComponent, testKind),
				},
			},
			want: ErrDuplicateComponentDeclaration,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catalog, err := Build(tt.input)
			if err == nil {
				t.Fatal("Build returned nil error")
			}
			if catalog != nil {
				t.Fatal("Build returned partial catalog")
			}
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestBuildRejectsUnknownComponentKind(t *testing.T) {
	catalog, err := Build(Input{
		Components: []ComponentDescriptor{componentDescriptor(testComponent, testKind)},
	})
	if err == nil {
		t.Fatal("Build returned nil error")
	}
	if catalog != nil {
		t.Fatal("Build returned partial catalog")
	}
	typed := requireErrorIs[UnknownComponentKindError](t, err, ErrUnknownComponentKind)
	if typed.Kind != testKind {
		t.Fatalf("Kind = %s, want %s", typed.Kind, testKind)
	}
}

func TestBuildDetachesInputSlices(t *testing.T) {
	input := validInput()
	catalog := mustCatalog(t, input)

	input.Reasons[0] = reasonDescriptor(admission.Reason("mutated_reason"))
	input.Kinds[0] = kindDescriptor(admission.ComponentKind("mutated_kind"))
	input.Components[0] = componentDescriptor(admission.ComponentID("mutated.component"), testKind)

	if !catalog.HasReason(testReason) {
		t.Fatal("catalog reason changed after input mutation")
	}
	if !catalog.HasKind(testKind) {
		t.Fatal("catalog kind changed after input mutation")
	}
	if !catalog.HasComponent(testComponent) {
		t.Fatal("catalog component changed after input mutation")
	}
}

func TestMustBuildPanicsOnInvalidInput(t *testing.T) {
	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("MustBuild did not panic")
		}
		err, ok := got.(error)
		if !ok {
			t.Fatalf("panic = %T, want error", got)
		}
		if !errors.Is(err, ErrInvalidReasonDescriptor) {
			t.Fatalf("panic error = %v, want ErrInvalidReasonDescriptor", err)
		}
	}()
	_ = MustBuild(Input{Reasons: []ReasonDescriptor{{Reason: admission.Reason("bad-reason")}}})
}

func TestBuildErrorPaths(t *testing.T) {
	tests := []struct {
		name  string
		input Input
		want  error
		path  string
	}{
		{
			name:  "reason",
			input: Input{Reasons: []ReasonDescriptor{{Reason: admission.Reason("bad-reason")}}},
			want:  ErrInvalidReasonDescriptor,
			path:  "input.reasons[0]",
		},
		{
			name:  "kind",
			input: Input{Kinds: []ComponentKindDescriptor{{Kind: admission.ComponentKind("bad-kind")}}},
			want:  ErrInvalidComponentKindDescriptor,
			path:  "input.kinds[0]",
		},
		{
			name:  "component",
			input: Input{Components: []ComponentDescriptor{{ID: admission.ComponentID("bad id"), Kind: testKind}}},
			want:  ErrInvalidComponentDescriptor,
			path:  "input.components[0]",
		},
		{
			name:  "unknown component kind",
			input: Input{Components: []ComponentDescriptor{componentDescriptor(testComponent, testKind)}},
			want:  ErrUnknownComponentKind,
			path:  "input.components[0]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Build(tt.input)
			if err == nil {
				t.Fatal("Build returned nil error")
			}
			switch tt.want {
			case ErrInvalidReasonDescriptor:
				typed := requireErrorIs[InvalidReasonDescriptorError](t, err, tt.want)
				if typed.Path != tt.path {
					t.Fatalf("Path = %q, want %q", typed.Path, tt.path)
				}
			case ErrInvalidComponentKindDescriptor:
				typed := requireErrorIs[InvalidComponentKindDescriptorError](t, err, tt.want)
				if typed.Path != tt.path {
					t.Fatalf("Path = %q, want %q", typed.Path, tt.path)
				}
			case ErrInvalidComponentDescriptor:
				typed := requireErrorIs[InvalidComponentDescriptorError](t, err, tt.want)
				if typed.Path != tt.path {
					t.Fatalf("Path = %q, want %q", typed.Path, tt.path)
				}
			case ErrUnknownComponentKind:
				typed := requireErrorIs[UnknownComponentKindError](t, err, tt.want)
				if typed.Path != tt.path {
					t.Fatalf("Path = %q, want %q", typed.Path, tt.path)
				}
			}
		})
	}
}

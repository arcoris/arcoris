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

import "testing"

func TestInputZeroValueBuildsEmptyCatalog(t *testing.T) {
	catalog, err := Build(Input{})
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if catalog.LenReasons() != 0 || catalog.LenKinds() != 0 || catalog.LenComponents() != 0 {
		t.Fatalf("catalog lengths = %d/%d/%d, want empty", catalog.LenReasons(), catalog.LenKinds(), catalog.LenComponents())
	}
}

func TestInputComponentsMayAppearBeforeKindsInStructOrder(t *testing.T) {
	input := Input{
		Components: []ComponentDescriptor{componentDescriptor(testComponent, testKind)},
		Kinds:      []ComponentKindDescriptor{kindDescriptor(testKind)},
	}

	catalog, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if !catalog.HasComponent(testComponent) {
		t.Fatal("catalog is missing component")
	}
}

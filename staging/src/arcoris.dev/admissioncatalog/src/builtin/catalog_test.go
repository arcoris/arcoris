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

package builtin

import "testing"

func TestNewCatalog(t *testing.T) {
	catalog := NewCatalog()
	if catalog.LenReasons() == 0 {
		t.Fatal("catalog has no reasons")
	}
	if catalog.LenKinds() == 0 {
		t.Fatal("catalog has no kinds")
	}
	if catalog.LenComponents() == 0 {
		t.Fatal("catalog has no components")
	}
	if !catalog.HasComponent(ComponentResilienceRetryBudget) {
		t.Fatal("catalog does not contain retry budget component")
	}
}

func TestNewCatalogReturnsFreshCatalogs(t *testing.T) {
	first := NewCatalog()
	second := NewCatalog()
	if first == second {
		t.Fatal("NewCatalog returned the same catalog pointer")
	}

	firstReasons := first.Reasons()
	firstReasons[0].Reason = "mutated_reason"
	if second.HasReason("mutated_reason") {
		t.Fatal("catalogs share mutable reason state")
	}
}

func TestBuiltinComponentKindsAreDeclared(t *testing.T) {
	catalog := NewCatalog()
	for _, component := range ComponentDescriptors() {
		if !catalog.HasKind(component.Kind) {
			t.Fatalf("component %s references undeclared kind %s", component.ID, component.Kind)
		}
	}
}

func TestBuiltinCatalogDescriptorsAreValidMetadata(t *testing.T) {
	catalog := NewCatalog()
	for _, descriptor := range catalog.Reasons() {
		if !descriptor.IsValid() {
			t.Fatalf("invalid reason descriptor: %+v", descriptor)
		}
		if descriptor.Summary == "" {
			t.Fatalf("reason descriptor has empty summary: %+v", descriptor)
		}
		if !descriptor.DeclaredCapabilities.IsValid() {
			t.Fatalf("reason descriptor has invalid capabilities: %+v", descriptor)
		}
	}
	for _, descriptor := range catalog.Kinds() {
		if !descriptor.IsValid() {
			t.Fatalf("invalid kind descriptor: %+v", descriptor)
		}
		if descriptor.Summary == "" {
			t.Fatalf("kind descriptor has empty summary: %+v", descriptor)
		}
		if !descriptor.DeclaredCapabilities.IsValid() {
			t.Fatalf("kind descriptor has invalid capabilities: %+v", descriptor)
		}
	}
	for _, descriptor := range catalog.Components() {
		if !descriptor.IsValid() {
			t.Fatalf("invalid component descriptor: %+v", descriptor)
		}
		if descriptor.Summary == "" {
			t.Fatalf("component descriptor has empty summary: %+v", descriptor)
		}
		if !descriptor.DeclaredCapabilities.IsValid() {
			t.Fatalf("component descriptor has invalid capabilities: %+v", descriptor)
		}
	}
}

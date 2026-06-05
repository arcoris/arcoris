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
	"testing"

	"arcoris.dev/admission"
)

func TestCatalogAccessLookups(t *testing.T) {
	catalog := mustCatalog(t, validInput())

	if descriptor, ok := catalog.Reason(testReason); !ok || descriptor.Reason != testReason {
		t.Fatalf("Reason = %+v, %v", descriptor, ok)
	}
	if descriptor, ok := catalog.Kind(testKind); !ok || descriptor.Kind != testKind {
		t.Fatalf("Kind = %+v, %v", descriptor, ok)
	}
	if descriptor, ok := catalog.Component(testComponent); !ok || descriptor.ID != testComponent {
		t.Fatalf("Component = %+v, %v", descriptor, ok)
	}
}

func TestCatalogAccessInvalidAndMissingLookupsReturnFalse(t *testing.T) {
	catalog := mustCatalog(t, validInput())

	if _, ok := catalog.Reason(admission.Reason("bad-reason")); ok {
		t.Fatal("invalid reason lookup returned true")
	}
	if _, ok := catalog.Kind(admission.ComponentKind("bad-kind")); ok {
		t.Fatal("invalid kind lookup returned true")
	}
	if _, ok := catalog.Component(admission.ComponentID("bad id")); ok {
		t.Fatal("invalid component lookup returned true")
	}
	if catalog.HasReason(admission.Reason("missing_reason")) {
		t.Fatal("missing reason returned true")
	}
	if catalog.HasKind(admission.ComponentKind("missing_kind")) {
		t.Fatal("missing kind returned true")
	}
	if catalog.HasComponent(admission.ComponentID("missing.component")) {
		t.Fatal("missing component returned true")
	}
}

func TestCatalogAccessListsAreSortedAndDetached(t *testing.T) {
	catalog := mustCatalog(t, validInput())

	reasons := catalog.Reasons()
	if got, want := reasons[0].Reason, testReason; got != want {
		t.Fatalf("first reason = %s, want %s", got, want)
	}
	reasons[0] = reasonDescriptor(admission.Reason("mutated_reason"))
	if catalog.HasReason(admission.Reason("mutated_reason")) {
		t.Fatal("mutating returned reasons changed catalog")
	}

	kinds := catalog.Kinds()
	if got, want := kinds[0].Kind, testKind; got != want {
		t.Fatalf("first kind = %s, want %s", got, want)
	}
	kinds[0] = kindDescriptor(admission.ComponentKind("mutated_kind"))
	if catalog.HasKind(admission.ComponentKind("mutated_kind")) {
		t.Fatal("mutating returned kinds changed catalog")
	}

	components := catalog.Components()
	if got, want := components[0].ID, testComponent; got != want {
		t.Fatalf("first component = %s, want %s", got, want)
	}
	components[0] = componentDescriptor(admission.ComponentID("mutated.component"), testKind)
	if catalog.HasComponent(admission.ComponentID("mutated.component")) {
		t.Fatal("mutating returned components changed catalog")
	}
}

func TestCatalogAccessLengths(t *testing.T) {
	catalog := mustCatalog(t, validInput())

	if got, want := catalog.LenReasons(), len(validInput().Reasons); got != want {
		t.Fatalf("LenReasons = %d, want %d", got, want)
	}
	if got, want := catalog.LenKinds(), len(validInput().Kinds); got != want {
		t.Fatalf("LenKinds = %d, want %d", got, want)
	}
	if got, want := catalog.LenComponents(), len(validInput().Components); got != want {
		t.Fatalf("LenComponents = %d, want %d", got, want)
	}
}

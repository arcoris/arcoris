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
	"fmt"
	"sync"
	"testing"

	"arcoris.dev/admission"
)

func TestCatalogLookupMethods(t *testing.T) {
	t.Parallel()

	catalog := testCatalog()

	if got, ok := catalog.Reason(admission.ReasonDenied); !ok || got.Reason != admission.ReasonDenied {
		t.Fatalf("Reason = (%+v, %v), want descriptor,true", got, ok)
	}
	if got, ok := catalog.Kind(testKindBulkhead); !ok || got.Kind != testKindBulkhead {
		t.Fatalf("Kind = (%+v, %v), want descriptor,true", got, ok)
	}
	if got, ok := catalog.Component("resilience.bulkhead"); !ok || got.ID != "resilience.bulkhead" {
		t.Fatalf("Component = (%+v, %v), want descriptor,true", got, ok)
	}
}

func TestCatalogListMethodsReturnSortedCopies(t *testing.T) {
	t.Parallel()

	catalog := testCatalog()

	reasons := catalog.Reasons()
	if len(reasons) == 0 {
		t.Fatal("Reasons should not be empty")
	}
	for i := 1; i < len(reasons); i++ {
		if reasons[i-1].Reason.String() > reasons[i].Reason.String() {
			t.Fatalf("Reasons order[%d:%d] = %q,%q, want sorted",
				i-1,
				i,
				reasons[i-1].Reason,
				reasons[i].Reason,
			)
		}
	}
	reasons[0].Reason = "mutated_reason"
	if catalog.LenReasons() == 0 || catalog.Reasons()[0].Reason == "mutated_reason" {
		t.Fatal("mutating Reasons result should not mutate catalog")
	}

	kinds := catalog.Kinds()
	if len(kinds) == 0 {
		t.Fatal("Kinds should not be empty")
	}
	for i := 1; i < len(kinds); i++ {
		if kinds[i-1].Kind.String() > kinds[i].Kind.String() {
			t.Fatalf("Kinds order[%d:%d] = %q,%q, want sorted",
				i-1,
				i,
				kinds[i-1].Kind,
				kinds[i].Kind,
			)
		}
	}
	kinds[0].Kind = "mutated_kind"
	if catalog.LenKinds() == 0 || catalog.Kinds()[0].Kind == "mutated_kind" {
		t.Fatal("mutating Kinds result should not mutate catalog")
	}

	components := catalog.Components()
	if len(components) == 0 {
		t.Fatal("Components should not be empty")
	}
	for i := 1; i < len(components); i++ {
		if components[i-1].ID.String() > components[i].ID.String() {
			t.Fatalf("Components order[%d:%d] = %q,%q, want sorted",
				i-1,
				i,
				components[i-1].ID,
				components[i].ID,
			)
		}
	}
	components[0].ID = "resilience.mutated"
	if catalog.LenComponents() == 0 || catalog.Components()[0].ID == "resilience.mutated" {
		t.Fatal("mutating Components result should not mutate catalog")
	}
}

func TestCatalogConcurrentAccess(t *testing.T) {
	catalog := testCatalog()

	var wg sync.WaitGroup
	errCh := make(chan error, 32*3)
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			suffix := string(rune('a'+i/26)) + string(rune('a'+i%26))
			reason := admission.Reason("custom_reason_" + suffix)
			kind := admission.ComponentKind("custom_kind_" + suffix)
			componentID := admission.ComponentID("custom.component_" + suffix)

			if err := catalog.RegisterReason(testReasonDescriptor(reason)); err != nil {
				errCh <- fmt.Errorf("register reason %q: %w", reason, err)
			}
			if err := catalog.RegisterKind(testKindDescriptor(kind)); err != nil {
				errCh <- fmt.Errorf("register kind %q: %w", kind, err)
			}
			if err := catalog.RegisterComponent(testComponentDescriptor(componentID, kind)); err != nil {
				errCh <- fmt.Errorf("register component %q: %w", componentID, err)
			}

			_, _ = catalog.Reason(reason)
			_, _ = catalog.Kind(kind)
			_, _ = catalog.Component(componentID)
			_ = catalog.Reasons()
			_ = catalog.Kinds()
			_ = catalog.Components()
			_ = catalog.LenReasons()
			_ = catalog.LenKinds()
			_ = catalog.LenComponents()
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected concurrent catalog error: %v", err)
		}
	}
}

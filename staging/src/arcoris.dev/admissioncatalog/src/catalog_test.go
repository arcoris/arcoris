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
	"sync"
	"testing"
)

func TestCatalogZeroValueIsEmpty(t *testing.T) {
	var catalog Catalog
	if catalog.LenReasons() != 0 || catalog.LenKinds() != 0 || catalog.LenComponents() != 0 {
		t.Fatalf("zero catalog lengths = %d/%d/%d", catalog.LenReasons(), catalog.LenKinds(), catalog.LenComponents())
	}
	if _, ok := catalog.Reason(testReason); ok {
		t.Fatal("zero catalog returned reason")
	}
	if _, ok := catalog.Kind(testKind); ok {
		t.Fatal("zero catalog returned kind")
	}
	if _, ok := catalog.Component(testComponent); ok {
		t.Fatal("zero catalog returned component")
	}
	if reasons := catalog.Reasons(); len(reasons) != 0 {
		t.Fatalf("zero catalog Reasons length = %d, want 0", len(reasons))
	}
	if kinds := catalog.Kinds(); len(kinds) != 0 {
		t.Fatalf("zero catalog Kinds length = %d, want 0", len(kinds))
	}
	if components := catalog.Components(); len(components) != 0 {
		t.Fatalf("zero catalog Components length = %d, want 0", len(components))
	}
}

func TestCatalogConcurrentReads(t *testing.T) {
	catalog := mustCatalog(t, validInput())

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = catalog.HasReason(testReason)
				_ = catalog.HasKind(testKind)
				_ = catalog.HasComponent(testComponent)
				_ = catalog.Reasons()
				_ = catalog.Kinds()
				_ = catalog.Components()
			}
		}()
	}
	wg.Wait()
}

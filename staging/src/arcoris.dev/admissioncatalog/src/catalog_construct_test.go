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

func TestNewCatalog(t *testing.T) {
	t.Parallel()

	kinds := testKindRegistry()
	catalog, err := NewCatalog(
		testReasonRegistry(),
		kinds,
		testComponentRegistry(kinds),
	)
	if err != nil {
		t.Fatalf("NewCatalog returned error: %v", err)
	}
	if catalog.LenReasons() == 0 {
		t.Fatal("catalog should expose reason registry")
	}
	if catalog.LenKinds() == 0 {
		t.Fatal("catalog should expose kind registry")
	}
	if catalog.LenComponents() == 0 {
		t.Fatal("catalog should expose component registry")
	}
}

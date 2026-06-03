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

func TestCatalogLenMethods(t *testing.T) {
	t.Parallel()

	catalog := testCatalog()

	if got := catalog.LenReasons(); got != len(catalog.Reasons()) {
		t.Fatalf("LenReasons = %d, want %d", got, len(catalog.Reasons()))
	}
	if got := catalog.LenKinds(); got != len(catalog.Kinds()) {
		t.Fatalf("LenKinds = %d, want %d", got, len(catalog.Kinds()))
	}
	if got := catalog.LenComponents(); got != len(catalog.Components()) {
		t.Fatalf("LenComponents = %d, want %d", got, len(catalog.Components()))
	}
}

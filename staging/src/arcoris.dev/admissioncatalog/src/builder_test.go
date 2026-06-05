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

func TestBuilderZeroValueIsUsable(t *testing.T) {
	var builder Builder
	if err := builder.DeclareReason(reasonDescriptor(testReason)); err != nil {
		t.Fatalf("DeclareReason returned error: %v", err)
	}
	if err := builder.DeclareKind(kindDescriptor(testKind)); err != nil {
		t.Fatalf("DeclareKind returned error: %v", err)
	}
	if err := builder.DeclareComponent(componentDescriptor(testComponent, testKind)); err != nil {
		t.Fatalf("DeclareComponent returned error: %v", err)
	}

	catalog, err := builder.Build()
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if !catalog.HasReason(testReason) || !catalog.HasKind(testKind) || !catalog.HasComponent(testComponent) {
		t.Fatal("built catalog is missing declarations")
	}
}

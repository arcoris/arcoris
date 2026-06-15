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

package objectquery

import "testing"

func TestIdentityRequirementZeroValues(t *testing.T) {
	var namespace NamespaceRequirement
	var name NameRequirement

	if !namespace.IsZero() {
		t.Fatal("zero namespace requirement is not zero")
	}
	if !name.IsZero() {
		t.Fatal("zero name requirement is not zero")
	}
	if got, ok := namespace.Namespace(); ok || got != "" {
		t.Fatalf("Namespace() = %q, %v; want empty, false", got, ok)
	}
	if got, ok := name.Name(); ok || got != "" {
		t.Fatalf("Name() = %q, %v; want empty, false", got, ok)
	}
}

func TestNamespaceRequirementAccessorPresent(t *testing.T) {
	req, err := NamespaceEquals("system")
	requireNoError(t, err)

	got, ok := req.Namespace()
	if !ok || got != "system" {
		t.Fatalf("Namespace() = %q, %v; want system, true", got, ok)
	}
}

func TestNamespaceRequirementAccessorExplicitZeroNamespace(t *testing.T) {
	req, err := NamespaceEquals("")
	requireNoError(t, err)

	got, ok := req.Namespace()
	if !ok || got != "" {
		t.Fatalf("Namespace() = %q, %v; want empty, true", got, ok)
	}
}

func TestNameRequirementAccessorPresent(t *testing.T) {
	req, err := NameEquals("worker")
	requireNoError(t, err)

	got, ok := req.Name()
	if !ok || got != "worker" {
		t.Fatalf("Name() = %q, %v; want worker, true", got, ok)
	}
}

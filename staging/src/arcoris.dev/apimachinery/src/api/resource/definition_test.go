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

package resource

import (
	"testing"

	"arcoris.dev/apimachinery/api/identity"
)

func TestDefinitionAccessors(t *testing.T) {
	def := validDefinition()
	requireEqual(t, def.Group(), identity.Group("control.arcoris.dev"))
	requireEqual(t, def.Kind(), identity.Kind("Worker"))
	requireEqual(t, def.Resource(), identity.Resource("workers"))
	requireEqual(t, def.Scope(), ScopeNamespaced)
	requireEqual(t, def.GroupKind().String(), "control.arcoris.dev#Worker")
	requireEqual(t, def.GroupResource().String(), "control.arcoris.dev:workers")

	version, ok := def.Version(identity.Version("v1"))
	if !ok {
		t.Fatalf("Version(v1) not found")
	}
	if !version.Exposed() || !version.Canonical() {
		t.Fatalf("version flags not preserved")
	}
	if _, ok := version.Observed(); !ok {
		t.Fatalf("Observed() missing")
	}

	gvk, ok := def.GroupVersionKind(identity.Version("v1"))
	if !ok {
		t.Fatalf("GroupVersionKind(v1) ok = false")
	}
	requireEqual(t, gvk.String(), "control.arcoris.dev/v1#Worker")

	gvr, ok := def.GroupVersionResource(identity.Version("v1"))
	if !ok {
		t.Fatalf("GroupVersionResource(v1) ok = false")
	}
	requireEqual(t, gvr.String(), "control.arcoris.dev/v1:workers")

	if _, ok := def.GroupVersionKind(identity.Version("v2")); ok {
		t.Fatalf("unknown GVK version returned ok=true")
	}
	if _, ok := def.GroupVersionResource(identity.Version("v2")); ok {
		t.Fatalf("unknown GVR version returned ok=true")
	}
}

func TestDefinitionIsZero(t *testing.T) {
	var zero Definition
	if !zero.IsZero() {
		t.Fatalf("zero Definition IsZero() = false")
	}

	if validDefinition().IsZero() {
		t.Fatalf("non-zero Definition IsZero() = true")
	}
}

func TestDefinitionVersionsReturnsDetachedSlice(t *testing.T) {
	def := validDefinition()
	versions := def.Versions()
	if len(versions) != 1 {
		t.Fatalf("Versions len = %d, want 1", len(versions))
	}
	versions[0] = NewVersion(identity.Version("v2"), objectType())
	again := def.Versions()
	if again[0].Version() != identity.Version("v1") {
		t.Fatalf("Versions() did not detach returned slice")
	}
}

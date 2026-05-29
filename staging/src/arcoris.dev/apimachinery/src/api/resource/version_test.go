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

func TestVersionDefinitionAccessors(t *testing.T) {
	version := validVersion()

	requireEqual(t, version.Version(), identity.Version("v1"))

	if version.Desired().Code().String() != "object" {
		t.Fatalf("Desired() code = %s, want object", version.Desired().Code())
	}

	observed, ok := version.Observed()
	if !ok {
		t.Fatalf("Observed() ok = false")
	}
	requireEqual(t, observed.Code().String(), "object")

	if !version.Exposed() {
		t.Fatalf("Exposed() = false")
	}

	if !version.Canonical() {
		t.Fatalf("Canonical() = false")
	}
}

func TestVersionDefinitionObservedMissing(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		objectType(),
		Exposed(),
		Canonical(),
	)

	observed, ok := version.Observed()
	if ok {
		t.Fatalf("Observed() ok = true")
	}
	if !observed.IsZero() {
		t.Fatalf("missing Observed() returned non-zero type")
	}
}

func TestVersionDefinitionIsZero(t *testing.T) {
	var zero VersionDefinition
	if !zero.IsZero() {
		t.Fatalf("zero VersionDefinition IsZero() = false")
	}

	if validVersion().IsZero() {
		t.Fatalf("non-zero VersionDefinition IsZero() = true")
	}
}

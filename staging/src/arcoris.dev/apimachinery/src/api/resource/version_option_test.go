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

func TestVersionOptions(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		objectType(),
		Observed(objectType()),
		Exposed(),
		Canonical(),
	)

	if !version.Exposed() {
		t.Fatalf("Exposed option was not applied")
	}

	if !version.Canonical() {
		t.Fatalf("Canonical option was not applied")
	}

	if _, ok := version.Observed(); !ok {
		t.Fatalf("Observed option was not applied")
	}
}

func TestNewVersionIgnoresNilOption(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		objectType(),
		nil,
	)

	requireEqual(t, version.Version(), identity.Version("v1"))
}

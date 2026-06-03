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

package objectownership

import "testing"

func TestVersionString(t *testing.T) {
	if VersionV1.String() != "v1" {
		t.Fatalf("VersionV1.String() = %q", VersionV1.String())
	}
}

func TestVersionIsZero(t *testing.T) {
	if !Version("").IsZero() {
		t.Fatalf("zero version IsZero() = false")
	}
	if VersionV1.IsZero() {
		t.Fatalf("VersionV1 IsZero() = true")
	}
}

func TestVersionIsSupported(t *testing.T) {
	if !VersionV1.IsSupported() {
		t.Fatalf("VersionV1 IsSupported() = false")
	}
	if Version("v2").IsSupported() {
		t.Fatalf("unknown version IsSupported() = true")
	}
}

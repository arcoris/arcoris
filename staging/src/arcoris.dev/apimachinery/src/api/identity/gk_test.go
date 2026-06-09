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

package identity

import "testing"

func TestGroupKindValue(t *testing.T) {
	gk := GroupKind{Group: "control.arcoris.dev", Kind: "Worker"}
	requireCanonicalText(t, gk, "control.arcoris.dev#Worker")

	if !(GroupKind{}).IsZero() {
		t.Fatalf("zero GroupKind should be zero")
	}
	if gk.IsZero() {
		t.Fatalf("complete GroupKind should not be zero")
	}
}

func TestGroupKindComposition(t *testing.T) {
	gk := GroupKind{Group: "control.arcoris.dev", Kind: "Worker"}
	want := GroupVersionKind{Group: "control.arcoris.dev", Version: "v1", Kind: "Worker"}

	requireEqual(t, "WithVersion()", gk.WithVersion("v1"), want)
}

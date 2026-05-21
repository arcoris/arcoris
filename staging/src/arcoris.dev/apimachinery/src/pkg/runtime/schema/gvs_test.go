/*
   Copyright 2026 The ARCORIS Authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package schema

import "testing"

// GroupVersions tests document that the receiver order is authoritative for
// preference matching and that no legacy best-match ambiguity is implemented.

// TestGroupVersionsValidateAndStrings verifies validation and ordered display.
func TestGroupVersionsValidateAndStrings(t *testing.T) {
	gvs := GroupVersions{
		{Group: "control.arcoris.dev", Version: "v1"},
		{Group: "control.arcoris.dev", Version: "v1beta1"},
	}
	if err := gvs.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	want := "[control.arcoris.dev/v1, control.arcoris.dev/v1beta1]"
	if gvs.Identifier() != want || gvs.String() != want {
		t.Fatalf("Identifier/String = %q/%q", gvs.Identifier(), gvs.String())
	}
	if (GroupVersions{}).Identifier() != "[]" {
		t.Fatalf("empty Identifier = %q", (GroupVersions{}).Identifier())
	}
	if err := (GroupVersions{{Version: "v01"}}).Validate(); err == nil {
		t.Fatalf("Validate on invalid entry expected error")
	}
}

// TestGroupVersionsContains verifies exact structural membership checks.
func TestGroupVersionsContains(t *testing.T) {
	gvs := GroupVersions{{Version: "v1"}, {Group: "control.arcoris.dev", Version: "v1alpha1"}}
	if !gvs.Contains(GroupVersion{Version: "v1"}) {
		t.Fatalf("Contains(v1) = false")
	}
	if gvs.Contains(GroupVersion{Group: "other.arcoris.dev", Version: "v1"}) {
		t.Fatalf("Contains(other/v1) = true")
	}
}

// TestGroupVersionsKindPreferenceOrderIsAuthoritative verifies preferred versions win first.
func TestGroupVersionsKindPreferenceOrderIsAuthoritative(t *testing.T) {
	preferred := GroupVersions{
		{Group: "control.arcoris.dev", Version: "v1"},
		{Group: "control.arcoris.dev", Version: "v1beta1"},
	}
	candidates := []GroupVersionKind{
		{Group: "control.arcoris.dev", Version: "v1beta1", Kind: "WorkloadClass"},
		{Group: "control.arcoris.dev", Version: "v1", Kind: "WorkloadClass"},
	}
	got, ok := preferred.KindForGroupVersionKinds(candidates)
	if !ok {
		t.Fatalf("KindForGroupVersionKinds did not find match")
	}
	want := GroupVersionKind{Group: "control.arcoris.dev", Version: "v1", Kind: "WorkloadClass"}
	if got != want {
		t.Fatalf("KindForGroupVersionKinds = %+v, want %+v", got, want)
	}
}

// TestGroupVersionsResourcePreferenceOrderIsAuthoritative verifies resource matching uses preference order.
func TestGroupVersionsResourcePreferenceOrderIsAuthoritative(t *testing.T) {
	preferred := GroupVersions{
		{Group: "control.arcoris.dev", Version: "v1"},
		{Group: "control.arcoris.dev", Version: "v1beta1"},
	}
	candidates := []GroupVersionResource{
		{Group: "control.arcoris.dev", Version: "v1beta1", Resource: "workloadclasses"},
		{Group: "control.arcoris.dev", Version: "v1", Resource: "workloadclasses"},
	}
	got, ok := preferred.ResourceForGroupVersionResources(candidates)
	if !ok {
		t.Fatalf("ResourceForGroupVersionResources did not find match")
	}
	want := GroupVersionResource{Group: "control.arcoris.dev", Version: "v1", Resource: "workloadclasses"}
	if got != want {
		t.Fatalf("ResourceForGroupVersionResources = %+v, want %+v", got, want)
	}
}

// TestGroupVersionsNoMatch verifies unmatched candidate lists return false and zero values.
func TestGroupVersionsNoMatch(t *testing.T) {
	preferred := GroupVersions{{Group: "control.arcoris.dev", Version: "v1"}}
	if got, ok := preferred.KindForGroupVersionKinds([]GroupVersionKind{{Group: "other.arcoris.dev", Version: "v1", Kind: "Other"}}); ok || got != (GroupVersionKind{}) {
		t.Fatalf("KindForGroupVersionKinds = %+v, %v", got, ok)
	}
	if got, ok := preferred.ResourceForGroupVersionResources([]GroupVersionResource{{Group: "other.arcoris.dev", Version: "v1", Resource: "others"}}); ok || got != (GroupVersionResource{}) {
		t.Fatalf("ResourceForGroupVersionResources = %+v, %v", got, ok)
	}
}

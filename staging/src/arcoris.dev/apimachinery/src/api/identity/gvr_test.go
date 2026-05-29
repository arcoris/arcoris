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

func TestGroupVersionResourceValue(t *testing.T) {
	gvr := GroupVersionResource{Group: "control.arcoris.dev", Version: "v1", Resource: "workers"}
	requireIdentifier(t, gvr, "control.arcoris.dev/v1:workers")

	if !(GroupVersionResource{}).IsZero() {
		t.Fatalf("zero GroupVersionResource should be zero")
	}
	if gvr.IsZero() {
		t.Fatalf("complete GroupVersionResource should not be zero")
	}
}

func TestGroupVersionResourceParts(t *testing.T) {
	gvr := GroupVersionResource{Group: "control.arcoris.dev", Version: "v1", Resource: "workers"}

	requireEqual(
		t,
		"GroupVersion()",
		gvr.GroupVersion(),
		GroupVersion{Group: "control.arcoris.dev", Version: "v1"},
	)

	requireEqual(
		t,
		"GroupResource()",
		gvr.GroupResource(),
		GroupResource{Group: "control.arcoris.dev", Resource: "workers"},
	)
}

func TestGroupVersionResourceComposition(t *testing.T) {
	gvr := GroupVersionResource{Group: "control.arcoris.dev", Version: "v1", Resource: "workers"}
	want := GroupVersionResourcePath{
		Group:       "control.arcoris.dev",
		Version:     "v1",
		Resource:    "workers",
		Subresource: "status",
	}

	requireEqual(t, "WithSubresource()", gvr.WithSubresource("status"), want)
}

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

func TestGroupVersionResourcePathValue(t *testing.T) {
	path := GroupVersionResourcePath{
		Group:       "control.arcoris.dev",
		Version:     "v1",
		Resource:    "workers",
		Subresource: "status",
	}
	requireCanonicalText(t, path, "control.arcoris.dev/v1:workers/status")

	if !(GroupVersionResourcePath{}).IsZero() {
		t.Fatalf("zero GroupVersionResourcePath should be zero")
	}
	if path.IsZero() {
		t.Fatalf("complete GroupVersionResourcePath should not be zero")
	}
	if !path.HasSubresource() {
		t.Fatalf("HasSubresource() = false, want true")
	}
}

func TestGroupVersionResourcePathParts(t *testing.T) {
	path := GroupVersionResourcePath{
		Group:       "control.arcoris.dev",
		Version:     "v1",
		Resource:    "workers",
		Subresource: "status",
	}

	requireEqual(
		t,
		"GroupVersion()",
		path.GroupVersion(),
		GroupVersion{Group: "control.arcoris.dev", Version: "v1"},
	)

	requireEqual(
		t,
		"GroupVersionResource()",
		path.GroupVersionResource(),
		GroupVersionResource{Group: "control.arcoris.dev", Version: "v1", Resource: "workers"},
	)

	requireEqual(
		t,
		"GroupResource()",
		path.GroupResource(),
		GroupResource{Group: "control.arcoris.dev", Resource: "workers"},
	)

	requireEqual(
		t,
		"ResourcePath()",
		path.ResourcePath(),
		ResourcePath{Resource: "workers", Subresource: "status"},
	)
}

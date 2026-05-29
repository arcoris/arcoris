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

// GroupResource identifies an API resource collection within a group.
//
// The core group uses only the resource text, for example "pods". A named
// group uses "group:resource", for example "control.arcoris.dev:workers".
type GroupResource struct {
	Group    Group
	Resource Resource
}

// String returns the canonical group/resource text without revalidating it.
func (gr GroupResource) String() string {
	return joinGroupResource(gr.Group, gr.Resource)
}

// Identifier returns the canonical group/resource identity string.
func (gr GroupResource) Identifier() string {
	return gr.String()
}

// IsZero reports whether group and resource are both absent.
func (gr GroupResource) IsZero() bool {
	return gr.Group.IsZero() && gr.Resource.IsZero()
}

// WithVersion composes this group/resource with a version.
func (gr GroupResource) WithVersion(version Version) GroupVersionResource {
	return GroupVersionResource{Group: gr.Group, Version: version, Resource: gr.Resource}
}

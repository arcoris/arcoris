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

// GroupVersionResource identifies a concrete versioned API resource collection.
//
// The canonical form is GroupVersion + ":" + Resource, for example "v1:pods"
// or "control.arcoris.dev/v1:workers". This identity names a collection and
// never includes a subresource.
type GroupVersionResource struct {
	Group    Group
	Version  Version
	Resource Resource
}

// String returns the canonical group/version/resource text without revalidating it.
func (gvr GroupVersionResource) String() string {
	return joinGroupVersionResource(gvr.GroupVersion(), gvr.Resource)
}

// CanonicalText validates the group/version/resource identity and returns its canonical text.
func (gvr GroupVersionResource) CanonicalText() (string, error) {
	if err := gvr.Validate(); err != nil {
		return "", err
	}

	return gvr.String(), nil
}

// IsZero reports whether group, version, and resource are all absent.
func (gvr GroupVersionResource) IsZero() bool {
	return gvr.Group.IsZero() &&
		gvr.Version.IsZero() &&
		gvr.Resource.IsZero()
}

// GroupVersion returns the group/version portion of the identity.
func (gvr GroupVersionResource) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvr.Group, Version: gvr.Version}
}

// GroupResource returns the group/resource portion of the identity.
func (gvr GroupVersionResource) GroupResource() GroupResource {
	return GroupResource{Group: gvr.Group, Resource: gvr.Resource}
}

// WithSubresource composes this collection identity with a subresource.
func (gvr GroupVersionResource) WithSubresource(subresource Subresource) GroupVersionResourcePath {
	return GroupVersionResourcePath{
		Group:       gvr.Group,
		Version:     gvr.Version,
		Resource:    gvr.Resource,
		Subresource: subresource,
	}
}

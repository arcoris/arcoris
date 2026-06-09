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

// GroupVersionResourcePath identifies a versioned resource and optional subresource.
//
// The canonical form is "group/version:resource" or
// "group/version:resource/subresource"; the core group omits the group prefix.
// URL-like "group/version/resource" forms are intentionally rejected because
// identity grammar is not REST routing syntax.
type GroupVersionResourcePath struct {
	Group       Group
	Version     Version
	Resource    Resource
	Subresource Subresource
}

// String returns the canonical versioned resource path without revalidating it.
func (gvrp GroupVersionResourcePath) String() string {
	return joinGroupVersionResourcePath(gvrp.GroupVersion(), gvrp.ResourcePath())
}

// CanonicalText validates the versioned resource path and returns its canonical text.
func (gvrp GroupVersionResourcePath) CanonicalText() (string, error) {
	if err := gvrp.Validate(); err != nil {
		return "", err
	}

	return gvrp.String(), nil
}

// IsZero reports whether all identity fields are absent.
func (gvrp GroupVersionResourcePath) IsZero() bool {
	return gvrp.Group.IsZero() &&
		gvrp.Version.IsZero() &&
		gvrp.Resource.IsZero() &&
		gvrp.Subresource.IsZero()
}

// GroupVersion returns the group/version portion of the identity.
func (gvrp GroupVersionResourcePath) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvrp.Group, Version: gvrp.Version}
}

// GroupVersionResource returns the collection portion of the identity.
func (gvrp GroupVersionResourcePath) GroupVersionResource() GroupVersionResource {
	return GroupVersionResource{
		Group:    gvrp.Group,
		Version:  gvrp.Version,
		Resource: gvrp.Resource,
	}
}

// GroupResource returns the group/resource portion of the identity.
func (gvrp GroupVersionResourcePath) GroupResource() GroupResource {
	return GroupResource{Group: gvrp.Group, Resource: gvrp.Resource}
}

// ResourcePath returns the resource/subresource portion of the identity.
func (gvrp GroupVersionResourcePath) ResourcePath() ResourcePath {
	return ResourcePath{Resource: gvrp.Resource, Subresource: gvrp.Subresource}
}

// HasSubresource reports whether the identity includes a subresource segment.
func (gvrp GroupVersionResourcePath) HasSubresource() bool {
	return !gvrp.Subresource.IsZero()
}

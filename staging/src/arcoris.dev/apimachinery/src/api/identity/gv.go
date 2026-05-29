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

// GroupVersion identifies a version within an API group.
//
// The core group uses only the version text, for example "v1". A named group
// uses "group/version", for example "control.arcoris.dev/v1".
type GroupVersion struct {
	Group   Group
	Version Version
}

// String returns the canonical group/version text without revalidating it.
func (gv GroupVersion) String() string {
	return joinGroupVersion(gv.Group, gv.Version)
}

// Identifier returns the canonical group/version identity string.
//
// It is equivalent to String and is intended for diagnostics and map keys.
func (gv GroupVersion) Identifier() string {
	return gv.String()
}

// IsZero reports whether group and version are both absent.
//
// The zero value is useful as an optional sentinel but is invalid as a complete
// group/version identity.
func (gv GroupVersion) IsZero() bool {
	return gv.Group.IsZero() && gv.Version.IsZero()
}

// WithKind composes this group/version with a kind.
//
// The method does not validate either side; callers can validate the completed
// GroupVersionKind at trust boundaries.
func (gv GroupVersion) WithKind(kind Kind) GroupVersionKind {
	return GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: kind}
}

// WithResource composes this group/version with a resource collection.
//
// The method preserves the exact fields without pluralization or route policy.
func (gv GroupVersion) WithResource(resource Resource) GroupVersionResource {
	return GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource}
}

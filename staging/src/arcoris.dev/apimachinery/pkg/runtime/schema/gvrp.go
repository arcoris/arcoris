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

import "strings"

// GroupVersionResourcePath identifies a versioned resource and optional subresource.
//
// The canonical form is "group/version:resource" or
// "group/version:resource/subresource"; the core group omits the group prefix.
// The resource segment remains a collection name, and the optional subresource
// is represented as a separate field and path suffix. This identity is the
// schema-level form a future REST mapper can lower into actual request paths.
type GroupVersionResourcePath struct {
	Group       Group
	Version     Version
	Resource    Resource
	Subresource Subresource
}

// ParseGroupVersionResourcePath parses a canonical versioned resource path.
//
// The parser requires a colon between group/version and resource path and then
// applies the strict ResourcePath parser to the right-hand side. URL-like
// "group/version/resource" input is rejected.
func ParseGroupVersionResourcePath(val string) (GroupVersionResourcePath, error) {
	if val == "" {
		return GroupVersionResourcePath{}, invalid("group/version/resource path", val, "group/version and resource are required")
	}
	gvPart, pathPart, ok := strings.Cut(val, ":")
	if !ok || gvPart == "" || pathPart == "" || strings.Contains(pathPart, ":") {
		return GroupVersionResourcePath{}, invalid("group/version/resource path", val, "expected canonical form group/version:resource[/subresource]")
	}
	gv, err := ParseGroupVersion(gvPart)
	if err != nil {
		return GroupVersionResourcePath{}, invalidValue("group/version/resource path", val, err)
	}
	path, err := ParseResourcePath(pathPart)
	if err != nil {
		return GroupVersionResourcePath{}, invalidValue("group/version/resource path", val, err)
	}
	return GroupVersionResourcePath{
		Group:       gv.Group,
		Version:     gv.Version,
		Resource:    path.Resource,
		Subresource: path.Subresource,
	}, nil
}

// String returns the canonical group/version/resource path string without revalidating it.
//
// The subresource slash is emitted only when Subresource is non-empty.
func (gvrp GroupVersionResourcePath) String() string {
	return gvrp.GroupVersion().String() + ":" + gvrp.ResourcePath().String()
}

// Identifier returns the canonical group/version/resource path identity string.
//
// It is equivalent to String and suitable for stable diagnostics and map keys.
func (gvrp GroupVersionResourcePath) Identifier() string {
	return gvrp.String()
}

// Validate checks that the versioned resource path is a complete canonical identity.
//
// Group may be empty for the core group. Version and Resource are required.
// Subresource may be empty, but if present it must be a DNS-1123 single label.
func (gvrp GroupVersionResourcePath) Validate() error {
	if gvrp.Version.IsZero() {
		return invalid("group/version/resource path", gvrp.String(), "version is required")
	}
	if gvrp.Resource.IsZero() {
		return invalid("group/version/resource path", gvrp.String(), "resource is required")
	}
	if err := gvrp.GroupVersion().Validate(); err != nil {
		return invalidValue("group/version/resource path", gvrp.String(), err)
	}
	if err := gvrp.ResourcePath().Validate(); err != nil {
		return invalidValue("group/version/resource path", gvrp.String(), err)
	}
	return nil
}

// IsZero reports whether all identity fields are empty.
//
// It is an optional sentinel check only; complete versioned resource paths
// require at least Version and Resource.
func (gvrp GroupVersionResourcePath) IsZero() bool {
	return gvrp.Group.IsZero() &&
		gvrp.Version.IsZero() &&
		gvrp.Resource.IsZero() &&
		gvrp.Subresource.IsZero()
}

// GroupVersion returns the group/version portion of the identity.
//
// The returned value preserves the exact group and version fields.
func (gvrp GroupVersionResourcePath) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvrp.Group, Version: gvrp.Version}
}

// GroupVersionResource returns the group/version/resource portion of the identity.
//
// The subresource field is intentionally dropped so callers can address the
// collection identity separately from the subresource path.
func (gvrp GroupVersionResourcePath) GroupVersionResource() GroupVersionResource {
	return GroupVersionResource{Group: gvrp.Group, Version: gvrp.Version, Resource: gvrp.Resource}
}

// GroupResource returns the group/resource portion of the identity.
//
// This supports comparisons where version preference has not yet been applied.
func (gvrp GroupVersionResourcePath) GroupResource() GroupResource {
	return GroupResource{Group: gvrp.Group, Resource: gvrp.Resource}
}

// ResourcePath returns the resource/subresource portion of the identity.
//
// The returned path preserves the structured subresource instead of reparsing
// the canonical string.
func (gvrp GroupVersionResourcePath) ResourcePath() ResourcePath {
	return ResourcePath{Resource: gvrp.Resource, Subresource: gvrp.Subresource}
}

// HasSubresource reports whether the identity includes a subresource.
//
// It is based on the Subresource field and does not inspect String output.
func (gvrp GroupVersionResourcePath) HasSubresource() bool {
	return !gvrp.Subresource.IsZero()
}

// MarshalText returns the canonical versioned resource path after validation.
//
// Invalid direct literals and partial identities are rejected before
// serialization.
func (gvrp GroupVersionResourcePath) MarshalText() ([]byte, error) {
	if err := gvrp.Validate(); err != nil {
		return nil, err
	}
	return []byte(gvrp.String()), nil
}

// UnmarshalText parses a canonical versioned resource path.
//
// The method rejects URL-like alternatives and resource paths with multiple
// subresource segments.
func (gvrp *GroupVersionResourcePath) UnmarshalText(data []byte) error {
	if gvrp == nil {
		return nilUnmarshalReceiver("group/version/resource path")
	}
	parsed, err := ParseGroupVersionResourcePath(string(data))
	if err != nil {
		return err
	}
	*gvrp = parsed
	return nil
}

// MarshalJSON returns the canonical versioned resource path as a JSON scalar.
//
// It never emits an object with group, version, resource, and subresource
// fields.
func (gvrp GroupVersionResourcePath) MarshalJSON() ([]byte, error) {
	return marshalJSONString("group/version/resource path", gvrp.String(), gvrp.Validate)
}

// UnmarshalJSON parses a canonical versioned resource path from a JSON scalar.
//
// Non-string JSON, including null, is rejected before parsing.
func (gvrp *GroupVersionResourcePath) UnmarshalJSON(data []byte) error {
	if gvrp == nil {
		return nilUnmarshalReceiver("group/version/resource path")
	}
	val, err := unmarshalJSONString("group/version/resource path", data)
	if err != nil {
		return err
	}
	return gvrp.UnmarshalText([]byte(val))
}

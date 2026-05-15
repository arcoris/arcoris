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

// GroupVersionResource identifies a concrete versioned API resource collection.
//
// The canonical form is "group/version:resource" or "version:resource" for the
// core group. This identity names a collection and never includes a
// subresource. It is a resource identity, not an object type identity.
type GroupVersionResource struct {
	Group    Group
	Version  Version
	Resource Resource
}

// ParseGroupVersionResource parses a canonical group/version/resource identity.
//
// Legacy dotted triplets and path-like forms are rejected because they are not
// canonical ARCORIS schema identifiers.
func ParseGroupVersionResource(val string) (GroupVersionResource, error) {
	if val == "" {
		return GroupVersionResource{}, invalid("group/version/resource", val, "group/version and resource are required")
	}
	gvPart, resourcePart, ok := strings.Cut(val, ":")
	if !ok || gvPart == "" || resourcePart == "" || strings.Contains(resourcePart, ":") {
		return GroupVersionResource{}, invalid("group/version/resource", val, "expected canonical form group/version:resource")
	}
	gv, err := ParseGroupVersion(gvPart)
	if err != nil {
		return GroupVersionResource{}, invalidValue("group/version/resource", val, err)
	}
	resource, err := ParseResource(resourcePart)
	if err != nil {
		return GroupVersionResource{}, invalidValue("group/version/resource", val, err)
	}
	return GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource}, nil
}

// String returns the canonical group/version/resource string without revalidating it.
//
// The colon separates the group/version namespace from the resource collection
// so the form cannot be confused with API URL paths.
func (gvr GroupVersionResource) String() string {
	return gvr.GroupVersion().String() + ":" + gvr.Resource.String()
}

// Identifier returns the canonical group/version/resource identity string.
//
// It is equivalent to String and suitable for stable diagnostics and map keys.
func (gvr GroupVersionResource) Identifier() string {
	return gvr.String()
}

// Validate checks that the group/version/resource is a complete canonical identity.
//
// Group may be empty for the core group. Version and Resource are required, so
// zero and partial GVR values are rejected.
func (gvr GroupVersionResource) Validate() error {
	if gvr.Version.IsZero() {
		return invalid("group/version/resource", gvr.String(), "version is required")
	}
	if gvr.Resource.IsZero() {
		return invalid("group/version/resource", gvr.String(), "resource is required")
	}
	if err := gvr.GroupVersion().Validate(); err != nil {
		return invalidValue("group/version/resource", gvr.String(), err)
	}
	if err := gvr.Resource.Validate(); err != nil {
		return invalidValue("group/version/resource", gvr.String(), err)
	}
	return nil
}

// IsZero reports whether group, version, and resource are all empty.
//
// It is intended for optional fields only and does not imply the value is valid
// for routing or serialization.
func (gvr GroupVersionResource) IsZero() bool {
	return gvr.Group.IsZero() && gvr.Version.IsZero() && gvr.Resource.IsZero()
}

// GroupVersion returns the group/version portion of the identity.
//
// The returned value preserves the exact group and version fields without
// validation or normalization.
func (gvr GroupVersionResource) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvr.Group, Version: gvr.Version}
}

// GroupResource returns the group/resource portion of the identity.
//
// This is useful for resource comparisons before version preference has been
// applied.
func (gvr GroupVersionResource) GroupResource() GroupResource {
	return GroupResource{Group: gvr.Group, Resource: gvr.Resource}
}

// WithSubresource returns the GroupVersionResourcePath formed from this resource and subresource.
//
// The method composes fields without validation so callers can attach the
// subresource first and validate the complete path later.
func (gvr GroupVersionResource) WithSubresource(subresource Subresource) GroupVersionResourcePath {
	return GroupVersionResourcePath{
		Group:       gvr.Group,
		Version:     gvr.Version,
		Resource:    gvr.Resource,
		Subresource: subresource,
	}
}

// MarshalText returns the canonical group/version/resource string after validation.
//
// Invalid direct literals and partial identities are rejected before
// serialization.
func (gvr GroupVersionResource) MarshalText() ([]byte, error) {
	if err := gvr.Validate(); err != nil {
		return nil, err
	}
	return []byte(gvr.String()), nil
}

// UnmarshalText parses a canonical group/version/resource string.
//
// The method rejects legacy dotted triplets, URL-like paths, and comma-based
// resource descriptions.
func (gvr *GroupVersionResource) UnmarshalText(data []byte) error {
	if gvr == nil {
		return nilUnmarshalReceiver("group/version/resource")
	}
	parsed, err := ParseGroupVersionResource(string(data))
	if err != nil {
		return err
	}
	*gvr = parsed
	return nil
}

// MarshalJSON returns the canonical group/version/resource string as a JSON scalar.
//
// It never emits an object with group, version, and resource fields.
func (gvr GroupVersionResource) MarshalJSON() ([]byte, error) {
	return marshalJSONString("group/version/resource", gvr.String(), gvr.Validate)
}

// UnmarshalJSON parses a canonical group/version/resource string from a JSON scalar.
//
// Non-string JSON, including null, is rejected before parsing.
func (gvr *GroupVersionResource) UnmarshalJSON(data []byte) error {
	if gvr == nil {
		return nilUnmarshalReceiver("group/version/resource")
	}
	val, err := unmarshalJSONString("group/version/resource", data)
	if err != nil {
		return err
	}
	return gvr.UnmarshalText([]byte(val))
}

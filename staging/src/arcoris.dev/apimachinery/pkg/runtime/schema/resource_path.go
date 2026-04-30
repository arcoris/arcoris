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

// ResourcePath identifies a resource collection and optional subresource.
//
// The canonical form is resource or resource/subresource. The slash form is
// used only when a non-empty subresource is present. ResourcePath deliberately
// omits group and version so routing layers can compose it with a
// GroupVersionResource without re-parsing strings.
type ResourcePath struct {
	Resource    Resource
	Subresource Subresource
}

// ParseResourcePath parses resource or resource/subresource in canonical form.
//
// The parser accepts exactly one optional slash. Empty resources, empty
// subresources after a slash, and deeper paths are rejected.
func ParseResourcePath(value string) (ResourcePath, error) {
	if value == "" {
		return ResourcePath{}, invalid("resource path", value, "resource is required")
	}
	parts := strings.Split(value, "/")
	if len(parts) > 2 {
		return ResourcePath{}, invalid("resource path", value, "expected canonical form resource or resource/subresource")
	}
	resource, err := ParseResource(parts[0])
	if err != nil {
		return ResourcePath{}, invalidValue("resource path", value, err)
	}
	if len(parts) == 1 {
		return ResourcePath{Resource: resource}, nil
	}
	if parts[1] == "" {
		return ResourcePath{}, invalid("resource path", value, "subresource is required after '/'")
	}
	subresource, err := ParseSubresource(parts[1])
	if err != nil {
		return ResourcePath{}, invalidValue("resource path", value, err)
	}
	return ResourcePath{Resource: resource, Subresource: subresource}, nil
}

// String returns resource or resource/subresource in canonical form without revalidating it.
//
// The subresource slash is emitted only when Subresource is non-empty.
func (rp ResourcePath) String() string {
	if rp.Subresource.IsZero() {
		return rp.Resource.String()
	}
	return rp.Resource.String() + "/" + rp.Subresource.String()
}

// Identifier returns the canonical resource path string.
//
// It is equivalent to String and exists for API symmetry with versioned
// resource identities.
func (rp ResourcePath) Identifier() string {
	return rp.String()
}

// Validate checks that the resource path is a complete canonical identity.
//
// Resource is required. Subresource may be empty, but if present it must satisfy
// the same single-label rule as Resource.
func (rp ResourcePath) Validate() error {
	if rp.Resource.IsZero() {
		return invalid("resource path", rp.String(), "resource is required")
	}
	if err := rp.Resource.Validate(); err != nil {
		return invalidValue("resource path", rp.String(), err)
	}
	if err := rp.Subresource.Validate(); err != nil {
		return invalidValue("resource path", rp.String(), err)
	}
	return nil
}

// IsZero reports whether both resource and subresource are empty.
//
// Zero is useful as an optional sentinel; it is rejected by Validate because a
// complete path requires a resource.
func (rp ResourcePath) IsZero() bool {
	return rp.Resource.IsZero() && rp.Subresource.IsZero()
}

// HasSubresource reports whether the path includes a subresource.
//
// It is a semantic check on the structured field, not a string search.
func (rp ResourcePath) HasSubresource() bool {
	return !rp.Subresource.IsZero()
}

// MarshalText returns the canonical resource path after validation.
//
// Invalid direct literals, including zero paths, are rejected.
func (rp ResourcePath) MarshalText() ([]byte, error) {
	if err := rp.Validate(); err != nil {
		return nil, err
	}
	return []byte(rp.String()), nil
}

// UnmarshalText parses a canonical resource path.
//
// The method delegates to ParseResourcePath and rejects path-like alternatives
// with more than one slash.
func (rp *ResourcePath) UnmarshalText(data []byte) error {
	if rp == nil {
		return nilUnmarshalReceiver("resource path")
	}
	parsed, err := ParseResourcePath(string(data))
	if err != nil {
		return err
	}
	*rp = parsed
	return nil
}

// MarshalJSON returns the canonical resource path as a JSON scalar.
//
// ResourcePath does not serialize as an object with resource/subresource
// fields.
func (rp ResourcePath) MarshalJSON() ([]byte, error) {
	return marshalJSONString("resource path", rp.String(), rp.Validate)
}

// UnmarshalJSON parses a canonical resource path from a JSON scalar.
//
// Non-string JSON, including null, is rejected.
func (rp *ResourcePath) UnmarshalJSON(data []byte) error {
	if rp == nil {
		return nilUnmarshalReceiver("resource path")
	}
	value, err := unmarshalJSONString("resource path", data)
	if err != nil {
		return err
	}
	return rp.UnmarshalText([]byte(value))
}

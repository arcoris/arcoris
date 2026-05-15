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

// GroupResource identifies an API resource collection without binding it to a version.
//
// The canonical form is "resource" for the core group and "resource.group" for
// named groups. Resource names do not include subresources; use ResourcePath or
// GroupVersionResourcePath when a subresource segment is part of the identity.
// This is a resource identity, not a type identity.
type GroupResource struct {
	Group    Group
	Resource Resource
}

// ParseGroupResource parses a canonical group/resource identity.
//
// Core resources use the bare resource form. Named groups use resource.group.
// The parser rejects dotted resource names because the first dot is always the
// boundary between resource and group.
func ParseGroupResource(val string) (GroupResource, error) {
	if val == "" {
		return GroupResource{}, invalid("group/resource", val, "resource is required")
	}
	resourcePart, groupPart, hasGroup := strings.Cut(val, ".")
	if !hasGroup {
		resource, err := ParseResource(val)
		if err != nil {
			return GroupResource{}, invalidValue("group/resource", val, err)
		}
		return GroupResource{Resource: resource}, nil
	}
	if resourcePart == "" || groupPart == "" {
		return GroupResource{}, invalid("group/resource", val, "expected canonical form resource.group")
	}
	resource, err := ParseResource(resourcePart)
	if err != nil {
		return GroupResource{}, invalidValue("group/resource", val, err)
	}
	group, err := ParseGroup(groupPart)
	if err != nil {
		return GroupResource{}, invalidValue("group/resource", val, err)
	}
	return GroupResource{Group: group, Resource: resource}, nil
}

// String returns the canonical group/resource string without revalidating it.
//
// The core group is emitted as the bare resource; named groups append "." and
// the group name.
func (gr GroupResource) String() string {
	if gr.Group.IsZero() {
		return gr.Resource.String()
	}
	return gr.Resource.String() + "." + gr.Group.String()
}

// Identifier returns the canonical group/resource identity string.
//
// It is equivalent to String and exists for API symmetry with other schema
// identities.
func (gr GroupResource) Identifier() string {
	return gr.String()
}

// Validate checks that the group/resource is a complete canonical identity.
//
// Group may be empty for the core group, but Resource is required.
func (gr GroupResource) Validate() error {
	if gr.Resource.IsZero() {
		return invalid("group/resource", gr.String(), "resource is required")
	}
	if err := gr.Group.Validate(); err != nil {
		return invalidValue("group/resource", gr.String(), err)
	}
	if err := gr.Resource.Validate(); err != nil {
		return invalidValue("group/resource", gr.String(), err)
	}
	return nil
}

// IsZero reports whether both group and resource are empty.
//
// Zero is an optional sentinel only; it is not a valid complete resource
// identity.
func (gr GroupResource) IsZero() bool {
	return gr.Group.IsZero() && gr.Resource.IsZero()
}

// MarshalText returns the canonical group/resource string after validation.
//
// Invalid direct literals are rejected before serialization.
func (gr GroupResource) MarshalText() ([]byte, error) {
	if err := gr.Validate(); err != nil {
		return nil, err
	}
	return []byte(gr.String()), nil
}

// UnmarshalText parses a canonical group/resource string.
//
// The method delegates to ParseGroupResource and therefore rejects alternative
// dotted triplet forms.
func (gr *GroupResource) UnmarshalText(data []byte) error {
	if gr == nil {
		return nilUnmarshalReceiver("group/resource")
	}
	parsed, err := ParseGroupResource(string(data))
	if err != nil {
		return err
	}
	*gr = parsed
	return nil
}

// MarshalJSON returns the canonical group/resource string as a JSON scalar.
//
// Object-form JSON is not part of the schema contract.
func (gr GroupResource) MarshalJSON() ([]byte, error) {
	return marshalJSONString("group/resource", gr.String(), gr.Validate)
}

// UnmarshalJSON parses a canonical group/resource string from a JSON scalar.
//
// Non-string JSON, including null, is rejected.
func (gr *GroupResource) UnmarshalJSON(data []byte) error {
	if gr == nil {
		return nilUnmarshalReceiver("group/resource")
	}
	val, err := unmarshalJSONString("group/resource", data)
	if err != nil {
		return err
	}
	return gr.UnmarshalText([]byte(val))
}

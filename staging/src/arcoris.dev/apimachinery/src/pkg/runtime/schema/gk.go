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

// GroupKind identifies an API object kind without binding it to a version.
//
// The canonical form is "Kind" for the core group and "Kind.group" for named
// groups. This identity is useful before version negotiation has selected a
// concrete GroupVersionKind. It is a type identity and never names a resource
// collection.
type GroupKind struct {
	Group Group
	Kind  Kind
}

// ParseGroupKind parses a canonical group/kind identity.
//
// Core kinds use the bare Kind form. Named groups use Kind.group, where the
// group suffix is a DNS-1123 subdomain. The parser rejects malformed kinds and
// groups without attempting alternate interpretations.
func ParseGroupKind(val string) (GroupKind, error) {
	if val == "" {
		return GroupKind{}, invalid("group/kind", val, "kind is required")
	}
	kindPart, groupPart, hasGroup := strings.Cut(val, ".")
	if !hasGroup {
		kind, err := ParseKind(val)
		if err != nil {
			return GroupKind{}, invalidValue("group/kind", val, err)
		}
		return GroupKind{Kind: kind}, nil
	}
	if kindPart == "" || groupPart == "" {
		return GroupKind{}, invalid("group/kind", val, "expected canonical form kind.group")
	}
	kind, err := ParseKind(kindPart)
	if err != nil {
		return GroupKind{}, invalidValue("group/kind", val, err)
	}
	group, err := ParseGroup(groupPart)
	if err != nil {
		return GroupKind{}, invalidValue("group/kind", val, err)
	}
	return GroupKind{Group: group, Kind: kind}, nil
}

// String returns the canonical group/kind string without revalidating it.
//
// The core group is emitted as the bare kind; named groups append "." and the
// group name.
func (gk GroupKind) String() string {
	if gk.Group.IsZero() {
		return gk.Kind.String()
	}
	return gk.Kind.String() + "." + gk.Group.String()
}

// Identifier returns the canonical group/kind identity string.
//
// It is equivalent to String and exists for API symmetry with other schema
// identities.
func (gk GroupKind) Identifier() string {
	return gk.String()
}

// Validate checks that the group/kind is a complete canonical identity.
//
// Group may be empty for the core group, but Kind is required.
func (gk GroupKind) Validate() error {
	if gk.Kind.IsZero() {
		return invalid("group/kind", gk.String(), "kind is required")
	}
	if err := gk.Group.Validate(); err != nil {
		return invalidValue("group/kind", gk.String(), err)
	}
	if err := gk.Kind.Validate(); err != nil {
		return invalidValue("group/kind", gk.String(), err)
	}
	return nil
}

// IsZero reports whether both group and kind are empty.
//
// Zero is useful only as an optional sentinel and is not a complete type
// identity.
func (gk GroupKind) IsZero() bool {
	return gk.Group.IsZero() && gk.Kind.IsZero()
}

// MarshalText returns the canonical group/kind string after validation.
//
// Invalid direct literals are rejected before serialization.
func (gk GroupKind) MarshalText() ([]byte, error) {
	if err := gk.Validate(); err != nil {
		return nil, err
	}
	return []byte(gk.String()), nil
}

// UnmarshalText parses a canonical group/kind string.
//
// The method uses ParseGroupKind exactly and performs no normalization.
func (gk *GroupKind) UnmarshalText(data []byte) error {
	if gk == nil {
		return nilUnmarshalReceiver("group/kind")
	}
	parsed, err := ParseGroupKind(string(data))
	if err != nil {
		return err
	}
	*gk = parsed
	return nil
}

// MarshalJSON returns the canonical group/kind string as a JSON scalar.
//
// Object-form JSON is not part of the schema contract.
func (gk GroupKind) MarshalJSON() ([]byte, error) {
	return marshalJSONString("group/kind", gk.String(), gk.Validate)
}

// UnmarshalJSON parses a canonical group/kind string from a JSON scalar.
//
// Non-string JSON, including null, is rejected.
func (gk *GroupKind) UnmarshalJSON(data []byte) error {
	if gk == nil {
		return nilUnmarshalReceiver("group/kind")
	}
	val, err := unmarshalJSONString("group/kind", data)
	if err != nil {
		return err
	}
	return gk.UnmarshalText([]byte(val))
}

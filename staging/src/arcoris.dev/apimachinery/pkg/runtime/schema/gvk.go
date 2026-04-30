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

// GroupVersionKind identifies a concrete versioned API object schema.
//
// The canonical form is "group/version, Kind=Kind" or "version, Kind=Kind" for
// the core group. This is a type identity, not a resource endpoint identity;
// serializers and runtime object metadata can use it without depending on
// resource routing.
type GroupVersionKind struct {
	Group   Group
	Version Version
	Kind    Kind
}

// ParseGroupVersionKind parses a canonical group/version/kind identity.
//
// The parser intentionally rejects dotted legacy forms such as
// "Kind.version.group"; ARCORIS uses a single unambiguous string grammar.
func ParseGroupVersionKind(value string) (GroupVersionKind, error) {
	if value == "" {
		return GroupVersionKind{}, invalid("group/version/kind", value, "group/version and kind are required")
	}
	gvPart, kindPart, ok := strings.Cut(value, ", Kind=")
	if !ok || gvPart == "" || kindPart == "" {
		return GroupVersionKind{}, invalid("group/version/kind", value, "expected canonical form group/version, Kind=Kind")
	}
	gv, err := ParseGroupVersion(gvPart)
	if err != nil {
		return GroupVersionKind{}, invalidValue("group/version/kind", value, err)
	}
	kind, err := ParseKind(kindPart)
	if err != nil {
		return GroupVersionKind{}, invalidValue("group/version/kind", value, err)
	}
	return GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: kind}, nil
}

// ParseAPIVersionAndKind parses separate apiVersion and kind field values.
//
// Both fields are validated strictly. This function exists for object payloads
// that carry apiVersion and kind as separate JSON fields while schema itself
// still stores one structured identity.
func ParseAPIVersionAndKind(apiVersion string, kind string) (GroupVersionKind, error) {
	gv, err := ParseAPIVersion(apiVersion)
	if err != nil {
		return GroupVersionKind{}, err
	}
	parsedKind, err := ParseKind(kind)
	if err != nil {
		return GroupVersionKind{}, err
	}
	return gv.WithKind(parsedKind), nil
}

// String returns the canonical group/version/kind string without revalidating it.
//
// The textual form is intentionally explicit about Kind so it cannot be
// confused with resource or path-like identities.
func (gvk GroupVersionKind) String() string {
	return gvk.GroupVersion().String() + ", Kind=" + gvk.Kind.String()
}

// Identifier returns the canonical group/version/kind identity string.
//
// It is equivalent to String and suitable for logs, maps, and diagnostics.
func (gvk GroupVersionKind) Identifier() string {
	return gvk.String()
}

// Validate checks that the group/version/kind is a complete canonical identity.
//
// Group may be empty for the core group. Version and Kind are required, so zero
// and partial GVK values are rejected.
func (gvk GroupVersionKind) Validate() error {
	if gvk.Version.IsZero() {
		return invalid("group/version/kind", gvk.String(), "version is required")
	}
	if gvk.Kind.IsZero() {
		return invalid("group/version/kind", gvk.String(), "kind is required")
	}
	if err := gvk.GroupVersion().Validate(); err != nil {
		return invalidValue("group/version/kind", gvk.String(), err)
	}
	if err := gvk.Kind.Validate(); err != nil {
		return invalidValue("group/version/kind", gvk.String(), err)
	}
	return nil
}

// IsZero reports whether group, version, and kind are all empty.
//
// It is intended only for optional-field checks; it does not mean the value is
// valid for serialization or discovery.
func (gvk GroupVersionKind) IsZero() bool {
	return gvk.Group.IsZero() && gvk.Version.IsZero() && gvk.Kind.IsZero()
}

// GroupVersion returns the group/version portion of the identity.
//
// The returned value preserves the exact group and version fields without
// validating or normalizing them.
func (gvk GroupVersionKind) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvk.Group, Version: gvk.Version}
}

// GroupKind returns the group/kind portion of the identity.
//
// This is useful before or after version negotiation when callers need to
// compare object types independent of version.
func (gvk GroupVersionKind) GroupKind() GroupKind {
	return GroupKind{Group: gvk.Group, Kind: gvk.Kind}
}

// ToAPIVersionAndKind returns values suitable for apiVersion and kind fields.
//
// The returned strings use the same canonical spelling as String and Kind.
func (gvk GroupVersionKind) ToAPIVersionAndKind() (apiVersion string, kind string) {
	return gvk.GroupVersion().APIVersion(), gvk.Kind.String()
}

// MarshalText returns the canonical group/version/kind string after validation.
//
// Invalid direct literals and partial identities are rejected before
// serialization.
func (gvk GroupVersionKind) MarshalText() ([]byte, error) {
	if err := gvk.Validate(); err != nil {
		return nil, err
	}
	return []byte(gvk.String()), nil
}

// UnmarshalText parses a canonical group/version/kind string.
//
// The method rejects legacy dotted triplets and any string that is missing the
// ", Kind=" separator.
func (gvk *GroupVersionKind) UnmarshalText(data []byte) error {
	if gvk == nil {
		return nilUnmarshalReceiver("group/version/kind")
	}
	parsed, err := ParseGroupVersionKind(string(data))
	if err != nil {
		return err
	}
	*gvk = parsed
	return nil
}

// MarshalJSON returns the canonical group/version/kind string as a JSON scalar.
//
// It never emits an object with group, version, and kind fields.
func (gvk GroupVersionKind) MarshalJSON() ([]byte, error) {
	return marshalJSONString("group/version/kind", gvk.String(), gvk.Validate)
}

// UnmarshalJSON parses a canonical group/version/kind string from a JSON scalar.
//
// Non-string JSON, including null, is rejected before parsing.
func (gvk *GroupVersionKind) UnmarshalJSON(data []byte) error {
	if gvk == nil {
		return nilUnmarshalReceiver("group/version/kind")
	}
	value, err := unmarshalJSONString("group/version/kind", data)
	if err != nil {
		return err
	}
	return gvk.UnmarshalText([]byte(value))
}

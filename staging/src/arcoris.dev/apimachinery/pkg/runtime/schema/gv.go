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

// GroupVersion identifies a versioned API namespace.
//
// The canonical form is "version" for the core group and "group/version" for
// named groups. A zero GroupVersion is useful as an optional field sentinel but
// is not a valid complete identity. Runtime, metadata, and discovery layers use
// this value as the namespace half of type and resource identities.
type GroupVersion struct {
	Group   Group
	Version Version
}

// ParseGroupVersion parses a canonical group/version identity.
//
// The parser accepts exactly "v1" for a core group version or
// "control.arcoris.dev/v1alpha1" for a named group version. It does not trim
// whitespace and rejects slash-heavy API paths such as "apps/v1/resources".
func ParseGroupVersion(val string) (GroupVersion, error) {
	if val == "" {
		return GroupVersion{}, invalid("group/version", val, "version is required")
	}
	if strings.Count(val, "/") > 1 {
		return GroupVersion{}, invalid("group/version", val, "expected canonical form version or group/version")
	}
	if !strings.Contains(val, "/") {
		version, err := ParseVersion(val)
		if err != nil {
			return GroupVersion{}, invalidValue("group/version", val, err)
		}
		return GroupVersion{Version: version}, nil
	}

	groupPart, versionPart, ok := strings.Cut(val, "/")
	if !ok || groupPart == "" || versionPart == "" {
		return GroupVersion{}, invalid("group/version", val, "expected canonical form group/version")
	}
	group, err := ParseGroup(groupPart)
	if err != nil {
		return GroupVersion{}, invalidValue("group/version", val, err)
	}
	version, err := ParseVersion(versionPart)
	if err != nil {
		return GroupVersion{}, invalidValue("group/version", val, err)
	}
	return GroupVersion{Group: group, Version: version}, nil
}

// ParseAPIVersion parses an apiVersion field value as a canonical GroupVersion.
//
// The function exists to make object-field parsing explicit while preserving
// exactly the same strict grammar as ParseGroupVersion.
func ParseAPIVersion(val string) (GroupVersion, error) {
	return ParseGroupVersion(val)
}

// String returns the canonical group/version string without revalidating it.
//
// It emits the core group as just the version and named groups as group/version.
func (gv GroupVersion) String() string {
	if gv.Group.IsZero() {
		return gv.Version.String()
	}
	return gv.Group.String() + "/" + gv.Version.String()
}

// APIVersion returns the canonical apiVersion field value.
//
// This is an alias for String with object-field intent made explicit.
func (gv GroupVersion) APIVersion() string {
	return gv.String()
}

// Identifier returns the canonical group/version identity string.
//
// Foundation and higher-level code can use Identifier when it needs a stable
// human-readable key without implying any Kubernetes-style compatibility.
func (gv GroupVersion) Identifier() string {
	return gv.String()
}

// Validate checks that the group/version is a complete canonical identity.
//
// Group may be empty for the core group, but Version must be present and valid.
// Zero GroupVersion is therefore invalid as a complete identity.
func (gv GroupVersion) Validate() error {
	if gv.Version.IsZero() {
		return invalid("group/version", gv.String(), "version is required")
	}
	if err := gv.Group.Validate(); err != nil {
		return invalidValue("group/version", gv.String(), err)
	}
	if err := gv.Version.Validate(); err != nil {
		return invalidValue("group/version", gv.String(), err)
	}
	return nil
}

// IsZero reports whether both group and version are empty.
//
// IsZero is intended for optional-field checks only; it is not a substitute for
// Validate when an identity is required.
func (gv GroupVersion) IsZero() bool {
	return gv.Group.IsZero() && gv.Version.IsZero()
}

// WithKind returns the GroupVersionKind formed from this group/version and kind.
//
// The method composes identity fields without validation so callers can build a
// value first and validate it at the boundary.
func (gv GroupVersion) WithKind(kind Kind) GroupVersionKind {
	return GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: kind}
}

// WithResource returns the GroupVersionResource formed from this group/version and resource.
//
// The method does not infer pluralization or resource naming; the provided
// Resource is used exactly.
func (gv GroupVersion) WithResource(resource Resource) GroupVersionResource {
	return GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource}
}

// MarshalText returns the canonical group/version string after validation.
//
// Invalid direct literals, including zero values and missing versions, are
// rejected.
func (gv GroupVersion) MarshalText() ([]byte, error) {
	if err := gv.Validate(); err != nil {
		return nil, err
	}
	return []byte(gv.String()), nil
}

// UnmarshalText parses a canonical group/version string.
//
// The method delegates to ParseGroupVersion and therefore rejects path-like
// alternatives such as "apps/v1/deployments".
func (gv *GroupVersion) UnmarshalText(data []byte) error {
	if gv == nil {
		return nilUnmarshalReceiver("group/version")
	}
	parsed, err := ParseGroupVersion(string(data))
	if err != nil {
		return err
	}
	*gv = parsed
	return nil
}

// MarshalJSON returns the canonical group/version string as a JSON scalar.
//
// Composite identifiers deliberately do not expose struct-field JSON shapes.
func (gv GroupVersion) MarshalJSON() ([]byte, error) {
	return marshalJSONString("group/version", gv.String(), gv.Validate)
}

// UnmarshalJSON parses a canonical group/version string from a JSON scalar.
//
// Non-string JSON, including null, is rejected before parsing.
func (gv *GroupVersion) UnmarshalJSON(data []byte) error {
	if gv == nil {
		return nilUnmarshalReceiver("group/version")
	}
	val, err := unmarshalJSONString("group/version", data)
	if err != nil {
		return err
	}
	return gv.UnmarshalText([]byte(val))
}

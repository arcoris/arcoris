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

// Resource identifies an API resource collection.
//
// A resource is a DNS-1123 single label. It never includes a group, version, or
// subresource segment. Use GroupResource, GroupVersionResource, ResourcePath,
// or GroupVersionResourcePath when additional identity dimensions are needed.
type Resource string

// ParseResource parses a canonical API resource name.
//
// The parser accepts only a non-empty DNS-1123 single label. Dots, slashes,
// underscores, uppercase letters, and surrounding whitespace are rejected.
func ParseResource(val string) (Resource, error) {
	resource := Resource(val)
	if err := resource.Validate(); err != nil {
		return "", err
	}
	return resource, nil
}

// String returns the canonical resource string without revalidating it.
//
// Direct Resource literals should be validated before they are used as schema
// identity.
func (r Resource) String() string {
	return string(r)
}

// Validate checks that the resource is a non-empty DNS-1123 single label.
//
// The single-label rule prevents resources from accidentally encoding groups,
// versions, or subresources.
func (r Resource) Validate() error {
	return validateDNS1123SingleLabel("resource", string(r), false)
}

// IsZero reports whether the resource is empty.
//
// Zero is allowed only as an incomplete construction state; complete resource
// identities require a non-empty Resource.
func (r Resource) IsZero() bool {
	return r == ""
}

// MarshalText returns the canonical resource string after validation.
//
// Invalid direct literals are rejected instead of being serialized.
func (r Resource) MarshalText() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return []byte(r.String()), nil
}

// UnmarshalText parses a canonical resource string.
//
// The method delegates to ParseResource and inherits its strict no-trim
// behavior.
func (r *Resource) UnmarshalText(data []byte) error {
	if r == nil {
		return nilUnmarshalReceiver("resource")
	}
	parsed, err := ParseResource(string(data))
	if err != nil {
		return err
	}
	*r = parsed
	return nil
}

// MarshalJSON returns the canonical resource string as a JSON scalar.
//
// Schema identifiers serialize as scalar strings, not objects.
func (r Resource) MarshalJSON() ([]byte, error) {
	return marshalJSONString("resource", r.String(), r.Validate)
}

// UnmarshalJSON parses a canonical resource string from a JSON scalar.
//
// Non-string JSON, including null, is rejected.
func (r *Resource) UnmarshalJSON(data []byte) error {
	if r == nil {
		return nilUnmarshalReceiver("resource")
	}
	val, err := unmarshalJSONString("resource", data)
	if err != nil {
		return err
	}
	return r.UnmarshalText([]byte(val))
}

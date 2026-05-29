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

// Group identifies an API group.
//
// The empty group is the canonical spelling of the core API group. Non-empty
// groups must be qualified lowercase ASCII DNS subdomains with at least one
// dot. The package never trims, lowercases, aliases, or normalizes group text,
// because a Group is a canonical identity token rather than user-facing input.
type Group string

const (
	// CoreGroup is the canonical empty API group.
	CoreGroup Group = ""
)

// String returns the canonical group text without revalidating it.
//
// Direct Group literals can contain invalid text. Call Validate at trust
// boundaries or use ParseGroup when accepting external input.
func (g Group) String() string {
	return string(g)
}

// IsZero reports whether the group is the core group.
//
// Group is the only atomic identity where the zero value is also a valid
// complete value. Composite identities still require their version, kind, or
// resource fields.
func (g Group) IsZero() bool {
	return g == CoreGroup
}

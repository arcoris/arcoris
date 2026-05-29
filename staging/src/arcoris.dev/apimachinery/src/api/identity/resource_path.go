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

// ResourcePath identifies a resource collection and optional subresource.
//
// The canonical form is "resource" or "resource/subresource". A slash requires
// a non-empty subresource and only one subresource segment is allowed.
type ResourcePath struct {
	Resource    Resource
	Subresource Subresource
}

// String returns the canonical resource path without revalidating it.
func (rp ResourcePath) String() string {
	return joinResourcePath(rp.Resource, rp.Subresource)
}

// Identifier returns the canonical resource path identity string.
func (rp ResourcePath) Identifier() string {
	return rp.String()
}

// IsZero reports whether resource and subresource are both absent.
func (rp ResourcePath) IsZero() bool {
	return rp.Resource.IsZero() && rp.Subresource.IsZero()
}

// HasSubresource reports whether the path includes a subresource segment.
func (rp ResourcePath) HasSubresource() bool {
	return !rp.Subresource.IsZero()
}

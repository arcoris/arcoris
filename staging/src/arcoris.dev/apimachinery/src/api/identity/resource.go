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

// Resource identifies an API resource collection.
//
// Resources are DNS-1123 single labels: lowercase ASCII letters, digits, and
// hyphen, starting and ending with a lowercase letter or digit. Resource does
// not enforce pluralization and never contains a subresource segment.
type Resource string

// String returns the canonical resource text without revalidating it.
func (r Resource) String() string {
	return string(r)
}

// IsZero reports whether the resource is absent.
func (r Resource) IsZero() bool {
	return r == ""
}

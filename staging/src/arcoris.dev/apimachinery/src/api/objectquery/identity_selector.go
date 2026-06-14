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

package objectquery

// IdentitySelector filters by authoritative objectstore storage identity.
//
// It intentionally matches objectstore.ListItem.Key.Object rather than
// arbitrary object payload metadata fields.
type IdentitySelector struct {
	// Namespace requires an exact object namespace when non-zero.
	Namespace NamespaceRequirement

	// Name requires an exact object name when non-zero.
	Name NameRequirement
}

// IsZero reports whether s carries no identity requirements.
func (s IdentitySelector) IsZero() bool {
	return s.Namespace.IsZero() && s.Name.IsZero()
}

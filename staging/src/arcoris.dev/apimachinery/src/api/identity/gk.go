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

// GroupKind identifies an API kind within a group.
//
// The core group uses only the kind text, for example "Pod". A named group
// uses "group#Kind", for example "control.arcoris.dev#Worker".
type GroupKind struct {
	Group Group
	Kind  Kind
}

// String returns the canonical group/kind text without revalidating it.
func (gk GroupKind) String() string {
	return joinGroupKind(gk.Group, gk.Kind)
}

// Identifier returns the canonical group/kind identity string.
func (gk GroupKind) Identifier() string {
	return gk.String()
}

// IsZero reports whether group and kind are both absent.
func (gk GroupKind) IsZero() bool {
	return gk.Group.IsZero() && gk.Kind.IsZero()
}

// WithVersion composes this group/kind with a version.
func (gk GroupKind) WithVersion(version Version) GroupVersionKind {
	return GroupVersionKind{Group: gk.Group, Version: version, Kind: gk.Kind}
}

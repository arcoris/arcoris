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

// GroupVersionKind identifies a concrete versioned API kind.
//
// The canonical form is GroupVersion + "#" + Kind, for example "v1#Pod" or
// "control.arcoris.dev/v1#Worker". This package does not expose object-field
// apiVersion/kind helpers; those belong to metadata or codec layers.
type GroupVersionKind struct {
	Group   Group
	Version Version
	Kind    Kind
}

// String returns the canonical group/version/kind text without revalidating it.
func (gvk GroupVersionKind) String() string {
	return joinGroupVersionKind(gvk.GroupVersion(), gvk.Kind)
}

// CanonicalText validates the group/version/kind identity and returns its canonical text.
func (gvk GroupVersionKind) CanonicalText() (string, error) {
	if err := gvk.Validate(); err != nil {
		return "", err
	}

	return gvk.String(), nil
}

// IsZero reports whether group, version, and kind are all absent.
func (gvk GroupVersionKind) IsZero() bool {
	return gvk.Group.IsZero() &&
		gvk.Version.IsZero() &&
		gvk.Kind.IsZero()
}

// GroupVersion returns the group/version portion of the identity.
func (gvk GroupVersionKind) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvk.Group, Version: gvk.Version}
}

// GroupKind returns the group/kind portion of the identity.
func (gvk GroupVersionKind) GroupKind() GroupKind {
	return GroupKind{Group: gvk.Group, Kind: gvk.Kind}
}

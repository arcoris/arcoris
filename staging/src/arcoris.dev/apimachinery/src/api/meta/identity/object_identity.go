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

// ObjectIdentity identifies a concrete object incarnation by name and UID.
type ObjectIdentity struct {
	// Namespace is optional and means absent when empty.
	Namespace Namespace `json:"namespace,omitempty"`
	// Name is the required object metadata name.
	Name Name `json:"name"`
	// UID identifies one concrete object incarnation.
	UID UID `json:"uid"`
}

// IsZero reports whether all identity fields are absent.
func (i ObjectIdentity) IsZero() bool {
	return i.Namespace.IsZero() && i.Name.IsZero() && i.UID.IsZero()
}

// ObjectName returns the namespace/name portion of the identity.
func (i ObjectIdentity) ObjectName() ObjectName {
	return ObjectName{Namespace: i.Namespace, Name: i.Name}
}

// String returns diagnostic text for the object identity.
//
// The result is intentionally diagnostic only. Storage keys, route keys, and
// cache keys need their own explicit formats in higher layers.
func (i ObjectIdentity) String() string {
	name := i.ObjectName().String()
	if i.UID.IsZero() {
		return name
	}
	return name + "#" + i.UID.String()
}

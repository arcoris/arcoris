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

// ObjectName identifies an object by namespace and metadata name.
//
// Its String form is canonical diagnostic text only. It is not a storage key,
// route path, cache key, or resource identity.
type ObjectName struct {
	// Namespace is optional and means absent when empty.
	Namespace Namespace `json:"namespace,omitempty"`
	// Name is the required object metadata name.
	Name Name `json:"name"`
}

// IsZero reports whether namespace and name are both absent.
func (n ObjectName) IsZero() bool {
	return n.Namespace.IsZero() && n.Name.IsZero()
}

// String returns "name" or "namespace/name" diagnostic text without validating it.
//
// The result is for logs and diagnostics. It is not a storage key, route path,
// cache key, or authorization resource string.
func (n ObjectName) String() string {
	if n.Namespace.IsZero() {
		return n.Name.String()
	}
	return n.Namespace.String() + "/" + n.Name.String()
}

// CanonicalText validates the object name and returns its canonical text.
func (n ObjectName) CanonicalText() (string, error) {
	if err := n.Validate(); err != nil {
		return "", err
	}

	return n.String(), nil
}

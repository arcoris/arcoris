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

import (
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
)

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

// InNamespace constructs an identity selector for one namespace.
func InNamespace(namespace metaidentity.Namespace) (IdentitySelector, error) {
	req, err := NamespaceEquals(namespace)
	if err != nil {
		return IdentitySelector{}, err
	}

	return IdentitySelector{Namespace: req}, nil
}

// WithName constructs an identity selector for one object name.
func WithName(name metaidentity.Name) (IdentitySelector, error) {
	req, err := NameEquals(name)
	if err != nil {
		return IdentitySelector{}, err
	}

	return IdentitySelector{Name: req}, nil
}

// WithObject constructs an identity selector for one namespace/name pair.
func WithObject(namespace metaidentity.Namespace, name metaidentity.Name) (IdentitySelector, error) {
	namespaceReq, err := NamespaceEquals(namespace)
	if err != nil {
		return IdentitySelector{}, err
	}
	nameReq, err := NameEquals(name)
	if err != nil {
		return IdentitySelector{}, err
	}

	return IdentitySelector{Namespace: namespaceReq, Name: nameReq}, nil
}

// IsZero reports whether s carries no identity requirements.
func (s IdentitySelector) IsZero() bool {
	return s.Namespace.IsZero() && s.Name.IsZero()
}

// Validate checks identity selector requirements.
func (s IdentitySelector) Validate() error {
	if err := s.Namespace.validate("query.identity.namespace"); err != nil {
		return err
	}
	if err := s.Name.validate("query.identity.name"); err != nil {
		return err
	}

	return nil
}

// match checks item against the identity selector.
func (s IdentitySelector) match(item objectstore.ListItem) bool {
	objectName := item.Key.Object
	if !s.Namespace.IsZero() && objectName.Namespace != s.Namespace.namespace {
		return false
	}
	if !s.Name.IsZero() && objectName.Name != s.Name.name {
		return false
	}

	return true
}

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

package objectstore

import (
	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

// NewKey validates and returns an object store key.
func NewKey(resource apiidentity.GroupVersionResource, object metaidentity.ObjectName) (Key, error) {
	key := Key{Resource: resource, Object: object}
	if err := ValidateKey(key); err != nil {
		return Key{}, err
	}

	return key, nil
}

// MustKey validates and returns an object store key.
//
// MustKey panics with a structured objectstore error when resource or object
// name is invalid. It is intended for package-level fixtures and examples.
func MustKey(resource apiidentity.GroupVersionResource, object metaidentity.ObjectName) Key {
	key, err := NewKey(resource, object)
	if err != nil {
		panic(err)
	}

	return key
}

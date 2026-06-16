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
	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
)

// ResourceEquals matches items whose authoritative storage key belongs to gvr.
func ResourceEquals(gvr apiidentity.GroupVersionResource) (Query, error) {
	if err := gvr.Validate(); err != nil {
		return Query{}, invalidTermError("invalid resource: %w", err)
	}

	return termQuery(term{kind: termResource, resource: gvr, operator: OperatorEquals}), nil
}

// ObjectInNamespace matches items by storage key namespace, not payload
// metadata.
func ObjectInNamespace(namespace metaidentity.Namespace) (Query, error) {
	if err := namespace.ValidateLexical(); err != nil {
		return Query{}, invalidTermError("invalid namespace: %w", err)
	}

	return termQuery(term{kind: termNamespace, namespace: namespace, operator: OperatorEquals}), nil
}

// ObjectWithName matches items by storage key object name, not payload
// metadata.
func ObjectWithName(name metaidentity.Name) (Query, error) {
	if err := name.ValidateLexical(); err != nil {
		return Query{}, invalidTermError("invalid name: %w", err)
	}

	return termQuery(term{kind: termName, name: name, operator: OperatorEquals}), nil
}

// ObjectEquals matches items by storage key namespace/name identity.
func ObjectEquals(namespace metaidentity.Namespace, name metaidentity.Name) (Query, error) {
	objectName := metaidentity.ObjectName{Namespace: namespace, Name: name}
	if err := objectName.ValidateLexical(); err != nil {
		return Query{}, invalidTermError("invalid object identity: %w", err)
	}

	return termQuery(term{
		kind:      termObject,
		namespace: namespace,
		name:      name,
		operator:  OperatorEquals,
	}), nil
}

// KeyEquals matches one exact object store key, including resource and object
// identity.
func KeyEquals(key objectstore.Key) (Query, error) {
	if err := objectstore.ValidateKey(key); err != nil {
		return Query{}, invalidTermError("invalid key: %w", err)
	}

	return termQuery(term{kind: termKey, key: key, operator: OperatorEquals}), nil
}

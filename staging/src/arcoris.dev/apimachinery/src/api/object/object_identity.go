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

package object

import (
	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

// GroupVersionKind returns the API group/version/kind stored in TypeMeta.
func (o Object[D, O]) GroupVersionKind() apiidentity.GroupVersionKind {
	return o.TypeMeta.GroupVersionKind()
}

// ObjectName returns the namespace/name stored in ObjectMeta.
func (o Object[D, O]) ObjectName() metaidentity.ObjectName {
	return o.ObjectMeta.ObjectName()
}

// ObjectIdentity returns the namespace/name/UID stored in ObjectMeta.
func (o Object[D, O]) ObjectIdentity() metaidentity.ObjectIdentity {
	return o.ObjectMeta.ObjectIdentity()
}

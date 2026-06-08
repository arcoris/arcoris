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

// Key identifies one committed live object inside a store.
//
// Resource names the versioned resource collection. Object names the concrete
// namespace/name identity. The store validates only structural key shape; it
// does not decide whether a resource is namespace-scoped or cluster-scoped.
type Key struct {
	// Resource identifies the API resource collection that owns the object.
	Resource apiidentity.GroupVersionResource
	// Object identifies the object by namespace and metadata name.
	Object metaidentity.ObjectName
}

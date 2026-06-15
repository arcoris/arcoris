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

package objectlifecycle

import (
	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

// ListRequest identifies one resource collection, structural scope, and
// optional semantic query to read.
type ListRequest struct {
	// Resource is the concrete resource collection identity to resolve.
	Resource identity.GroupVersionResource

	// Scope is the explicit structural storage collection scope.
	Scope objectstore.ListScope

	// Query filters already-loaded list items above objectstore.
	//
	// Resource and Scope define the storage read sent to objectstore.List. Query
	// is compiled by objectlifecycle, evaluated after the store result is cloned,
	// and never pushed into storage. It does not affect ListResult.Revision.
	Query objectquery.Query
}

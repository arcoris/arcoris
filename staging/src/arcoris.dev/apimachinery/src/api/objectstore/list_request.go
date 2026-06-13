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

import apiidentity "arcoris.dev/apimachinery/api/identity"

// ListRequest identifies one resource collection read.
type ListRequest struct {
	// Resource is the concrete versioned resource collection to list.
	Resource apiidentity.GroupVersionResource

	// Scope is the structural namespace filter to apply to object keys.
	Scope ListScope
}

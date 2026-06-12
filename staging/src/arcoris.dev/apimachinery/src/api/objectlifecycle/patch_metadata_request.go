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
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
)

// PatchMetadataRequest patches labels and annotations on one object.
type PatchMetadataRequest struct {
	// Resource is the concrete resource collection identity to resolve.
	Resource identity.GroupVersionResource

	// Object is the namespace/name identity to patch.
	Object metaidentity.ObjectName

	// Labels maps label keys to replacement values. A nil value deletes the key.
	Labels map[string]*string

	// Annotations maps annotation keys to replacement values. A nil value deletes the key.
	Annotations map[string]*string

	// Owner receives ownership of every patched label and annotation key.
	Owner fieldownership.Owner

	// Expected is the committed live revision that must match before patch.
	Expected objectstore.Revision
}

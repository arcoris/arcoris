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

package objectsurface

import "arcoris.dev/apimachinery/api/apidocument"

// MetadataKindTree is an immutable navigation view over ObjectMeta surfaces.
type MetadataKindTree struct{}

// Labels returns the ObjectMeta labels map surface.
func (MetadataKindTree) Labels() Kind {
	return metadataKind(apidocument.Fields().ObjectMeta().Labels())
}

// Annotations returns the ObjectMeta annotations map surface.
func (MetadataKindTree) Annotations() Kind {
	return metadataKind(apidocument.Fields().ObjectMeta().Annotations())
}

// Finalizers returns the reserved finalizers lifecycle surface.
func (MetadataKindTree) Finalizers() Kind {
	return metadataKind(apidocument.Fields().ObjectMeta().Finalizers())
}

// OwnerReferences returns the reserved ownerReferences governance surface.
func (MetadataKindTree) OwnerReferences() Kind {
	return metadataKind(apidocument.Fields().ObjectMeta().OwnerReferences())
}

// metadataKind converts an ObjectMeta field into an object-root-relative
// metadata surface identifier.
func metadataKind(field apidocument.FieldName) Kind {
	return Kind(
		apidocument.Path("").
			Child(apidocument.Fields().Object().Metadata()).
			Child(field).
			String(),
	)
}

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

// KindTree is an immutable navigation view over known object surfaces.
//
// Surface IDs are relative to the object document root. For example, the
// apidocument object path object.metadata.labels maps to the surface ID
// metadata.labels.
type KindTree struct{}

// Kinds returns an immutable tree of known object surface identifiers.
func Kinds() KindTree {
	return KindTree{}
}

// Desired returns the declarative user or manager intent surface.
func (KindTree) Desired() Kind {
	return rootKind(apidocument.Fields().Object().Desired())
}

// Observed returns the runtime, controller, or agent truth surface.
func (KindTree) Observed() Kind {
	return rootKind(apidocument.Fields().Object().Observed())
}

// Metadata returns known ObjectMeta surfaces.
func (KindTree) Metadata() MetadataKindTree {
	return MetadataKindTree{}
}

// rootKind converts a top-level object field into an object-root-relative
// surface identifier.
func rootKind(field apidocument.FieldName) Kind {
	return Kind(apidocument.Path("").Child(field).String())
}

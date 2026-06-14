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

// KindFromObjectPath returns the known surface identified by path.
//
// The input path must be a full apidocument object path such as
// object.metadata.labels. Paths to TypeMeta, ObjectMeta identity fields, and
// unknown object fields are intentionally not mapped to generic surfaces.
func KindFromObjectPath(path apidocument.Path) (kind Kind, ok bool) {
	kinds := Kinds()
	metadata := kinds.Metadata()
	paths := apidocument.Paths().Object()

	switch path {
	case paths.Desired():
		return kinds.Desired(), true
	case paths.Observed():
		return kinds.Observed(), true
	case paths.Metadata().Labels():
		return metadata.Labels(), true
	case paths.Metadata().Annotations():
		return metadata.Annotations(), true
	case paths.Metadata().Finalizers():
		return metadata.Finalizers(), true
	case paths.Metadata().OwnerReferences():
		return metadata.OwnerReferences(), true
	default:
		return "", false
	}
}

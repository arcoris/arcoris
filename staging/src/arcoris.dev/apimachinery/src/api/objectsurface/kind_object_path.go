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

// ObjectPath returns the canonical API document path for k.
//
// ObjectPath reports ok=false for unknown surfaces. Known surface IDs are
// object-root-relative; returned paths are full apidocument object paths.
func (k Kind) ObjectPath() (path apidocument.Path, ok bool) {
	kinds := Kinds()
	metadata := kinds.Metadata()
	paths := apidocument.Paths().Object()

	switch k {
	case kinds.Desired():
		return paths.Desired(), true
	case kinds.Observed():
		return paths.Observed(), true
	case metadata.Labels():
		return paths.Metadata().Labels(), true
	case metadata.Annotations():
		return paths.Metadata().Annotations(), true
	case metadata.Finalizers():
		return paths.Metadata().Finalizers(), true
	case metadata.OwnerReferences():
		return paths.Metadata().OwnerReferences(), true
	default:
		return "", false
	}
}

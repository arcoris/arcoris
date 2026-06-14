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

package apidocument

// ObjectPaths groups paths for object envelope documents.
type ObjectPaths struct {
	// base is the document path root for this object view.
	base Path
}

// APIVersion returns the top-level object apiVersion path.
func (p ObjectPaths) APIVersion() Path {
	return p.base.Child(Fields().Object().APIVersion())
}

// Kind returns the top-level object kind path.
func (p ObjectPaths) Kind() Path {
	return p.base.Child(Fields().Object().Kind())
}

// MetadataPath returns the top-level object metadata path.
func (p ObjectPaths) MetadataPath() Path {
	return p.base.Child(Fields().Object().Metadata())
}

// Metadata returns paths below the top-level object metadata field.
func (p ObjectPaths) Metadata() ObjectMetaPaths {
	return ObjectMetaPaths{base: p.MetadataPath()}
}

// Desired returns the top-level object desired path.
func (p ObjectPaths) Desired() Path {
	return p.base.Child(Fields().Object().Desired())
}

// Observed returns the top-level object observed path.
func (p ObjectPaths) Observed() Path {
	return p.base.Child(Fields().Object().Observed())
}

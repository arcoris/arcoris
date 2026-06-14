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

// OwnershipPaths groups paths for object ownership state documents.
type OwnershipPaths struct {
	// base is the document path root for this ownership view.
	base Path
}

// Desired returns the Desired ownership surface path.
func (p OwnershipPaths) Desired() Path {
	return p.base.Child(Fields().Ownership().Desired())
}

// DesiredSurface returns paths below the Desired ownership surface object.
func (p OwnershipPaths) DesiredSurface() OwnershipSurfacePaths {
	return OwnershipSurfacePaths{base: p.Desired()}
}

// Observed returns the Observed ownership surface path.
func (p OwnershipPaths) Observed() Path {
	return p.base.Child(Fields().Ownership().Observed())
}

// ObservedSurface returns paths below the Observed ownership surface object.
func (p OwnershipPaths) ObservedSurface() OwnershipSurfacePaths {
	return OwnershipSurfacePaths{base: p.Observed()}
}

// MetadataPath returns the metadata ownership object path.
func (p OwnershipPaths) MetadataPath() Path {
	return p.base.Child(Fields().Ownership().MetadataField())
}

// Metadata returns paths below the metadata ownership object.
func (p OwnershipPaths) Metadata() OwnershipMetadataPaths {
	return OwnershipMetadataPaths{base: p.MetadataPath()}
}

// OwnershipMetadataPaths groups paths inside the metadata ownership object.
type OwnershipMetadataPaths struct {
	// base is the metadata ownership object path.
	base Path
}

// Labels returns the metadata.labels ownership surface path.
func (p OwnershipMetadataPaths) Labels() Path {
	return p.base.Child(Fields().Ownership().Metadata().Labels())
}

// LabelsSurface returns paths below the metadata.labels ownership surface object.
func (p OwnershipMetadataPaths) LabelsSurface() OwnershipSurfacePaths {
	return OwnershipSurfacePaths{base: p.Labels()}
}

// Annotations returns the metadata.annotations ownership surface path.
func (p OwnershipMetadataPaths) Annotations() Path {
	return p.base.Child(Fields().Ownership().Metadata().Annotations())
}

// AnnotationsSurface returns paths below the metadata.annotations ownership surface object.
func (p OwnershipMetadataPaths) AnnotationsSurface() OwnershipSurfacePaths {
	return OwnershipSurfacePaths{base: p.Annotations()}
}

// OwnershipSurfacePaths groups paths inside one ownership surface object.
type OwnershipSurfacePaths struct {
	// base is the ownership surface object path.
	base Path
}

// Entries returns the entries path inside one ownership surface object.
func (p OwnershipSurfacePaths) Entries() Path {
	return p.base.Child(Fields().Ownership().Surface().Entries())
}

// Entry returns template paths below one ownership entry object.
func (p OwnershipSurfacePaths) Entry() OwnershipEntryPaths {
	return OwnershipEntryPaths{base: p.Entries()}
}

// OwnershipEntryPaths groups paths inside one ownership entry object.
type OwnershipEntryPaths struct {
	// base is the ownership entry object path.
	base Path
}

// Owner returns the owner path inside one ownership entry object.
func (p OwnershipEntryPaths) Owner() Path {
	return p.base.Child(Fields().Ownership().Entry().Owner())
}

// Fields returns the fields path inside one ownership entry object.
func (p OwnershipEntryPaths) Fields() Path {
	return p.base.Child(Fields().Ownership().Entry().Fields())
}

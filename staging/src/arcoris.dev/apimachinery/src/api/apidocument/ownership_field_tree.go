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

// OwnershipFields groups field names for object ownership state documents.
type OwnershipFields struct{}

// Desired returns the Desired ownership surface field.
func (OwnershipFields) Desired() FieldName {
	return OwnershipFieldDesired
}

// Observed returns the Observed ownership surface field.
func (OwnershipFields) Observed() FieldName {
	return OwnershipFieldObserved
}

// MetadataField returns the metadata ownership object field.
func (OwnershipFields) MetadataField() FieldName {
	return OwnershipFieldMetadata
}

// Metadata returns fields inside the metadata ownership object.
func (OwnershipFields) Metadata() OwnershipMetadataFields {
	return OwnershipMetadataFields{}
}

// Surface returns fields common to each emitted ownership surface.
func (OwnershipFields) Surface() OwnershipSurfaceFields {
	return OwnershipSurfaceFields{}
}

// Entry returns fields inside one ownership entry.
func (OwnershipFields) Entry() OwnershipEntryFields {
	return OwnershipEntryFields{}
}

// OwnershipMetadataFields groups fields inside the metadata ownership object.
type OwnershipMetadataFields struct{}

// Labels returns the metadata.labels ownership surface field.
func (OwnershipMetadataFields) Labels() FieldName {
	return OwnershipFieldLabels
}

// Annotations returns the metadata.annotations ownership surface field.
func (OwnershipMetadataFields) Annotations() FieldName {
	return OwnershipFieldAnnotations
}

// OwnershipSurfaceFields groups fields common to each ownership surface object.
type OwnershipSurfaceFields struct{}

// Entries returns the ownership surface entries field.
func (OwnershipSurfaceFields) Entries() FieldName {
	return OwnershipFieldEntries
}

// OwnershipEntryFields groups fields inside one ownership entry object.
type OwnershipEntryFields struct{}

// Owner returns the ownership entry owner field.
func (OwnershipEntryFields) Owner() FieldName {
	return OwnershipFieldOwner
}

// Fields returns the ownership entry fields field.
func (OwnershipEntryFields) Fields() FieldName {
	return OwnershipFieldFields
}

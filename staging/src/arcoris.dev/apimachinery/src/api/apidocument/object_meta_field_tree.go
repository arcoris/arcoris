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

// ObjectMetaFields groups field names for object metadata documents.
type ObjectMetaFields struct{}

// Name returns the object metadata name field.
func (ObjectMetaFields) Name() FieldName {
	return ObjectMetaFieldName
}

// GenerateName returns the object metadata generateName field.
func (ObjectMetaFields) GenerateName() FieldName {
	return ObjectMetaFieldGenerateName
}

// Namespace returns the object metadata namespace field.
func (ObjectMetaFields) Namespace() FieldName {
	return ObjectMetaFieldNamespace
}

// UID returns the object metadata uid field.
func (ObjectMetaFields) UID() FieldName {
	return ObjectMetaFieldUID
}

// ResourceVersion returns the object metadata resourceVersion field.
func (ObjectMetaFields) ResourceVersion() FieldName {
	return ObjectMetaFieldResourceVersion
}

// Generation returns the object metadata generation field.
func (ObjectMetaFields) Generation() FieldName {
	return ObjectMetaFieldGeneration
}

// CreatedAt returns the object metadata createdAt field.
func (ObjectMetaFields) CreatedAt() FieldName {
	return ObjectMetaFieldCreatedAt
}

// Deletion returns the object metadata deletion field.
func (ObjectMetaFields) Deletion() FieldName {
	return ObjectMetaFieldDeletion
}

// Labels returns the object metadata labels field.
func (ObjectMetaFields) Labels() FieldName {
	return ObjectMetaFieldLabels
}

// Annotations returns the object metadata annotations field.
func (ObjectMetaFields) Annotations() FieldName {
	return ObjectMetaFieldAnnotations
}

// OwnerReferences returns the object metadata ownerReferences field.
func (ObjectMetaFields) OwnerReferences() FieldName {
	return ObjectMetaFieldOwnerReferences
}

// Finalizers returns the object metadata finalizers field.
func (ObjectMetaFields) Finalizers() FieldName {
	return ObjectMetaFieldFinalizers
}

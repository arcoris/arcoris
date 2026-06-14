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

// ObjectFields groups field names for object envelope documents.
type ObjectFields struct{}

// APIVersion returns the top-level object apiVersion field.
func (ObjectFields) APIVersion() FieldName {
	return ObjectFieldAPIVersion
}

// Kind returns the top-level object kind field.
func (ObjectFields) Kind() FieldName {
	return ObjectFieldKind
}

// Metadata returns the top-level object metadata field.
func (ObjectFields) Metadata() FieldName {
	return ObjectFieldMetadata
}

// Desired returns the top-level object desired field.
func (ObjectFields) Desired() FieldName {
	return ObjectFieldDesired
}

// Observed returns the top-level object observed field.
func (ObjectFields) Observed() FieldName {
	return ObjectFieldObserved
}

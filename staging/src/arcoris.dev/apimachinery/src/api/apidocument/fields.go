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

// FieldTree is an immutable navigation view over canonical API document fields.
//
// The package-level FieldName constants remain the single source of spelling
// truth. FieldTree groups those constants by the API document shape they belong
// to, so call sites can say Fields().Object().Desired() instead of scanning the
// whole flat constant set.
type FieldTree struct{}

// Fields returns an immutable tree of canonical API document fields.
func Fields() FieldTree {
	return FieldTree{}
}

// TypeMeta returns fields that belong to standalone TypeMeta documents.
func (FieldTree) TypeMeta() TypeMetaFields {
	return TypeMetaFields{}
}

// Object returns fields that belong to object envelope documents.
func (FieldTree) Object() ObjectFields {
	return ObjectFields{}
}

// ObjectMeta returns fields that belong to standalone ObjectMeta documents.
func (FieldTree) ObjectMeta() ObjectMetaFields {
	return ObjectMetaFields{}
}

// Ownership returns fields that belong to object ownership state documents.
func (FieldTree) Ownership() OwnershipFields {
	return OwnershipFields{}
}

// PageMeta returns fields that belong to list/page metadata documents.
func (FieldTree) PageMeta() PageMetaFields {
	return PageMetaFields{}
}

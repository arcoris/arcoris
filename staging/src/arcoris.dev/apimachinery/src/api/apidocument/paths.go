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

const (
	pathRootTypeMeta   Path = "typeMeta"
	pathRootObject     Path = "object"
	pathRootObjectMeta Path = "objectMeta"
	pathRootOwnership  Path = "ownership"
	pathRootPageMeta   Path = "pageMeta"
)

// PathTree is an immutable navigation view over canonical API document paths.
type PathTree struct{}

// Paths returns an immutable tree of stable API document diagnostic paths.
func Paths() PathTree {
	return PathTree{}
}

// TypeMeta returns paths for standalone TypeMeta documents.
func (PathTree) TypeMeta() TypeMetaPaths {
	return TypeMetaPaths{base: pathRootTypeMeta}
}

// Object returns paths for object envelope documents.
func (PathTree) Object() ObjectPaths {
	return ObjectPaths{base: pathRootObject}
}

// ObjectMeta returns paths for standalone ObjectMeta documents.
func (PathTree) ObjectMeta() ObjectMetaPaths {
	return ObjectMetaPaths{base: pathRootObjectMeta}
}

// Ownership returns paths for object ownership state documents.
func (PathTree) Ownership() OwnershipPaths {
	return OwnershipPaths{base: pathRootOwnership}
}

// PageMeta returns paths for list/page metadata documents.
func (PathTree) PageMeta() PageMetaPaths {
	return PageMetaPaths{base: pathRootPageMeta}
}

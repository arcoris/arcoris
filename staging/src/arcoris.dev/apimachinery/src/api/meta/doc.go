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

// Package meta defines metadata primitives and metadata structs for ARCORIS API
// objects.
//
// The package owns TypeMeta, ObjectMeta, ListMeta, PageToken, and metadata
// subpackages for object identity, stamps, labels, annotations, owner
// references, and finalizers. Concrete resource packages embed or compose these
// metadata values and define their own desired and observed payload fields.
//
// api/meta is a metadata foundation layer, not a full object model. It does not
// define resource contracts, structural descriptors, selectors, status
// conditions, storage keys, watches, admission, patch/apply, codecs, runtime
// schemes, or global registries. api/resource owns resource contracts and
// desired/observed descriptors; api/meta remains independent from that layer.
package meta

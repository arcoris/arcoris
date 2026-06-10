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

// Package meta defines the canonical metadata semantics used by ARCORIS API
// objects.
//
// The package owns TypeMeta, ObjectMeta, ListMeta, PageToken, and metadata
// subpackages for object identity, stamps, labels, annotations, owner
// references, and finalizers. Concrete resource packages embed or compose these
// metadata values and define their own desired and observed payload fields.
// api/meta is the single source of truth for what these metadata values mean
// and how their lexical forms are validated.
//
// api/apidocument is the separate source of truth for canonical API document
// field spelling such as apiVersion, metadata, resourceVersion, and continue.
// Metadata structs carry literal JSON tags that mirror those document field
// constants because Go struct tags cannot reference constants. Drift is enforced
// by api/apidocument tests, not by importing api/apidocument from production
// metadata code.
//
// The package is versioned with the Go module. It is not itself an API
// group-version package, and it should not be copied into api/meta/v1 without a
// real conversion and compatibility design. Versioned resource packages may
// embed these metadata structs in their external object types.
//
// Scalar metadata types may provide small JSON/text helpers where the scalar
// owns its own wire contract, but this package does not implement object
// document codecs, envelope parsing, runtime schemes, conversion, admission,
// storage validation, or apply behavior.
//
// Metadata validation is lexical metadata validation, not lifecycle-phase
// validation. ObjectMeta.Validate and ListMeta.Validate do not decide whether
// metadata is valid for create, update, storage, apply, delete, list response,
// or pagination phases.
//
// api/meta is a metadata foundation layer, not a full object model. It does not
// define resource contracts, structural descriptors, selectors, status
// conditions, storage keys, watches, admission, patch/apply, codecs, runtime
// schemes, or global registries. api/resource owns resource contracts and
// desired/observed descriptors; api/meta remains independent from that layer.
package meta

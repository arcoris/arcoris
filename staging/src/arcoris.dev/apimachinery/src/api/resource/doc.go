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

// Package resource defines versioned API resource definition descriptors.
//
// The package composes api/identity and api/types. Package identity owns the
// lexical API identity algebra: groups, versions, kinds, resources, and their
// canonical combinations. Package types owns structural value descriptors.
// Package resource binds those two layers into durable resource-family
// contracts.
//
// A Definition describes a resource family. It fixes the API group, kind,
// resource collection name, instance scope, and the set of versioned structural
// surfaces for that family. A VersionDefinition describes one API contract
// version of that family.
//
// The package intentionally uses ARCORIS terminology instead of Kubernetes CRD
// terminology:
//
//   - Desired is the user-authored/requested state surface. Kubernetes-like
//     adapters may later map it to spec, while SDK/IaC adapters may map it to
//     input arguments.
//   - Observed is the optional system-computed/read surface. Kubernetes-like
//     adapters may later map it to status, while SDK/IaC adapters may map it to
//     output properties.
//   - Exposed marks an API version as part of the public API surface.
//   - Canonical marks the single canonical version of a resource family.
//     Canonical does not define persistence or storage behavior.
//   - ScopeGlobal and ScopeNamespaced describe future object-addressing scope
//     without defining object metadata.
//
// Resource definitions are descriptor contracts only. This package does not
// define live objects, object metadata, names, namespaces, UIDs, labels,
// annotations, resource versions, generations, storage, watches, REST routes,
// codecs, JSON/YAML manifest formats, OpenAPI export, discovery documents,
// conversion, defaulting, pruning, admission, field ownership, controllers, or
// provider lifecycle handlers.
//
// # Model
//
// The core model is deliberately small:
//
//	Definition
//	  Group + Kind + Resource + Scope + []VersionDefinition
//
//	VersionDefinition
//	  Version + Desired Descriptor + optional Observed Descriptor + Exposed/Canonical flags
//
// Definitions validate their identity values, version set, exposed/canonical
// invariants, and structural surface descriptors. Desired and Observed root
// descriptors MUST be objects or references resolving to objects. Reference
// roots require a non-nil types.Resolver so resource validation can prove the
// root object shape.
//
// # Construction
//
// Constructors keep required fields visible and leave validation as an explicit
// trust boundary:
//
//	worker := resource.NewDefinition(
//		identity.Group("control.arcoris.dev"),
//		identity.Kind("Worker"),
//		identity.Resource("workers"),
//		resource.ScopeNamespaced,
//		resource.NewVersion(
//			identity.Version("v1"),
//			types.Ref("control.arcoris.dev.WorkerDesired").Descriptor(),
//			resource.Observed(
//				types.Ref("control.arcoris.dev.WorkerObserved").Descriptor(),
//			),
//			resource.Exposed(),
//			resource.Canonical(),
//		),
//	)
//
//	_ = worker
//
// Concrete storage for resource definitions is intentionally outside this
// package. A future resource catalog can index Definitions by GVK/GVR without
// making this package a global registry.
//
// # File Orientation
//
// The package is split by descriptor responsibility rather than by call path.
// Definition accessors live with Definition. VersionDefinition accessors live
// with VersionDefinition. Scope parsing, validation, and scalar encoding each
// have their own files because those operations are distinct trust boundaries.
//
// Validation follows the same rule. The top-level definition validator only
// coordinates validation. Identity validation checks group/kind/resource/scope.
// Version-set validation checks duplicate/exposed/canonical invariants.
// Version validation checks the single version descriptor, then delegates
// Desired/Observed surface checks to the surface validator. Object-like probing
// is isolated because references need resolver-aware traversal.
//
// Error files are similarly separated. The Error type and broad sentinels live
// together; precise reason constants and private construction helpers live in
// their own files. This keeps diagnostics consistent while avoiding hard-coded
// paths or repeated string fragments across validators.
//
// Tests mirror this layout: each production file has a matching test file for
// its responsibility, with shared test-only helpers kept in test_helpers_test.go.
package resource

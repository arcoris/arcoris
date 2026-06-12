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

// Package resource defines durable API resource-family descriptor contracts.
//
// Package identity owns lexical API identity. Package types owns structural
// value descriptors. Package resource binds those two layers into durable
// resource-family contracts. Package resourcecatalog owns owner-created
// cataloging, indexing, and lookup of those contracts. Object validation and
// apply packages consume resource definitions; future codec, storage,
// discovery, and runtime layers may also consume them without becoming
// resource responsibilities.
//
// A Definition describes one resource family:
//
//   - Group
//   - Kind
//   - Resource
//   - Scope
//   - []VersionDefinition
//
// A VersionDefinition describes one versioned contract:
//
//   - Version
//   - Desired descriptor
//   - optional Observed descriptor
//   - Exposed flag
//   - Canonical flag
//
// Desired is the user-authored/requested state surface. Observed is the
// optional system-computed/read surface. Both are descriptor surfaces and both
// must be object-like when validation is resolved. Desired/Observed root
// DescriptorRef values are accepted by local validation when their ref syntax
// is valid; resolved validation additionally requires them to resolve to
// object-like descriptors through a resolver.
//
// Exposed means the version is part of the public API surface. Canonical means
// the family's single canonical descriptor version for
// version-independent consumers. Canonical does not imply storage version,
// conversion hub, persistence format, migration target, or defaulting behavior.
// A Definition represents an API-facing resource family, so validation requires
// at least one exposed version and exactly one canonical version. Internal-only
// structural descriptors belong in api/types and api/typecatalog, not here.
//
// ScopeGlobal and ScopeNamespaced describe future instance-addressing scope.
// They do not define metadata, namespace objects, routing, authorization,
// tenancy, storage partitioning, or discovery behavior.
//
// Local validation checks resource identity, scope, version-set invariants, and
// descriptor-local surface structure. Resolved validation does the same but
// also resolves DescriptorRef roots and nested refs through a resolver.
//
// Resource definitions are descriptor contracts only. This package does not
// define live objects, object metadata, names, namespaces, UIDs, labels,
// annotations, resource versions, generations, storage, watches, REST routes,
// codecs, JSON/YAML definition formats, OpenAPI export, discovery documents,
// conversion, defaulting, pruning, admission, field ownership, controllers, or
// provider lifecycle handlers.
package resource

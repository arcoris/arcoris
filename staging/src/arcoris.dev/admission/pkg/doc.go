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

// Package admission defines generic admission-result contracts for ARCORIS
// control paths.
//
// Admission is the boundary at which a component decides whether work may enter
// the next execution or control path. The package provides a small algebra for
// those decisions: closed outcome and side-effect vocabularies, open-world
// reason and component identifiers, typed generic results, and generic admitter
// contracts.
//
// The package owns the common language of admission results and owner-created
// catalogs for stable admission metadata. It does not own capacity accounting,
// leases, retry budgets, token buckets, queues, workers, schedulers, health
// state, metrics, logging, tracing, distributed coordination, dynamic service
// discovery, or rollback chains. Domain packages compose Result with their own
// request, grant, and metadata types.
//
// # Open-world model
//
// Outcome and Effect are closed because Result validity depends on their exact
// semantics. Reason, ComponentID, and ComponentKind are open-world string value
// types: packages may define their own stable values without changing this core
// package, provided those values satisfy the package syntax rules and do not
// contain dynamic data such as tenant IDs, request IDs, timestamps, addresses,
// raw errors, secrets, or resource-specific identifiers.
//
// ComponentID identifies a component kind or role, not a runtime instance. For
// example, resilience.bulkhead and resilience.retrybudget are valid component
// IDs; resilience.bulkhead.tenant_123 is not. Instance naming, metrics labels,
// tracing attributes, health state, config ownership, and dynamic discovery
// belong to higher-level owners.
//
// # Catalog registries
//
// ReasonRegistry catalogs stable ReasonDescriptor values. KindRegistry catalogs
// stable ComponentKindDescriptor values. ComponentRegistry catalogs stable
// ComponentDescriptor values and validates them against a KindRegistry. Catalog
// aggregates those owner-created registries behind one metadata boundary.
// Catalog construction rejects a ComponentRegistry backed by a different
// KindRegistry than the one passed to NewCatalog. That reference check
// guarantees RegisterKind and RegisterComponent operate over one consistent
// kind catalog.
//
// Built-in catalog constructors return ordinary owner-created registries. They
// do not create global mutable state. Admission does not provide package-level
// Register functions, init-time registration, global mutable registries, or
// process-wide singleton catalogs.
//
// Catalog registries are metadata catalogs, not enforcement engines. Result
// validity still depends on Decision, Effect, and grant-shape invariants.
// Registry membership can support deterministic listing, docs/tests/config
// validation, and future chain validation foundations, but ordinary Result
// construction does not require a registry lookup.
//
// Runtime instances, live admitters, chains, queues, health state, metrics,
// tracing attributes, dynamic discovery, config ownership, scheduling, and
// service discovery remain outside admission.
//
// # Side effects
//
// Effect records the side-effect semantics of an admission decision. This is
// deliberately separate from Outcome. An admitted retry budget may commit a
// spend-only side effect, while an admitted bulkhead may return a caller-owned
// lease that must later be released. Chains that compose multiple effectful
// stages must release, cancel, commit, roll back, or reconcile domain grants when
// later stages deny or defer an attempt. This package validates result shape but
// does not implement generic rollback.
package admission

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

// Package resourcecatalog provides owner-created catalogs for API resource
// definition descriptors.
//
// The package stores api/resource Definition values and indexes them by
// GroupResource, GroupKind, GroupVersionResource, and GroupVersionKind. It is
// the resource-level analogue of api/typecatalog: a small descriptor
// composition helper with stable enumeration order, atomic registration, and no
// package-level global state.
//
// The package stores API descriptors only. It deliberately does not bind
// descriptors to Go runtime types, constructors, codecs, handlers, storage
// backends, watches, routes, discovery documents, OpenAPI exporters,
// controllers, provider lifecycle hooks, or process-wide extension state.
//
// Catalogs validate resource definitions at registration time. Definitions that
// use DescriptorRef roots are validated against the types.Resolver supplied to New.
// A zero-value Catalog is usable and behaves as if it had a nil type resolver.
package resourcecatalog

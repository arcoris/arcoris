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

// Package typecatalog provides an owner-created catalog of named structural
// type definitions.
//
// The catalog is the concrete storage counterpart to the Resolver contract in
// package arcoris.dev/apimachinery/api/types. Package types owns DescriptorRef and
// therefore owns the minimal resolution interface. Package typecatalog owns the
// mutable, concurrency-safe collection used by descriptor authors who want a
// local catalog of Definition values.
//
// Catalog is an API descriptor composition helper. It is not a runtime scheme,
// codec registry, resource registry, storage registry, discovery cache, or
// global extension registry. It has no package-level singleton and performs no
// init-time registration, so callers keep catalog scope, versioning, and test
// fixtures explicit.
package typecatalog

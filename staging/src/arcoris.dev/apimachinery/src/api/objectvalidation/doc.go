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

// Package objectvalidation validates API object envelopes against already
// resolved resource contracts.
//
// The package checks contract conformance, not request admissibility. It
// coordinates object metadata validation, object GVK to resource group/kind
// matching, resource version lookup, minimal scope compatibility, and desired
// or observed surface validation through typed SurfaceValidator implementations.
//
// Resource definitions supplied through Plan are expected to have been
// validated at construction, registration, or catalog boundaries. This package
// performs defensive plan-shape checks, but it does not repeatedly validate the
// whole resource descriptor graph for every object.
//
// The package deliberately does not perform catalog lookup, admission,
// defaulting, conversion, pruning, mutation, storage validation, serving
// validation, selector matching, status or subresource handling, runtime scheme
// registration, codecs, clients, controllers, or global registration.
package objectvalidation

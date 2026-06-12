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

// Package objectvalidation validates generic object envelopes against already
// resolved resource contracts.
//
// Package object defines generic envelopes and metadata-only validation.
// Package resource defines resource-family contracts and version descriptors.
// Package types defines structural descriptors. Package valuevalidation can
// implement SurfaceValidator[value.Value]. Package objectapply uses this
// package for value-backed apply validation. Package objectstore stores
// committed state but does not validate resource contracts.
//
// Validation returns the first failure in this deterministic order:
//
//   - static plan shape
//   - object metadata
//   - object group/kind match against the resource family
//   - object API version lookup in the resource definition
//   - minimal scope compatibility
//   - desired surface validation
//   - observed surface validation
//
// Plan.Resource must already be resolved and prevalidated by construction,
// registration, or catalog code. Plan.Resolver is passed through to surface
// validators. Plan.DesiredValidator is required. Plan.ObservedValidator is
// required only when the selected resource version defines an observed
// descriptor and the object actually carries observed data.
//
// Scope validation is intentionally minimal. Global resources must not carry a
// namespace. Namespaced resources may carry a namespace or omit it; requiring a
// namespace belongs to request admission, serving, storage, or lifecycle
// layers.
//
// Version validation only checks that the object's apiVersion exists in the
// resource definition. It does not require the version to be exposed,
// canonical, served, preferred, or storage-backed.
//
// Desired is always validated against the selected resource version's desired
// descriptor. Observed is optional. Observed is rejected when the selected
// version does not define an observed descriptor, and validated only when
// present.
//
// Non-goals: no catalog lookup, admission, serving/exposed-version policy,
// storage validation, defaulting, conversion, pruning, apply, ownership,
// lifecycle, codecs, runtime schemes, clients, controllers, or global
// registration.
package objectvalidation

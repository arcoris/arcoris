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

// Package objectvalidation validates API object envelopes against resolved
// resource contracts.
//
// The package sits between api/object, api/resource, and api/types. It checks
// that an object carries valid metadata, matches the resource family group and
// kind, uses a version defined by that resource, satisfies minimal scope rules,
// and delegates desired/observed payload checks to explicit typed surface
// validators.
//
// The package deliberately works with an already-resolved resource.Definition.
// It does not perform catalog lookup, defaulting, conversion, pruning, storage
// validation, selector matching, operation validation, runtime scheme work,
// codecs, controllers, clients, or global registration.
package objectvalidation

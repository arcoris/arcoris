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

// Package object defines generic ARCORIS API object envelopes.
//
// The package composes api/meta TypeMeta, ObjectMeta, and PageMeta with
// resource-specific desired and observed payload values. It is intentionally
// generic: desired and observed values may be handwritten structs, generated
// structs, decoded adapter values, or test fixtures.
//
// api/object validates metadata only through ValidateMeta. It does not validate
// desired or observed payloads because it does not know the resource contract or
// structural descriptors for those surfaces. Resource contracts live in
// arcoris.dev/apimachinery/api/resource. Structural descriptors live in
// arcoris.dev/apimachinery/api/types.
//
// A future validation package may bind Object[D, O] values to resource
// definitions and type descriptors. That layer is intentionally outside
// api/object so this package remains a small reusable envelope model.
//
// The package is not runtime machinery. It does not implement storage keys,
// watches, admission, defaulting, conversion, selectors, status conventions,
// codecs, runtime schemes, clients, informers, controllers, catalogs, or global
// registries.
package object

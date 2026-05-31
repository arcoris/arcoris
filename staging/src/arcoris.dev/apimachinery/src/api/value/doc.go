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

// Package value defines concrete ARCORIS API payload values.
//
// The package is the value half of the descriptor/value split:
// arcoris.dev/apimachinery/api/types describes allowed shapes and constraints,
// while api/value stores actual payload data. This package deliberately does
// not import descriptor packages and does not validate values against type
// descriptors. Descriptor-aware validation belongs to a future valuevalidation
// package.
//
// Value is immutable by API convention. Constructors and accessors clone
// mutable data such as byte slices and composite payloads, so callers cannot
// mutate stored values through returned slices, maps, or views. The zero Value
// is invalid and represents missing initialization; Null is an explicit API
// value.
//
// Object and map views use linear lookup rather than storing lookup indexes in
// payloads. API values are expected to be small, and this keeps construction,
// cloning, and view creation allocation-light with fewer invariants to maintain.
//
// The package does not implement JSON or YAML codecs, Go-struct introspection,
// object/resource validation, defaulting, pruning, conversion, admission,
// storage, patch/apply behavior, selectors, status conventions, runtime
// schemes, code generation, or catalog lookup.
package value

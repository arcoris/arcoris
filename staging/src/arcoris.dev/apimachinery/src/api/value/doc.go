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
// Value is immutable by API convention. Constructors clone mutable inputs, and
// accessors clone mutable outputs such as byte slices and nested composite
// values. The zero Value is invalid and represents missing initialization; Null
// is an explicit API value.
//
// Constructors that return Value use the <Kind>Value naming form, such as
// StringValue, ObjectValue, and ListValue. Constructors that return supporting
// domain types keep New<Type> names, such as NewDate and NewDecimal. Must
// variants exist only for fallible Value or domain constructors.
//
// Value Object is a concrete keyed payload node. It does not decide whether the
// payload is a descriptor object or descriptor map. Descriptor-aware validation
// interprets the same concrete payload according to the expected api/types.Type.
//
// Object views use linear lookup rather than storing lookup indexes in payloads.
// API values are expected to be small, and this keeps construction, cloning, and
// view creation allocation-light with fewer invariants to maintain. Empty object
// and list views return non-nil empty slices from bulk accessors even though
// internal empty payload storage may be nil.
//
// Object member names are concrete non-empty strings. Empty name rejection is a
// base value grammar invariant; name regexes, prefixes, and semantic key
// constraints belong to descriptor-aware validation outside this package.
//
// The package does not implement JSON or YAML codecs, Go-struct introspection,
// object/resource validation, defaulting, pruning, conversion, admission,
// storage, patch/apply behavior, selectors, status conventions, runtime
// schemes, code generation, or catalog lookup.
package value

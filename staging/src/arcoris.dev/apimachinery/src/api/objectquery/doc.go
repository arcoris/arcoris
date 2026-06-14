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

// Package objectquery defines typed, format-neutral query predicates for API
// object list items.
//
// The package evaluates already-loaded objectstore.ListItem values. It does
// not read from stores, push filters into storage, parse selector strings,
// decode HTTP query parameters, build indexes, watch changes, run controllers,
// validate resource descriptors, or perform admission or authorization.
//
// Query is the declarative value callers build. Compile validates and
// canonicalizes Query into Predicate, which is the deterministic evaluator.
// Query sections use AND semantics: identity, labels, and annotations must all
// match. Label and annotation selector requirements are also ANDed. There is no
// OR, nested boolean grouping, regex, numeric comparison, or Desired/Observed
// field traversal in v1.
//
// Negative metadata requirements intentionally match absent keys. NotEquals is
// true when the key is absent or has a different value. NotIn is true when the
// key is absent or its value is outside the set.
//
// Predicate.Filter preserves input order and does not clone item state. Store
// and list result APIs own detachment; objectquery is only a pure selection
// layer.
package objectquery

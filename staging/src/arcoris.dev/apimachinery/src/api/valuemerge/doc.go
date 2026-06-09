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

// Package valuemerge merges concrete API payload values under api/types
// descriptor semantics.
//
// The package is a pure value-level merge primitive. It copies only selected
// semantic field paths from an overlay value into a base value and returns the
// merged value. Selection is expressed as an api/fieldpath.Set.
//
// A selected path replaces the overlay subtree at that path. A selected
// descendant path recursively merges only that descendant. An unselected path
// preserves the base value. A selected path absent from the overlay removes the
// corresponding base field where the descriptor node supports removal.
//
// The package is descriptor-aware and path-aware, but policy-free. It does not
// validate complete values, extract field sets, compare values, detect ownership
// conflicts, update field ownership, apply configurations, authorize requests,
// perform admission, decode wire formats, normalize values, apply defaults,
// access storage, or validate object/resource metadata.
//
// Callers are expected to validate base and overlay values with
// api/valuevalidation before merge when full descriptor conformance is required.
// valuemerge performs only the defensive checks needed for selected traversal
// and replacement shape. Exact subtree replacement checks zero values, DescriptorRef
// resolution, and descriptor/value kind compatibility; it does not fully
// validate the replacement subtree.
//
// Higher layers decide which fields are selected for merge. For example, a
// future apply layer may combine valuevalidation, valuefieldset, valuecompare,
// fieldownership, and valuemerge.
package valuemerge

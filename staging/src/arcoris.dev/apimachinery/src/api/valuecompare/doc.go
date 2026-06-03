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

// Package valuecompare compares concrete ARCORIS API payload values under
// api/types descriptor semantics and reports semantic field changes.
//
// The package is descriptor-aware and path-aware, but side-effect free. It
// reads descriptors and immutable-by-convention payload values, then returns
// paths that were added, removed, or modified. It does not validate complete
// values, decode wire formats, normalize values, apply defaults, prune values,
// merge values, manage ownership, validate API object metadata, access storage,
// or perform admission.
//
// Compare results are expressed as added, removed, and modified api/fieldpath
// sets. Object fields, map keys, ordered list indexes, and ListMap selectors are
// compared according to descriptor semantics. Added and removed subtrees are
// discovered through api/valuefieldset so compare output stays aligned with
// future field ownership and apply planning layers.
package valuecompare

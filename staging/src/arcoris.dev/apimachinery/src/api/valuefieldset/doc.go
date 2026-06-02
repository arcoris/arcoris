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

// Package valuefieldset extracts semantic field paths mentioned by concrete
// ARCORIS API payload values under api/types descriptor semantics.
//
// The package is descriptor-aware but side-effect free. It does not perform
// full value validation, decode wire formats, normalize values, apply defaults,
// prune fields, compare values, merge values, manage ownership, validate API
// object metadata, access storage, or perform admission.
//
// Callers should validate payloads with api/valuevalidation before extracting
// field sets. This package performs only defensive checks needed to traverse
// values and build stable semantic paths.
//
// Extracted paths use api/fieldpath semantics: object descriptor fields use
// field elements, map entries use key elements, ordered list items use index
// elements, and ListMap items use selector elements.
//
// List extraction follows merge/ownership intent. Atomic and set-like lists
// produce the list path as one semantic field. Ordered lists produce index
// paths because item position is part of the descriptor contract. List-map
// values produce selector paths based on their declared identity fields.
//
// Unknown preserved object fields are included as opaque leaves because no
// descriptor exists for nested traversal. Unknown pruned fields are omitted, and
// rejected unknown fields fail extraction.
package valuefieldset

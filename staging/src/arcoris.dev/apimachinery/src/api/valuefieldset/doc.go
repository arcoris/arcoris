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

// Package valuefieldset extracts ownership-relevant semantic field paths from
// concrete ARCORIS API payload values under api/types descriptor semantics.
//
// The value pipeline keeps related but distinct responsibilities separate:
// api/valuevalidation validates values, valuefieldset extracts ownership fields,
// api/valuecompare computes semantic changes, api/valuemerge merges values,
// api/valueapply orchestrates apply and ownership updates, and
// api/fieldownership stores ownership state and detects ownership conflicts.
//
// Extracted fields are fields explicitly mentioned by the payload under
// descriptor semantics. They are suitable as applied ownership fields. They are
// not necessarily changed fields, conflict fields, merge fields, validation
// diagnostic paths, all possible leaf fields, or all fields in a live value.
//
// Explicit NullValue mentions the current semantic field. Nullability acceptance
// belongs to api/valuevalidation, not valuefieldset. Empty composite payloads
// also mention their composite container path: empty object, empty map, empty
// ordered list, empty ListMap, empty ListSet, and empty atomic list all extract
// the current path.
//
// Object descriptor fields produce fieldpath field elements. Unknown preserved
// fields produce opaque leaf paths because no descriptor exists for nested
// traversal. Unknown pruned fields are omitted. Unknown rejected fields fail
// extraction.
//
// Map entries produce fieldpath key elements. Concrete API maps are represented
// as string-keyed api/value records. Map key descriptor validation belongs to
// api/types and api/valuevalidation; valuefieldset performs only defensive
// descriptor checks needed for traversal.
//
// Malformed payload-derived record member names are reported as invalid values:
// this package extracts from a concrete payload and the malformed name is part
// of that payload. Malformed descriptor-declared field names remain invalid
// descriptors.
//
// List extraction follows ownership/apply intent. ListAtomic extracts the list
// path. ListSet also extracts the list path until a stable value-based set item
// ownership model exists. ListOrdered extracts index paths because position is
// part of the semantic contract. ListMap extracts selector paths from declared
// identity fields using api/internal/listmapkey.
//
// Descriptors are expected to have been validated at construction,
// registration, or catalog boundaries. valuefieldset does not call
// types.ValidateResolved on every extraction. It performs local defensive checks
// required for traversal and path construction. DescriptorRef values are
// resolved through Options.Resolver with MaxDepth recursion protection.
//
// The package does not perform full validation, scalar constraint checks,
// comparison, merging, ownership mutation, conflict detection, force behavior,
// defaulting, pruning, normalization, wire decoding, object metadata validation,
// storage access, authorization, or admission.
package valuefieldset

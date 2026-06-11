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

// Package valuemerge applies selected descriptor-semantic payload fields from
// one concrete value into another.
//
// valuemerge is a pure value-level merge primitive. valuevalidation validates
// values, valuefieldset extracts ownership fields, valuecompare computes
// semantic changed fields, valueapply decides selected merge fields and
// deletion policy, and fieldownership stores ownership state. valuemerge applies
// the selected field set only.
//
// Selection is expressed as an api/fieldpath.Set. Selecting the current path
// replaces the current subtree with the overlay subtree. Selecting a descendant
// recursively merges only that descendant. Unselected paths preserve base data.
// Selecting a path absent from overlay removes the corresponding base field
// where the parent container supports removal.
//
// Explicit NullValue is concrete payload data, not a deletion marker.
// Nullability is validated by valuevalidation, not valuemerge.
//
// Root removal is unsupported because a removed root cannot be represented as a
// value.Value. Record, map, and ListMap item removal is supported because those
// containers preserve remaining field identities. Ordered-list item removal is
// restricted to a selected tail-contiguous suffix; deleting from the middle is
// rejected because physical indexes are semantic. ListAtomic and ListSet allow
// exact parent replacement only and reject descendant merges.
//
// UnknownReject rejects undeclared record members. UnknownPreserveOpaque treats
// undeclared members as opaque leaves. UnknownPrune omits undeclared members
// during record merge. valuemerge does not perform generic pruning or
// normalization; UnknownPrune is descriptor-directed record merge behavior.
//
// ListAtomic treats the complete list as one field. ListSet is also merged only
// as one field until stable value-based set item identity exists. ListOrdered
// uses index-based semantics. ListMap uses selector identity from declared key
// fields; physical order is not semantic identity.
//
// Callers are expected to validate base and overlay values with
// valuevalidation before merge when full descriptor conformance is required.
// Descriptors are expected to have been validated upstream. valuemerge performs
// only local defensive checks needed for selected traversal and replacement
// shape. Exact subtree replacement checks zero values, DescriptorRef resolution,
// and descriptor/value kind compatibility; it does not fully validate the
// replacement subtree.
//
// Malformed payload-derived record member names that cannot become semantic
// field or map-key path elements are reported as invalid paths. The merge layer
// is constructing selected semantic paths; full payload validation remains in
// api/valuevalidation.
//
// The package does not own full validation, field extraction, semantic
// comparison, ownership mutation, conflict detection, apply or force policy,
// object metadata, storage, codec behavior, defaulting, or generic
// normalization.
package valuemerge

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

// Package valuecompare computes descriptor-semantic changes between concrete
// API payload values.
//
// valuevalidation validates values. valuefieldset extracts ownership fields.
// valuecompare computes semantic changed fields. valuemerge merges values.
// valueapply orchestrates apply, conflict policy, force policy, and ownership
// updates. fieldownership stores ownership state and detects ownership conflicts.
//
// Result values are expressed as Added, Removed, and Modified fieldpath sets.
// Added contains semantic paths present in the new value but absent in the old
// value. Removed contains semantic paths present in the old value but absent in
// the new value. Modified contains semantic paths present in both values whose
// descriptor-semantic payload changed. Changed is the union of Added, Removed,
// and Modified.
//
// Added and removed subtrees are expanded through valuefieldset so compare
// output stays aligned with ownership and apply field semantics. Result fields
// are not applied fields, ownership fields, merge fields, validation diagnostic
// paths, or representation-level diffs.
//
// Absent and explicit null are different. Absent-to-present produces Added.
// Present-to-absent produces Removed. Null-to-null produces no change. Null to
// non-null, and non-null to null, produce Modified at the current path.
//
// Decimal comparison is numeric, bytes compare by content, temporal values use
// their concrete equality methods, and float values compare exactly. Unknown
// preserved members are compared as opaque leaves. Unknown rejected members
// fail comparison. Unknown pruned members are ignored. Record member order is
// ignored for object and map descriptors. Ordered list item order is semantic.
// ListMap selector identity is semantic. ListSet comparison is intentionally
// conservative and currently treats the whole list as one order-sensitive field
// until a stable set-element identity model exists.
//
// Descriptors are expected to have been validated at construction, registration,
// or catalog boundaries. valuecompare does not call full descriptor validation on
// every comparison. It performs only local defensive checks needed for
// traversal.
//
// valuecompare does not validate full values, merge values, apply values, mutate
// ownership, detect ownership conflicts, authorize requests, or touch storage,
// codecs, defaulting, pruning, or normalization.
package valuecompare

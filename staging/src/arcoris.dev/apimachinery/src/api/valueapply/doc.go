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

// Package valueapply orchestrates value-level apply and field ownership policy.
//
// valueapply coordinates lower-level value packages without taking over their
// responsibilities. valuevalidation validates concrete values against prepared
// descriptors. valuefieldset extracts AppliedFields from the applied payload.
// valuecompare computes descriptor-semantic Changes between live and applied
// values. valuemerge applies selected MergeFields. fieldownership stores owner
// state and detects structural ownership conflicts. valueapply decides how
// those pieces are sequenced and which ownership policy is applied.
//
// The pipeline order is fixed:
//
//  1. validate request shape and ownership scope;
//  2. validate Live and Applied values;
//  3. extract AppliedFields;
//  4. compare Live and Applied values;
//  5. compute ChangedAppliedFields;
//  6. detect ownership conflicts;
//  7. enforce force/takeover policy;
//  8. plan DroppedFields and DeletedFields;
//  9. compute MergeFields;
//  10. merge selected fields;
//  11. replace ownership state;
//  12. assemble the public Result.
//
// AppliedFields are ownership fields explicitly mentioned by Applied under
// descriptor semantics. They are not necessarily changed fields, conflict
// fields, merge fields, or all fields in the merged value.
//
// ChangedAppliedFields are AppliedFields that structurally overlap
// Changes.Changed(). This set is the conflict-attempt set. Only changed applied
// fields conflict; applying the same value to a field owned by another owner is
// allowed and can create shared ownership.
//
// DroppedFields are fields previously owned by Owner but no longer covered by
// AppliedFields. Dropping releases ownership. DeletedFields are dropped fields
// that no other owner protects with an exact, ancestor, or descendant overlap.
// Dropping ownership is therefore distinct from deleting value data.
//
// MergeFields are AppliedFields union DeletedFields. valueapply delegates the
// actual selected-field merge to valuemerge. Unsupported merge shapes, such as
// unsafe ordered-list index removals, fail before ownership is updated.
//
// Force bypasses conflict failure only for representable conflicting changed
// applied fields. It removes overlapping ownership from other owners only for
// conflicting attempted paths. It does not remove unrelated ownership, and it
// does not remove ownership that protects dropped fields. If a conflict would
// require subtracting a child path from another owner's ancestor ownership,
// Force is rejected as an unsupported takeover.
//
// Apply is side-effect free. It does not mutate request values or the input
// ownership state, and it updates ownership only after merge succeeds.
//
// The package does not handle object metadata, storage, admission,
// authorization, codecs, defaulting, pruning, normalization, event emission,
// resource catalog access, or runtime lifecycle behavior.
package valueapply

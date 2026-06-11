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

// Package valueapply applies one concrete API payload value to another under
// api/types descriptor semantics.
//
// The package is a pure value-level orchestration layer. It validates live and
// applied values, extracts applied field sets, compares live and applied values,
// checks field ownership conflicts, merges selected fields, and returns updated
// field ownership state.
//
// The pipeline is deliberately staged: validate the request, validate both
// values, extract applied ownership fields, compare for semantic changes, check
// ownership conflicts, plan dropped/deleted fields, merge selected fields, and
// write a replacement ownership state.
//
// Conflicts are checked only for applied fields whose values would actually
// change. Applying the same value to a field owned by another owner is not a
// conflict and can create shared ownership.
//
// Fields previously owned by the applying owner but omitted from the new
// applied value are released. Such dropped fields are deleted from the live
// value only if no other owner structurally overlaps them.
//
// Force only takes conflicting fields from other owners. It does not remove
// unrelated ownership or ownership for dropped fields.
//
// valueapply requires both live and applied values to be valid for the
// descriptor before apply proceeds. It is not a repair, pruning, defaulting, or
// normalization operation. Invalid live data is rejected even if the applied
// value would otherwise remove the invalid field.
//
// The package does not read or write storage, run admission, authorize request
// subjects, manage object metadata, serialize managed fields, decode wire
// formats, emit events, access resource catalogs, or execute runtime lifecycle
// behavior. Runtime and object-level layers are responsible for those concerns.
package valueapply

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

// Package fieldownership models ownership state over semantic API field paths.
//
// The package is a foundation layer over api/fieldpath.Set. It stores
// owner-to-field-set entries, normalizes deterministic ownership state, detects
// structural ownership conflicts for attempted field sets, and provides state
// transformation helpers for higher layers.
//
// Owner identities are local field-ownership identities. They are not admission
// roles, authorization subjects, RBAC principals, runtime component instances,
// request actors, storage keys, audit identities, or policy identities. Higher
// layers decide which request subject may act as which Owner.
//
// State is normalized by owner order. Duplicate owners are merged by exact
// field-set union and empty entries are pruned. State deliberately preserves
// parent/child path pairs and overlapping ownership across owners. It is a
// representation model; conflicts are contextual.
//
// A conflict means structural overlap between attempted paths and paths owned by
// a different owner. Exact, ancestor, and descendant relations conflict. Sibling
// paths do not conflict, and the same owner never conflicts with itself.
//
// Transform helpers separate exact and overlap semantics. RemoveFields and
// RemoveFieldsFromOthers remove exact paths only. RemoveOverlappingFields and
// RemoveOverlappingFieldsFromOthers remove exact, ancestor, and descendant
// paths that overlap the supplied field set. SetFields replaces one owner's
// fields exactly; an empty replacement removes that owner.
//
// The package does not decide what fields are attempted, whether conflicts are
// fatal, whether force takeover is allowed, whether omitted fields release
// ownership, or how ownership is serialized or stored. It does not inspect
// api/value payloads, traverse api/types descriptors, compare values, extract
// field sets, merge values, apply configurations, serialize managed fields,
// validate object metadata, access storage, perform admission, or authorize
// owners.
package fieldownership

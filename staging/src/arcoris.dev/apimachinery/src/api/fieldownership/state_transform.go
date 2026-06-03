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

package fieldownership

import "arcoris.dev/apimachinery/api/fieldpath"

// WithFields replaces owner fields with fields exactly.
//
// Empty fields remove owner from the state. Other owners are preserved unchanged.
func (s State) WithFields(owner Owner, fields fieldpath.Set) (State, error) {
	if err := validateOwnerFields(
		owner,
		fields,
		"replacement field path is invalid",
	); err != nil {
		return State{}, err
	}

	return s.replaceOwnerFields(owner, fields)
}

// AddFields adds fields to owner without dropping existing owner fields.
func (s State) AddFields(owner Owner, fields fieldpath.Set) (State, error) {
	return s.transformOwnerFields(
		owner,
		fields,
		"added field path is invalid",
		unionTransform,
	)
}

// RemoveFields removes exact fields from owner only.
//
// Ancestor and descendant paths are left intact. Use RemoveOverlaps when the
// caller has explicitly chosen structural-overlap removal semantics.
func (s State) RemoveFields(owner Owner, fields fieldpath.Set) (State, error) {
	return s.transformOwnerFields(
		owner,
		fields,
		"removed field path is invalid",
		removeExactTransform,
	)
}

// RemoveOverlaps removes owner paths that structurally overlap fields.
//
// Exact matches, owned ancestors, and owned descendants of the supplied fields
// are removed from owner only.
func (s State) RemoveOverlaps(owner Owner, fields fieldpath.Set) (State, error) {
	return s.transformOwnerFields(
		owner,
		fields,
		"overlap field path is invalid",
		removeOverlapTransform,
	)
}

// RemoveFieldsFromOthers removes exact fields from every owner except owner.
//
// This helper supports higher-layer force/takeover policies without deciding
// when those policies should run.
func (s State) RemoveFieldsFromOthers(owner Owner, fields fieldpath.Set) (State, error) {
	return s.transformOtherOwnerFields(
		owner,
		fields,
		"other-owner field path is invalid",
		removeExactTransform,
	)
}

// RemoveOverlapsFromOthers removes overlapping paths from every owner except owner.
//
// This is the structural-overlap counterpart to RemoveFieldsFromOthers.
func (s State) RemoveOverlapsFromOthers(owner Owner, fields fieldpath.Set) (State, error) {
	return s.transformOtherOwnerFields(
		owner,
		fields,
		"other-owner overlap field path is invalid",
		removeOverlapTransform,
	)
}

// WithoutOwner removes owner entirely from s.
func (s State) WithoutOwner(owner Owner) (State, error) {
	if err := owner.Validate(); err != nil {
		return State{}, err
	}

	return s.replaceOwnerFields(owner, fieldpath.EmptySet())
}

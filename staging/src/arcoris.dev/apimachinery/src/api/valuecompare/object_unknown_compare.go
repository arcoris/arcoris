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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// compareUnknownObjectMembers applies object unknown-field policy.
func (c *comparer) compareUnknownObjectMembers(
	path fieldpath.Path,
	oldObject value.ObjectView,
	newObject value.ObjectView,
	declared map[string]types.FieldDescriptor,
	policy types.UnknownFieldPolicy,
) (Result, error) {
	switch policy {
	case types.UnknownReject:
		return c.rejectUnknownObjectMembers(path, oldObject, newObject, declared)
	case types.UnknownPreserve:
		return c.comparePreservedUnknownObjectMembers(path, oldObject, newObject, declared)
	case types.UnknownPrune:
		return EmptyResult(), nil
	default:
		return Result{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"object unknown-field policy is invalid",
		)
	}
}

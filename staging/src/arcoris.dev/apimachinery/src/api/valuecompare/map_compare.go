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
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// compareMap interprets value.KindObject as a dynamic string-keyed map.
//
// TypeMap shares the concrete object payload with TypeObject, but uses
// path.Key(key) rather than path.Field(name). Dynamic keys are sorted before
// traversal so result construction is deterministic.
func (c *comparer) compareMap(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
	depth int,
) (Result, error) {
	if err := requireKind(path, oldValue, value.KindObject, descriptor.Code()); err != nil {
		return Result{}, err
	}
	if err := requireKind(path, newValue, value.KindObject, descriptor.Code()); err != nil {
		return Result{}, err
	}

	mapView, ok := descriptor.Map()
	if !ok {
		return Result{}, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a map")
	}
	if !mapView.Key().IsValid() {
		return Result{}, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "map key type is invalid")
	}

	valueType := mapView.Value()
	if !valueType.IsValid() {
		return Result{}, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "map value descriptor is invalid")
	}

	oldObject, _ := oldValue.Object()
	newObject, _ := newValue.Object()
	oldMembers := membersByName(oldObject.Members())
	newMembers := membersByName(newObject.Members())

	result := EmptyResult()
	for _, key := range unionSortedKeys(oldMembers, newMembers) {
		memberPath := path.Key(key)
		child, err := c.compare(
			memberPath,
			mapOperand(oldMembers, key),
			mapOperand(newMembers, key),
			valueType,
			depth+1,
		)
		if err != nil {
			return Result{}, err
		}

		result = result.merge(child)
	}

	return result, nil
}

// mapOperand converts a dynamic map lookup into presence-aware compare input.
func mapOperand(members map[string]value.Value, key string) valuepresence.Operand {
	val, ok := members[key]
	return valuepresence.From(val, ok)
}

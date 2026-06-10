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

// equalMap reports descriptor-aware map equality without building result sets.
//
// Map entries use key-path semantics even when equality is called from a
// whole-value decision such as atomic list comparison.
func (c *comparer) equalMap(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Descriptor,
	depth int,
) (bool, error) {
	if err := requireKind(path, oldValue, value.KindRecord, descriptor.Code()); err != nil {
		return false, err
	}
	if err := requireKind(path, newValue, value.KindRecord, descriptor.Code()); err != nil {
		return false, err
	}

	mapView, ok := descriptor.AsMap()
	if !ok {
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a map")
	}
	if !mapView.Key().IsValid() {
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "map key descriptor is invalid")
	}

	valueDescriptor := mapView.Value()
	if !valueDescriptor.IsValid() {
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "map value descriptor is invalid")
	}

	oldObject, _ := oldValue.AsRecord()
	newObject, _ := newValue.AsRecord()
	if oldObject.Len() != newObject.Len() {
		return false, nil
	}

	oldMembers := membersByName(oldObject.Members())
	for _, newMember := range newObject.Members() {
		name := newMember.Name.String()
		oldMember, found := oldMembers[name]
		if !found {
			return false, nil
		}

		equal, err := c.equalValue(path.Key(name), oldMember, newMember.Value, valueDescriptor, depth+1)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}

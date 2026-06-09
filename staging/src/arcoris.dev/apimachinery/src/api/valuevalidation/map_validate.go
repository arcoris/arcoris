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

package valuevalidation

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// validateMap interprets value.KindObject as a dynamic string-keyed map descriptor.
func (v *validator) validateMap(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	depth int,
) {
	if !v.requireKind(path, val, value.KindObject, descriptor.Code()) {
		return
	}

	mapView, ok := descriptor.AsMap()
	if !ok {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a map")
		return
	}

	keyDescriptor := mapView.Key()
	if !keyDescriptor.IsValid() {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "map key descriptor is invalid")
		return
	}

	valueDescriptor := mapView.Value()
	if !valueDescriptor.IsValid() {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "map value descriptor is invalid")
		return
	}

	valueView, _ := val.Object()
	length := valueView.Len()
	if minEntries, ok := mapView.MinEntries(); ok && length < minEntries {
		v.addf(
			path,
			ErrLengthOutOfRange,
			ErrorReasonTooShort,
			"map entry count %d is below minimum %d",
			length,
			minEntries,
		)
	}
	if maxEntries, ok := mapView.MaxEntries(); ok && length > maxEntries {
		v.addf(
			path,
			ErrLengthOutOfRange,
			ErrorReasonTooLong,
			"map entry count %d is above maximum %d",
			length,
			maxEntries,
		)
	}

	for _, mapMember := range valueView.Members() {
		memberPath := path.Key(mapMember.Name)
		v.validate(memberPath, value.StringValue(mapMember.Name), keyDescriptor, depth+1)
		v.validate(memberPath, mapMember.Value, valueDescriptor, depth+1)
	}
}

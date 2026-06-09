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

package valuefieldset

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// extractMap interprets value.KindObject as a dynamic string-keyed map descriptor.
func (e *extractor) extractMap(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	depth int,
) (fieldpath.Set, error) {
	if err := requireKind(path, val, value.KindObject, descriptor.Code()); err != nil {
		return fieldpath.Set{}, err
	}

	mapView, ok := descriptor.AsMap()
	if !ok {
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not a map",
		)
	}
	if !mapView.Key().IsValid() {
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"map key descriptor is invalid",
		)
	}

	valueDescriptor := mapView.Value()
	if !valueDescriptor.IsValid() {
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"map value descriptor is invalid",
		)
	}

	valueView, _ := val.Object()
	if valueView.IsEmpty() {
		return setAt(path)
	}

	out := fieldpath.EmptySet()
	for _, mapMember := range valueView.Members() {
		memberSet, err := e.extract(
			path.Key(mapMember.Name),
			mapMember.Value,
			valueDescriptor,
			depth+1,
		)
		if err != nil {
			return fieldpath.Set{}, err
		}

		out = out.Union(memberSet)
	}

	return out, nil
}

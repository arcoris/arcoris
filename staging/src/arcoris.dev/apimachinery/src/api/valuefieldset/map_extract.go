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

// extractMap interprets value.KindRecord as a dynamic string-keyed map descriptor.
func (e *extractor) extractMap(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	depth int,
) (fieldpath.Set, error) {
	if err := requireKind(path, val, value.KindRecord, descriptor.Code()); err != nil {
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

	valueView, _ := val.AsRecord()
	if valueView.IsEmpty() {
		return setAt(path)
	}

	var out setBuilder
	var extractErr error
	valueView.ForEach(func(_ int, mapMember value.RecordMember) bool {
		memberPath, err := mapMemberPath(path, mapMember.Name.String())
		if err != nil {
			extractErr = err
			return false
		}

		memberSet, err := e.extract(
			memberPath,
			mapMember.Value,
			valueDescriptor,
			depth+1,
		)
		if err != nil {
			extractErr = err
			return false
		}

		out.AddSet(memberSet)
		return true
	})
	if extractErr != nil {
		return fieldpath.Set{}, extractErr
	}

	return out.Build(path)
}

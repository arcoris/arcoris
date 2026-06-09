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
	"arcoris.dev/apimachinery/api/internal/typekind"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// extractScalar records the current path for a concrete leaf payload.
//
// Field-set extraction intentionally ignores scalar constraints such as
// min/max, pattern, enum, precision, and scale. Those are validation rules, not
// path-discovery rules.
func (e *extractor) extractScalar(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
) (fieldpath.Set, error) {
	expected, ok := scalarKind(descriptor.Code())
	if !ok {
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not a scalar type",
		)
	}

	if err := requireKind(path, val, expected, descriptor.Code()); err != nil {
		return fieldpath.Set{}, err
	}

	return setAt(path)
}

// scalarKind maps descriptor scalar codes to concrete value kinds.
func scalarKind(code types.DescriptorKind) (value.Kind, bool) {
	return typekind.Scalar(code)
}

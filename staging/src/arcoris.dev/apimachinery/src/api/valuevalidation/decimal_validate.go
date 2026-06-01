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

// validateDecimal checks decimal kind, precision, and scale constraints.
//
// api/types does not currently expose decimal bounds or enum rules, and
// api/value.Decimal intentionally has no ordering API yet. The validator
// therefore keeps this first pass to exact descriptor rules that are already
// defined: maximum precision and maximum scale.
func (v *validator) validateDecimal(path fieldpath.Path, val value.Value, descriptor types.Type) {
	if !v.requireKind(path, val, value.KindDecimal, descriptor.Code()) {
		return
	}

	got, _ := val.Decimal()
	view, ok := descriptor.Decimal()
	if !ok {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not decimal")
		return
	}

	precision := len(got.Coefficient())
	if maxPrecision, ok := view.Precision(); ok && precision > maxPrecision {
		v.addf(
			path,
			ErrValueOutOfRange,
			ErrorReasonAboveMaximum,
			"decimal precision %d is above maximum %d",
			precision,
			maxPrecision,
		)
	}
	if scale, ok := view.Scale(); ok && int(got.Scale()) > scale {
		v.addf(
			path,
			ErrValueOutOfRange,
			ErrorReasonAboveMaximum,
			"decimal scale %d is above maximum %d",
			got.Scale(),
			scale,
		)
	}
}

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

package valuemerge

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// Merge copies selected semantic fields from overlay into base at root.
func Merge(
	base value.Value,
	overlay value.Value,
	descriptor types.Descriptor,
	fields fieldpath.Set,
	opts Options,
) (value.Value, error) {
	return MergeAt(fieldpath.Root(), base, overlay, descriptor, fields, opts)
}

// MergeAt copies selected semantic fields from overlay into base below path.
//
// Empty selected fields are a no-op: MergeAt validates path, fields, and base,
// then returns a clone of base without inspecting overlay. For non-empty
// selections, callers should validate base and overlay with api/valuevalidation
// when full descriptor conformance is required. valuemerge performs only the
// defensive checks needed for selected traversal and replacement shape.
func MergeAt(
	path fieldpath.Path,
	base value.Value,
	overlay value.Value,
	descriptor types.Descriptor,
	fields fieldpath.Set,
	opts Options,
) (value.Value, error) {
	if err := path.ValidateStructure(); err != nil {
		return value.Value{}, wrapAt(
			path,
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"base field path is invalid",
			err,
		)
	}
	if err := validateFieldsAt(path, fields); err != nil {
		return value.Value{}, err
	}
	if err := requireValidValue(path, valuepresence.Present(base)); err != nil {
		return value.Value{}, err
	}
	if fields.IsEmpty() {
		return base.Clone(), nil
	}
	if err := requireValidValue(path, valuepresence.Present(overlay)); err != nil {
		return value.Value{}, err
	}

	result, err := newMerger(opts).merge(
		path,
		valuepresence.Present(base),
		valuepresence.Present(overlay),
		descriptor,
		fields,
		0,
	)

	if err != nil {
		return value.Value{}, err
	}
	if result.Absent() {
		return value.Value{}, errorAt(
			path,
			ErrUnsupportedMerge,
			ErrorReasonUnsupportedMerge,
			"root removal is not representable",
		)
	}

	return result.Value(), nil
}

// validateFieldsAt rejects malformed paths and selections outside base.
func validateFieldsAt(base fieldpath.Path, fields fieldpath.Set) error {
	var fieldErr error
	fields.ForEach(func(_ int, path fieldpath.Path) bool {
		if err := path.ValidateStructure(); err != nil {
			fieldErr = wrapAt(
				path,
				ErrInvalidPath,
				ErrorReasonInvalidPath,
				"selected field path is invalid",
				err,
			)
			return false
		}
		if !path.Equal(base) && !path.IsDescendantOf(base) {
			fieldErr = errorfAt(
				path,
				ErrInvalidPath,
				ErrorReasonInvalidPath,
				"selected field path %s is outside merge base path %s",
				path.CanonicalText(),
				base.CanonicalText(),
			)
			return false
		}

		return true
	})

	return fieldErr
}

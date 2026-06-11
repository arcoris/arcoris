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

import "arcoris.dev/apimachinery/api/fieldpath"

// recordFieldPath converts a concrete record member name into a semantic field path.
func recordFieldPath(base fieldpath.Path, name string) (fieldpath.Path, error) {
	fieldName, err := fieldpath.NewFieldName(name)
	if err != nil {
		return fieldpath.Path{}, wrapAt(
			base,
			ErrInvalidPath,
			ErrorReasonInvalidFieldName,
			"record member name cannot become a field path element",
			err,
		)
	}

	return base.Field(fieldName), nil
}

// mapKeyPath converts a concrete record member name into a semantic map-key path.
func mapKeyPath(base fieldpath.Path, key string) (fieldpath.Path, error) {
	mapKey, err := fieldpath.NewMapKey(key)
	if err != nil {
		return fieldpath.Path{}, wrapAt(
			base,
			ErrInvalidPath,
			ErrorReasonInvalidMapKey,
			"map member name cannot become a map-key path element",
			err,
		)
	}

	return base.Key(mapKey), nil
}

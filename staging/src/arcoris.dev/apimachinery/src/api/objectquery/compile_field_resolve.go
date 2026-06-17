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

package objectquery

import "fmt"

// resolveAndValidateField resolves field metadata and validates operator and
// literal compatibility before a field term can be evaluated.
func resolveAndValidateField(t term, opts compileOptions) (SelectableField, error) {
	if err := t.fieldRef.Validate(); err != nil {
		return SelectableField{}, invalidFieldError(err, "invalid field ref")
	}
	field, err := resolveField(t.fieldRef, opts)
	if err != nil {
		return SelectableField{}, err
	}
	if err := field.Validate(); err != nil {
		return SelectableField{}, invalidFieldError(err, "invalid selectable field")
	}
	if !field.Operators.Supports(t.operator) {
		return SelectableField{}, unsupportedOperatorError(t.operator, "field "+t.fieldRef.String())
	}
	if err := validateFieldLiterals(field, t.operator, t.values); err != nil {
		return SelectableField{}, err
	}

	return field, nil
}

// resolveField is the narrow registry boundary for selectable fields. Missing
// registries and unknown fields are both query validation failures.
func resolveField(ref FieldRef, opts compileOptions) (SelectableField, error) {
	if opts.fields == nil {
		return SelectableField{}, unresolvedFieldError(
			ref,
			fmt.Errorf("no selectable field set supplied for %s", ref.String()),
		)
	}

	field, ok := opts.fields.ResolveSelectableField(ref)
	if !ok {
		return SelectableField{}, unresolvedFieldError(
			ref,
			fmt.Errorf("selectable field %s is not registered", ref.String()),
		)
	}

	return field, nil
}

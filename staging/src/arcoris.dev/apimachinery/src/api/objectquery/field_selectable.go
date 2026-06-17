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

import (
	"fmt"

	"arcoris.dev/apimachinery/api/value"
)

// SelectableField describes one field that objectquery may evaluate.
//
// Resource or adapter code owns the registry of selectable fields. objectquery
// only enforces this declaration and evaluates registered refs.
type SelectableField struct {
	// Ref is the stable field identity.
	Ref FieldRef
	// Kind is the non-null payload kind expected for the field.
	Kind value.Kind
	// Operators is the finite operator set allowed for this field.
	Operators OperatorSet
	// Index is an advisory indexing hint for future caches/storage adapters.
	Index IndexHint
}

// Validate checks whether f is usable for field term compilation.
func (f SelectableField) Validate() error {
	if err := f.Ref.Validate(); err != nil {
		return err
	}
	if !f.Kind.IsValid() {
		return fmt.Errorf("%w: invalid field kind", ErrInvalidField)
	}
	if f.Operators == 0 {
		return fmt.Errorf("%w: empty operator set", ErrInvalidField)
	}

	return nil
}

// SelectableFieldSet resolves registered selectable fields by exact FieldRef.
type SelectableFieldSet interface {
	ResolveSelectableField(FieldRef) (SelectableField, bool)
}

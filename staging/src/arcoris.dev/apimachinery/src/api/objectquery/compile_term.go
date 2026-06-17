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

import "arcoris.dev/apimachinery/api/objectstore"

// validateTerm checks the active payload for a leaf term and returns a detached
// canonical term. It intentionally does not inspect object payloads.
func validateTerm(t term, opts compileOptions) (term, error) {
	switch t.kind {
	case termResource:
		if err := t.resource.Validate(); err != nil {
			return term{}, invalidTermError("invalid resource: %w", err)
		}
	case termNamespace:
		if err := t.namespace.ValidateLexical(); err != nil {
			return term{}, invalidTermError("invalid namespace: %w", err)
		}
	case termName:
		if err := t.name.ValidateLexical(); err != nil {
			return term{}, invalidTermError("invalid name: %w", err)
		}
	case termObject:
		if _, err := ObjectEquals(t.namespace, t.name); err != nil {
			return term{}, err
		}
	case termKey:
		if err := objectstore.ValidateKey(t.key); err != nil {
			return term{}, invalidTermError("invalid key: %w", err)
		}
	case termMetadata:
		values, err := validateMetadataTerm(
			t.metadataDomain,
			t.operator,
			t.metadataKey,
			t.stringValues,
		)
		if err != nil {
			return term{}, err
		}
		t.stringValues = values
	case termField:
		field, err := resolveAndValidateField(t, opts)
		if err != nil {
			return term{}, err
		}
		t.fieldRef = field.Ref
		t.field = field
	default:
		return term{}, invalidTermError("unknown term kind")
	}

	return t.clone(), nil
}

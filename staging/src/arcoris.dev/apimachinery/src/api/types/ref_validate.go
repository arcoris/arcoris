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

package types

import "errors"

// validateRef checks TypeRef syntax, optional resolver lookup, and cycles.
func validateRef(t Type, resolver Resolver, path string, resolving map[TypeName]bool) error {
	name := t.ref.name

	if !name.IsValid() {
		return typeErrorf(
			path,
			ErrInvalidTypeReference,
			TypeErrorReasonInvalidReferenceName,
			"reference name %q is not a valid TypeName",
			name,
		)
	}

	if resolver == nil {
		return nil
	}

	if resolving[name] {
		return typeErrorf(
			path,
			ErrInvalidTypeReference,
			TypeErrorReasonReferenceCycle,
			"reference %q creates a recursive TypeDefinition graph",
			name,
		)
	}

	def, ok := resolver.ResolveType(name)

	if !ok {
		return typeErrorf(
			path,
			ErrUnknownTypeReference,
			TypeErrorReasonUnknownReference,
			"reference %q was not found in resolver",
			name,
		)
	}

	next := copyResolving(resolving)
	next[name] = true

	if err := validateType(def.Type(), resolver, path, next); err != nil {
		var typeErr *TypeError

		if errors.As(err, &typeErr) && typeErr.Reason == TypeErrorReasonReferenceCycle {
			return err
		}

		return typeErrorf(
			path,
			err,
			TypeErrorReasonInvalidResolvedDefinition,
			"resolved reference %q is structurally invalid: %v",
			name,
			err,
		)
	}

	return nil
}

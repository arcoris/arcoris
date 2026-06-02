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

// validateListMapKeyIdentityType checks the descriptor-level identity contract
// for one ListMap key field.
//
// ListMap keys are stronger than ordinary required object fields. Future value
// validation, diff, and apply layers must be able to extract them into stable
// selector identity values. This first pass limits that identity surface to
// non-nullable bool, string, and integer descriptors, including references that
// resolve to those descriptors.
func validateListMapKeyIdentityType(
	field FieldDescriptor,
	resolver Resolver,
	path string,
	resolving map[TypeName]bool,
) error {
	return validateListMapKeyType(field.Type(), resolver, path, resolving)
}

// validateListMapKeyType validates a concrete or referenced key type.
func validateListMapKeyType(
	typ Type,
	resolver Resolver,
	path string,
	resolving map[TypeName]bool,
) error {
	if typ.Nullable() {
		return typeErrorf(
			path,
			ErrInvalidField,
			TypeErrorReasonInvalidListMapKeyType,
			"ListMap key field must be non-nullable",
		)
	}

	switch typ.code {
	case TypeBool,
		TypeString,
		TypeInt8,
		TypeInt16,
		TypeInt32,
		TypeInt64,
		TypeUint8,
		TypeUint16,
		TypeUint32,
		TypeUint64:
		return nil
	case TypeRef:
		return validateListMapKeyRef(typ.ref.name, resolver, path, resolving)
	default:
		return typeErrorf(
			path,
			ErrInvalidField,
			TypeErrorReasonInvalidListMapKeyType,
			"ListMap key field type %s cannot be represented as a stable identity value",
			typ.code,
		)
	}
}

// validateListMapKeyRef resolves a referenced key type without introducing a
// catalog dependency or a public identity-type API.
func validateListMapKeyRef(
	name TypeName,
	resolver Resolver,
	path string,
	resolving map[TypeName]bool,
) error {
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
		return typeErrorf(
			path,
			ErrUnknownTypeReference,
			TypeErrorReasonUnknownReference,
			"reference %q cannot be resolved without a resolver",
			name,
		)
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

	return validateListMapKeyType(def.Type(), resolver, path, next)
}

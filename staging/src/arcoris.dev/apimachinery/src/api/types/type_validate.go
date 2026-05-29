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

// ValidateType checks the structural integrity of t.
//
// Validation is descriptor validation only. It checks type-code/payload
// consistency, field descriptors, list-map key shape, references, and resolver
// lookup when a resolver is supplied. It does not validate concrete object
// values, apply defaults, prune fields, export schemas, or run arbitrary Go
// validators.
func ValidateType(t Type, resolver Resolver) error {
	return validateType(t, resolver, "type", make(map[TypeName]bool))
}

// ValidateDefinition checks the structural integrity of def.
//
// A provided resolver is used to resolve TypeRef descriptors. The definition
// name itself is treated as part of the active reference stack so direct and
// indirect reference cycles are rejected. Recursive TypeDefinition graphs are
// not supported by api/types; recursive schemas require a future explicit
// design pass.
func ValidateDefinition(def TypeDefinition, resolver Resolver) error {
	if !def.name.IsValid() {
		return typeErrorf(
			"definition.name",
			ErrInvalidTypeReference,
			TypeErrorReasonInvalidReferenceName,
			"definition name %q is not a valid TypeName",
			def.name,
		)
	}
	resolving := map[TypeName]bool{def.name: true}
	if err := validateType(def.typ, resolver, "definition.type", resolving); err != nil {
		return err
	}
	return nil
}

// validateType validates t at a descriptor path.
//
// The path parameter describes descriptor structure, not a future object-value
// path. It is threaded through recursive validation so callers receive precise
// errors such as type.fields[spec].type or ref(example.Name).
func validateType(t Type, resolver Resolver, path string, resolving map[TypeName]bool) error {
	if !t.code.IsValid() {
		return typeErrorf(
			path,
			errors.Join(ErrInvalidType, ErrInvalidTypeCode),
			TypeErrorReasonInvalidTypeCode,
			"type code %d is not supported",
			t.code,
		)
	}
	if err := validateInactivePayloads(t, path); err != nil {
		return err
	}
	if t.code == TypeNull && t.Nullable() {
		return typeErrorf(
			path,
			ErrInvalidType,
			TypeErrorReasonInvalidNullability,
			"TypeNull is the null literal and cannot be nullable",
		)
	}
	switch t.code {
	case TypeNull:
		return nil
	case TypeBool:
		return nil
	case TypeString:
		return validateString(t, path)
	case TypeBytes:
		return validateBytes(t, path)
	case TypeInt8:
		return validateInt8(t, path)
	case TypeInt16:
		return validateInt16(t, path)
	case TypeInt32:
		return validateInt32(t, path)
	case TypeInt64:
		return validateInt64(t, path)
	case TypeUint8:
		return validateUint8(t, path)
	case TypeUint16:
		return validateUint16(t, path)
	case TypeUint32:
		return validateUint32(t, path)
	case TypeUint64:
		return validateUint64(t, path)
	case TypeFloat32:
		return validateFloat32(t, path)
	case TypeFloat64:
		return validateFloat64(t, path)
	case TypeDecimal:
		return validateDecimal(t, path)
	case TypeTimestamp:
		return validateTimestamp(t, path)
	case TypeDate:
		return validateDate(t, path)
	case TypeTime:
		return validateTime(t, path)
	case TypeDuration:
		return validateDuration(t, path)
	case TypeObject:
		return validateObject(t, resolver, path, resolving)
	case TypeList:
		return validateList(t, resolver, path, resolving)
	case TypeMap:
		return validateMap(t, resolver, path, resolving)
	case TypeRef:
		return validateRef(t, resolver, path, resolving)
	default:
		return typeErrorf(
			path,
			ErrInvalidType,
			TypeErrorReasonInvalidTypeCode,
			"type code %d has no validator",
			t.code,
		)
	}
}

// copyResolving detaches the active reference stack.
//
// Reference validation is recursive. Each branch receives its own resolving map
// so sibling references cannot accidentally affect one another while cycle
// detection still catches the active chain.
func copyResolving(in map[TypeName]bool) map[TypeName]bool {
	out := make(map[TypeName]bool, len(in)+1)
	for name, resolving := range in {
		out[name] = resolving
	}
	return out
}

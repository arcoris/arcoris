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

// ValidateLocal checks the local structural integrity of d.
//
// Local validation validates descriptor shape and reference-name syntax, but it
// does not require referenced definitions to exist. Use ValidateResolved when a
// concrete Resolver is available and unresolved references should fail.
func ValidateLocal(d Descriptor) error {
	return validateDescriptor(d, nil, "descriptor", make(map[TypeName]bool))
}

// ValidateResolved checks the structural integrity of d against resolver.
//
// Resolved validation validates local descriptor shape, resolves every Ref, and
// rejects unresolved, invalid, or recursive definition graphs. Concrete values,
// defaults, pruning, schema export, codecs, and arbitrary Go validators remain
// outside api/types.
func ValidateResolved(d Descriptor, resolver Resolver) error {
	if resolver == nil {
		return descriptorErrorf(
			"resolver",
			ErrInvalidDescriptorReference,
			DescriptorErrorReasonUnknownReference,
			"resolved descriptor validation requires a non-nil resolver",
		)
	}

	return validateDescriptor(d, resolver, "descriptor", make(map[TypeName]bool))
}

// ValidateDefinitionLocal checks the local structural integrity of def.
func ValidateDefinitionLocal(def Definition) error {
	if !def.name.IsValid() {
		return descriptorErrorf(
			"definition.name",
			ErrInvalidDescriptorReference,
			DescriptorErrorReasonInvalidReferenceName,
			"definition name %q is not a valid TypeName",
			def.name,
		)
	}

	resolving := map[TypeName]bool{def.name: true}
	return validateDescriptor(def.descriptor, nil, "definition.descriptor", resolving)
}

// ValidateDefinitionResolved checks the resolved structural integrity of def.
//
// A provided resolver is used to resolve DescriptorRef descriptors. The definition
// name itself is treated as part of the active reference stack so direct and
// indirect reference cycles are rejected. Recursive Definition graphs are
// not supported by api/types; recursive schemas require a future explicit
// design pass.
func ValidateDefinitionResolved(def Definition, resolver Resolver) error {
	if !def.name.IsValid() {
		return descriptorErrorf(
			"definition.name",
			ErrInvalidDescriptorReference,
			DescriptorErrorReasonInvalidReferenceName,
			"definition name %q is not a valid TypeName",
			def.name,
		)
	}
	if resolver == nil {
		return descriptorErrorf(
			"resolver",
			ErrInvalidDescriptorReference,
			DescriptorErrorReasonUnknownReference,
			"resolved definition validation requires a non-nil resolver",
		)
	}

	resolving := map[TypeName]bool{def.name: true}

	if err := validateDescriptor(def.descriptor, resolver, "definition.descriptor", resolving); err != nil {
		return err
	}

	return nil
}

// validateDescriptor validates desc at a descriptor path.
//
// The path parameter describes descriptor structure, not a future object-value
// path. It is threaded through recursive validation so callers receive precise
// errors such as descriptor.fields[spec].type or ref(example.dev.Name).
func validateDescriptor(desc Descriptor, resolver Resolver, path string, resolving map[TypeName]bool) error {
	if !desc.code.IsValid() {
		return descriptorErrorf(
			path,
			errors.Join(ErrInvalidDescriptor, ErrInvalidDescriptorKind),
			DescriptorErrorReasonInvalidDescriptorKind,
			"descriptor kind %d is not supported",
			desc.code,
		)
	}

	if err := validateInactiveDescriptorPayloads(desc, path); err != nil {
		return err
	}

	if desc.code == DescriptorNull && desc.Nullable() {
		return descriptorErrorf(
			path,
			ErrInvalidDescriptor,
			DescriptorErrorReasonInvalidNullability,
			"DescriptorNull is the null literal and cannot be nullable",
		)
	}

	switch desc.code {
	case DescriptorNull:
		return nil
	case DescriptorBool:
		return nil
	case DescriptorString:
		return validateString(desc, path)
	case DescriptorBytes:
		return validateBytes(desc, path)
	case DescriptorInt8:
		return validateInt8(desc, path)
	case DescriptorInt16:
		return validateInt16(desc, path)
	case DescriptorInt32:
		return validateInt32(desc, path)
	case DescriptorInt64:
		return validateInt64(desc, path)
	case DescriptorUint8:
		return validateUint8(desc, path)
	case DescriptorUint16:
		return validateUint16(desc, path)
	case DescriptorUint32:
		return validateUint32(desc, path)
	case DescriptorUint64:
		return validateUint64(desc, path)
	case DescriptorFloat32:
		return validateFloat32(desc, path)
	case DescriptorFloat64:
		return validateFloat64(desc, path)
	case DescriptorDecimal:
		return validateDecimal(desc, path)
	case DescriptorTimestamp:
		return validateTimestamp(desc, path)
	case DescriptorDate:
		return validateDate(desc, path)
	case DescriptorTime:
		return validateTime(desc, path)
	case DescriptorDuration:
		return validateDuration(desc, path)
	case DescriptorObject:
		return validateObject(desc, resolver, path, resolving)
	case DescriptorList:
		return validateList(desc, resolver, path, resolving)
	case DescriptorMap:
		return validateMap(desc, resolver, path, resolving)
	case DescriptorRef:
		return validateRef(desc, resolver, path, resolving)
	default:
		return descriptorErrorf(
			path,
			ErrInvalidDescriptor,
			DescriptorErrorReasonInvalidDescriptorKind,
			"descriptor kind %d has no validator",
			desc.code,
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

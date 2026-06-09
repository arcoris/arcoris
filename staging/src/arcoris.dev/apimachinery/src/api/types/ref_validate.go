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

// validateRef checks DescriptorRef syntax, optional resolver lookup, and cycles.
func validateRef(desc Descriptor, resolver Resolver, path string, resolving map[TypeName]bool) error {
	name := desc.ref.name

	if !name.IsValid() {
		return descriptorErrorf(
			path,
			ErrInvalidDescriptorReference,
			DescriptorErrorReasonInvalidReferenceName,
			"reference name %q is not a valid TypeName",
			name,
		)
	}

	if resolver == nil {
		return nil
	}

	if resolving[name] {
		return descriptorErrorf(
			path,
			ErrInvalidDescriptorReference,
			DescriptorErrorReasonReferenceCycle,
			"reference %q creates a recursive Definition graph",
			name,
		)
	}

	def, ok := resolver.Resolve(name)

	if !ok {
		return descriptorErrorf(
			path,
			ErrUnresolvedDescriptorReference,
			DescriptorErrorReasonUnknownReference,
			"reference %q was not found in resolver",
			name,
		)
	}

	next := copyResolving(resolving)
	next[name] = true

	if err := validateDescriptor(def.Descriptor(), resolver, path, next); err != nil {
		var descriptorErr *DescriptorError

		if errors.As(err, &descriptorErr) && descriptorErr.Reason == DescriptorErrorReasonReferenceCycle {
			return err
		}

		return descriptorErrorf(
			path,
			err,
			DescriptorErrorReasonInvalidResolvedDefinition,
			"resolved reference %q is structurally invalid: %v",
			name,
			err,
		)
	}

	return nil
}

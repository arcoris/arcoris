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

package resource

import (
	"fmt"

	"arcoris.dev/apimachinery/api/types"
)

// validateSurfaceLocal checks one resource API surface descriptor without a resolver.
//
// Local validation accepts direct objects and reference roots with valid syntax,
// but it does not prove that references resolve to object-like descriptors.
func validateSurfaceLocal(
	desc types.Descriptor,
	path string,
	invalidReason ErrorReason,
	objectReason ErrorReason,
	label string,
) error {
	if err := types.ValidateLocal(desc); err != nil {
		return nestedVersionError(
			path,
			invalidReason,
			invalidSurfaceDetail(label, err),
			err,
		)
	}

	switch desc.Code() {
	case types.DescriptorObject, types.DescriptorRef:
		return nil
	default:
		return versionError(
			path,
			objectReason,
			fmt.Sprintf("%s root must be object or reference to object, got %s", label, desc.Code()),
		)
	}
}

// validateSurfaceResolved checks one resource API surface descriptor with a resolver.
//
// Resolved validation accepts direct objects or references that resolve to
// object-like descriptors through resolver. It does not define live objects,
// metadata, storage, codecs, or runtime behavior.
func validateSurfaceResolved(
	desc types.Descriptor,
	resolver types.Resolver,
	path string,
	invalidReason ErrorReason,
	objectReason ErrorReason,
	label string,
) error {
	if err := types.ValidateResolved(desc, resolver); err != nil {
		return nestedVersionError(
			path,
			invalidReason,
			invalidSurfaceDetail(label, err),
			err,
		)
	}

	ok, detail := objectLikeResolved(desc, resolver, make(map[types.TypeName]bool), label, 0)
	if ok {
		return nil
	}

	return versionError(path, objectReason, detail)
}

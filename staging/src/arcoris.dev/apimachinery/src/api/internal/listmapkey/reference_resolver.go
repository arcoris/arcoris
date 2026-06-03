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

package listmapkey

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/typeref"
	"arcoris.dev/apimachinery/api/types"
)

// referenceResolver resolves descriptor references for one ListMap key extraction.
//
// The wrapper keeps listmapkey's failure taxonomy independent from typeref's
// internal taxonomy. Callers of listmapkey see ListMap selector failures, not a
// lower-level TypeRef traversal model.
type referenceResolver struct {
	refs *typeref.Resolver
}

// newReferenceResolver creates bounded resolver state for one extraction call.
//
// Options are normalized by listmapkey so this package keeps ownership of its
// default depth policy.
func newReferenceResolver(opts Options) referenceResolver {
	return referenceResolver{
		refs: typeref.New(opts.Resolver, opts.normalizedMaxDepth()),
	}
}

// resolve resolves a descriptor until it reaches a non-ref descriptor.
//
// Non-reference descriptors return unchanged and do not require a resolver.
// Reference failures are remapped to listmapkey FailureKind values so callers can
// still decide whether a failure is descriptor-related or payload-related.
func (r referenceResolver) resolve(
	path fieldpath.Path,
	descriptor types.Type,
	depth int,
) (types.Type, error) {
	if descriptor.Code() != types.TypeRef {
		return descriptor, nil
	}

	resolvedDescriptor, err := r.refs.ResolveFinal(path, descriptor, depth)
	if err != nil {
		return types.Type{}, listMapRefError(err)
	}

	return resolvedDescriptor, nil
}

// listMapRefError maps shared TypeRef traversal errors into ListMap key errors.
func listMapRefError(err error) error {
	refError, ok := typeref.AsError(err)
	if !ok {
		return err
	}

	return failure(refError.Path, listMapRefFailureKind(refError.Kind), refError.Detail)
}

// listMapRefFailureKind keeps TypeRef classifications inside listmapkey's
// existing FailureKind vocabulary.
func listMapRefFailureKind(kind typeref.FailureKind) FailureKind {
	switch kind {
	case typeref.FailureInvalidDescriptor:
		return FailureInvalidDescriptor
	case typeref.FailureUnresolvedRef:
		return FailureUnresolvedRef
	case typeref.FailureReferenceCycle:
		return FailureReferenceCycle
	default:
		return FailureInvalidDescriptor
	}
}

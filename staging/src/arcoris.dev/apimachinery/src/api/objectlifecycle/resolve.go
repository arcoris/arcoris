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

package objectlifecycle

import (
	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/resource"
)

// resolvedResource is one exact resource/version lookup result.
type resolvedResource struct {
	// definition is the resource family definition.
	definition resource.Definition

	// version is the selected version contract.
	version resource.VersionDefinition

	// gvr is the concrete store resource identity for the selected version.
	gvr identity.GroupVersionResource
}

// resolveObjectResource resolves the resource selected by object TypeMeta.
func (e *Executor) resolveObjectResource(op Operation, obj objectapply.ValueObject) (resolvedResource, error) {
	gvk := obj.GroupVersionKind()
	if err := gvk.Validate(); err != nil {
		return resolvedResource{}, errorFor(op, ErrorReasonInvalidRequest, objectstore.Key{}, ErrInvalidRequest, err)
	}

	def, version, ok := e.resources.ResolveVersionKind(gvk)
	if !ok {
		return resolvedResource{}, errorFor(op, ErrorReasonResourceNotFound, objectstore.Key{}, ErrResourceNotFound, nil)
	}

	gvr, ok := def.GroupVersionResource(version.Version())
	if !ok {
		return resolvedResource{}, errorFor(
			op,
			ErrorReasonInvalidResourceContract,
			objectstore.Key{},
			ErrValidationFailed,
			ErrInvalidResourceContract,
		)
	}

	return resolvedResource{definition: def, version: version, gvr: gvr}, nil
}

// resolveKeyResource resolves the resource selected by request GVR.
func (e *Executor) resolveKeyResource(op Operation, gvr identity.GroupVersionResource) (resolvedResource, error) {
	if err := gvr.Validate(); err != nil {
		return resolvedResource{}, errorFor(op, ErrorReasonInvalidRequest, objectstore.Key{}, ErrInvalidRequest, err)
	}

	def, version, ok := e.resources.ResolveVersionResource(gvr)
	if !ok {
		return resolvedResource{}, errorFor(op, ErrorReasonResourceNotFound, objectstore.Key{}, ErrResourceNotFound, nil)
	}

	return resolvedResource{definition: def, version: version, gvr: gvr}, nil
}

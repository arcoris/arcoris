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

package valuefieldset

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// extractor carries one extraction run's resolver and recursion state.
type extractor struct {
	resolver  types.Resolver
	maxDepth  int
	resolving map[types.TypeName]bool
}

// newExtractor normalizes options into executable extraction state.
func newExtractor(opts Options) *extractor {
	return &extractor{
		resolver:  opts.Resolver,
		maxDepth:  opts.normalizedMaxDepth(),
		resolving: make(map[types.TypeName]bool),
	}
}

// extract dispatches field-set extraction by descriptor type code.
func (e *extractor) extract(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Type,
	depth int,
) (fieldpath.Set, error) {
	if val.IsZero() {
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidValue,
			ErrorReasonInvalidZero,
			"value is the invalid zero Value",
		)
	}

	if !descriptor.IsValid() {
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has no valid type code",
		)
	}

	if val.IsNull() {
		return setAt(path)
	}

	switch descriptor.Code() {
	case types.TypeNull:
		return e.extractNull(path, val, descriptor)
	case types.TypeBool,
		types.TypeString,
		types.TypeBytes,
		types.TypeInt8,
		types.TypeInt16,
		types.TypeInt32,
		types.TypeInt64,
		types.TypeUint8,
		types.TypeUint16,
		types.TypeUint32,
		types.TypeUint64,
		types.TypeFloat32,
		types.TypeFloat64,
		types.TypeDecimal,
		types.TypeTimestamp,
		types.TypeDate,
		types.TypeTime,
		types.TypeDuration:
		return e.extractScalar(path, val, descriptor)
	case types.TypeObject:
		return e.extractObject(path, val, descriptor, depth)
	case types.TypeMap:
		return e.extractMap(path, val, descriptor, depth)
	case types.TypeList:
		return e.extractList(path, val, descriptor, depth)
	case types.TypeRef:
		return e.extractRef(path, val, descriptor, depth)
	default:
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has an unsupported type code",
		)
	}
}

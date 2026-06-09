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
	"arcoris.dev/apimachinery/api/internal/typeref"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// extractor carries one extraction run's resolver and recursion state.
type extractor struct {
	resolver types.Resolver
	maxDepth int
	refs     *typeref.Resolver
}

// newExtractor normalizes options into executable extraction state.
func newExtractor(opts Options) *extractor {
	maxDepth := opts.normalizedMaxDepth()

	return &extractor{
		resolver: opts.Resolver,
		maxDepth: maxDepth,
		refs:     typeref.New(opts.Resolver, maxDepth),
	}
}

// extract dispatches field-set extraction by descriptor kind.
func (e *extractor) extract(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
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
			"descriptor has no valid kind",
		)
	}

	if val.IsNull() {
		return setAt(path)
	}

	switch descriptor.Code() {
	case types.DescriptorNull:
		return e.extractNull(path, val, descriptor)
	case types.DescriptorBool,
		types.DescriptorString,
		types.DescriptorBytes,
		types.DescriptorInt8,
		types.DescriptorInt16,
		types.DescriptorInt32,
		types.DescriptorInt64,
		types.DescriptorUint8,
		types.DescriptorUint16,
		types.DescriptorUint32,
		types.DescriptorUint64,
		types.DescriptorFloat32,
		types.DescriptorFloat64,
		types.DescriptorDecimal,
		types.DescriptorTimestamp,
		types.DescriptorDate,
		types.DescriptorTime,
		types.DescriptorDuration:
		return e.extractScalar(path, val, descriptor)
	case types.DescriptorObject:
		return e.extractObject(path, val, descriptor, depth)
	case types.DescriptorMap:
		return e.extractMap(path, val, descriptor, depth)
	case types.DescriptorList:
		return e.extractList(path, val, descriptor, depth)
	case types.DescriptorRef:
		return e.extractRef(path, val, descriptor, depth)
	default:
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has an unsupported kind",
		)
	}
}

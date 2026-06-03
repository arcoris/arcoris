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

package valuemerge

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/typeref"
	"arcoris.dev/apimachinery/api/types"
)

// merger owns one merge traversal's shared descriptor state.
type merger struct {
	resolver types.Resolver
	maxDepth int
	refs     *typeref.Resolver
}

// newMerger constructs traversal-local state from public options.
func newMerger(opts Options) *merger {
	maxDepth := normalizedMaxDepth(opts.MaxDepth)

	return &merger{
		resolver: opts.Resolver,
		maxDepth: maxDepth,
		refs:     typeref.New(opts.Resolver, maxDepth),
	}
}

// merge applies selected overlay fields at path according to descriptor semantics.
func (m *merger) merge(
	path fieldpath.Path,
	base operand,
	overlay operand,
	descriptor types.Type,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	selection := selectAt(fields, path)
	if !selection.selected() {
		return base.Clone(), nil
	}
	if selection.exact {
		return m.replaceSubtree(path, overlay, descriptor, depth)
	}
	if descriptor.IsZero() {
		return operand{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is invalid",
		)
	}

	switch descriptor.Code() {
	case types.TypeRef:
		return m.mergeRef(path, base, overlay, descriptor, fields, depth)
	case types.TypeObject:
		return m.mergeObject(path, base, overlay, descriptor, fields, depth)
	case types.TypeMap:
		return m.mergeMap(path, base, overlay, descriptor, fields, depth)
	case types.TypeList:
		return m.mergeList(path, base, overlay, descriptor, fields, depth)
	default:
		return m.mergeScalar(path, base, overlay, descriptor)
	}
}

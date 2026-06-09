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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/internal/typeref"
	"arcoris.dev/apimachinery/api/types"
)

// comparer holds one isolated comparison run.
//
// The resolver and depth limit are normalized once, then shared with the
// internal DescriptorRef helper for this one comparison run.
type comparer struct {
	// resolver loads named descriptor definitions for DescriptorRef nodes.
	resolver types.Resolver
	// maxDepth is the effective DescriptorRef hop limit for this run.
	maxDepth int
	// refs tracks the active DescriptorRef stack for cycle detection.
	refs *typeref.Resolver
}

// newComparer copies user options into fresh per-run state.
func newComparer(opts Options) *comparer {
	maxDepth := opts.normalizedMaxDepth()

	return &comparer{
		resolver: opts.Resolver,
		maxDepth: maxDepth,
		refs:     typeref.New(opts.Resolver, maxDepth),
	}
}

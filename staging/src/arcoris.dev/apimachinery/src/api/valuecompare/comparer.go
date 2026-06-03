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

import "arcoris.dev/apimachinery/api/types"

// comparer holds per-run state that must not leak between Compare calls.
//
// The resolver and depth limit are normalized once. The resolving map is the
// active TypeRef stack, used only while descending through references.
type comparer struct {
	resolver  types.Resolver
	maxDepth  int
	resolving map[types.TypeName]bool
}

// newComparer converts user options into immutable run configuration.
func newComparer(opts Options) *comparer {
	return &comparer{
		resolver:  opts.Resolver,
		maxDepth:  opts.normalizedMaxDepth(),
		resolving: make(map[types.TypeName]bool),
	}
}

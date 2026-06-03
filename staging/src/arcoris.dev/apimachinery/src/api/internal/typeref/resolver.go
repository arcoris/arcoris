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

package typeref

import "arcoris.dev/apimachinery/api/types"

// Resolver carries one traversal's TypeRef state.
//
// The active stack is per operation. Sharing a Resolver across independent
// validation, extraction, or comparison runs would make cycle detection leak
// across calls.
type Resolver struct {
	resolver types.Resolver
	maxDepth int
	active   map[types.TypeName]bool
}

// New creates TypeRef resolver state for one descriptor traversal.
//
// maxDepth is expected to be normalized by the caller. Keeping normalization in
// public packages lets each package expose its own defaults without coupling
// this internal helper to those APIs.
func New(resolver types.Resolver, maxDepth int) *Resolver {
	return &Resolver{
		resolver: resolver,
		maxDepth: maxDepth,
		active:   make(map[types.TypeName]bool),
	}
}

// Enter marks name active until the returned cleanup function is called.
//
// Callers normally pair Enter with defer. Calling the returned function more
// than once is harmless because deleting a missing map key is a no-op.
func (r *Resolver) Enter(name types.TypeName) func() {
	r.active[name] = true

	return func() {
		delete(r.active, name)
	}
}

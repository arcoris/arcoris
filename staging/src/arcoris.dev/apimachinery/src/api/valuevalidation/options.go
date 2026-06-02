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

package valuevalidation

import "arcoris.dev/apimachinery/api/types"

const (
	// defaultMaxDepth bounds recursive TypeRef traversal when callers do not
	// provide an explicit limit.
	defaultMaxDepth = 64

	// defaultMaxErrors bounds diagnostic collection in the common case.
	defaultMaxErrors = 100
)

// Options configures descriptor-aware value validation.
//
// Zero Options are valid. The validator uses conservative defaults for
// recursion depth and diagnostic collection, and TypeRef validation reports
// unresolved-reference diagnostics when no resolver is available.
//
// Options do not request full descriptor validation. Descriptor preparation is
// an upstream responsibility of api/types validation and catalog registration.
type Options struct {
	// Resolver resolves TypeRef descriptors. It may be nil for descriptor trees
	// that do not contain references.
	Resolver types.Resolver

	// MaxDepth prevents runaway TypeRef recursion. Values <= 0 use the package
	// default.
	MaxDepth int

	// MaxErrors limits collected diagnostics. Values <= 0 use the package
	// default.
	MaxErrors int
}

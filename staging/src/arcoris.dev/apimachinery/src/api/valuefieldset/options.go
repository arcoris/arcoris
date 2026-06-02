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

import "arcoris.dev/apimachinery/api/types"

const defaultMaxDepth = 64

// Options configures one field-set extraction run.
//
// Options intentionally stays small because extraction assumes descriptors and
// values have already been prepared by their owning layers. The extractor only
// needs reference resolution and a recursion guard to produce stable semantic
// paths safely.
type Options struct {
	// Resolver resolves api/types TypeRef descriptors.
	Resolver types.Resolver

	// MaxDepth prevents runaway TypeRef recursion.
	//
	// Zero or a negative value uses the package default.
	MaxDepth int
}

// normalizedMaxDepth returns the effective TypeRef traversal budget.
func (o Options) normalizedMaxDepth() int {
	if o.MaxDepth > 0 {
		return o.MaxDepth
	}

	return defaultMaxDepth
}

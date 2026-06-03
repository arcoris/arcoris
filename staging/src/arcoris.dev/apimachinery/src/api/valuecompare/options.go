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

const defaultMaxDepth = 64

// Options configures one value comparison run.
//
// Zero Options are valid for descriptor graphs that do not contain TypeRef.
type Options struct {
	// Resolver resolves api/types TypeRef descriptors. It may be nil when the
	// compared descriptor tree contains no references.
	Resolver types.Resolver

	// MaxDepth prevents runaway TypeRef recursion.
	//
	// Values <= 0 use the package default.
	MaxDepth int
}

// normalizedMaxDepth returns the effective TypeRef traversal budget.
func (o Options) normalizedMaxDepth() int {
	if o.MaxDepth > 0 {
		return o.MaxDepth
	}

	return defaultMaxDepth
}

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

// defaultMaxDepth bounds DescriptorRef traversal when callers do not provide a limit.
const defaultMaxDepth = 64

// Options configures one value comparison run.
//
// Zero Options are valid for descriptor graphs that do not contain DescriptorRef.
// Options are copied into a comparer at the start of Compare or CompareAt.
type Options struct {
	// Resolver resolves api/types DescriptorRef descriptors during traversal.
	//
	// A nil resolver is accepted only when no reachable descriptor is DescriptorRef.
	Resolver types.Resolver

	// MaxDepth limits DescriptorRef hops before comparison reports ErrReferenceCycle.
	//
	// Values <= 0 use the package default.
	MaxDepth int
}

// normalizedMaxDepth returns the per-run DescriptorRef hop budget.
func (o Options) normalizedMaxDepth() int {
	if o.MaxDepth > 0 {
		return o.MaxDepth
	}

	return defaultMaxDepth
}

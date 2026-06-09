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

package valueapply

import "arcoris.dev/apimachinery/api/types"

// Options configures one value-level apply operation.
type Options struct {
	// Resolver resolves api/types DescriptorRef descriptors.
	Resolver types.Resolver

	// MaxDepth prevents runaway DescriptorRef recursion.
	//
	// Zero means package defaults in lower-level packages.
	MaxDepth int

	// Force allows Owner to take ownership of conflicting changed fields from
	// other owners.
	Force bool
}

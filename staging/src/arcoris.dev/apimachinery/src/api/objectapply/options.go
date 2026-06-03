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

package objectapply

import "arcoris.dev/apimachinery/api/types"

// Options configures one object-level apply operation.
//
// The fields intentionally mirror valueapply.Options so objectapply can pass
// traversal and ownership policy through without reinterpreting value-level
// semantics.
type Options struct {
	// Resolver resolves api/types TypeRef descriptors.
	Resolver types.Resolver

	// MaxDepth bounds TypeRef traversal in lower-level validators and apply.
	//
	// Zero means lower-level package defaults. objectapply does not define its
	// own recursion limit.
	MaxDepth int

	// Force allows value-level apply to take supported conflicting Desired
	// fields from other owners according to api/valueapply semantics. Unsupported
	// value-level takeovers remain errors and are preserved as lower-level causes.
	Force bool
}

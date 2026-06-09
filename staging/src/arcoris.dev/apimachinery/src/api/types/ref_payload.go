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

package types

// refPayload stores the named definition target for DescriptorRef.
//
// A reference points at a Definition through an owner-created Resolver. It
// is not a Go type, reflection handle, global registry lookup, or callback hook.
type refPayload struct {
	// name is resolved through an owner-created Resolver.
	//
	// ValidateResolved checks syntax and, when a resolver is supplied, resolution and
	// reference cycles.
	name TypeName
}

// cloneRefPayload copies the reference target.
func cloneRefPayload(p refPayload) refPayload {
	return p
}

// emptyRefPayload reports whether p has no configured DescriptorRef state.
func emptyRefPayload(p refPayload) bool {
	return p.name == ""
}

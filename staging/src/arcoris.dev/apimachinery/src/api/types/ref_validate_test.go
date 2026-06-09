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

import "testing"

func TestRefValidateResolutionAndCycles(t *testing.T) {
	requireErrorIs(t, ValidateLocal(Ref("bad").Descriptor()), ErrInvalidDescriptorReference)

	missing := resolverFunc(func(TypeName) (Definition, bool) {
		return Definition{}, false
	})
	requireErrorIs(t, ValidateResolved(Ref("example.Name").Descriptor(), missing), ErrUnresolvedDescriptorReference)

	cycle := resolverFunc(func(name TypeName) (Definition, bool) {
		switch name {
		case "example.A":
			return Define("example.A", Ref("example.B")), true
		case "example.B":
			return Define("example.B", Ref("example.A")), true
		default:
			return Definition{}, false
		}
	})
	requireErrorIs(t, ValidateResolved(Ref("example.A").Descriptor(), cycle), ErrInvalidDescriptorReference)
}

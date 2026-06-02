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

package listmapkey

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestReferenceResolverReturnsNonRefDescriptor(t *testing.T) {
	resolvedDescriptor, err := newReferenceResolver(Options{}).resolve(
		conditionPath(0),
		types.String().Type(),
		0,
	)

	requireNoError(t, err)
	requireEqual(t, resolvedDescriptor.Code(), types.TypeString)
}

func TestReferenceResolverResolvesReference(t *testing.T) {
	typeResolver := resolverFunc(func(name types.TypeName) (types.TypeDefinition, bool) {
		if name == "example.Name" {
			return types.Define("example.Name", types.String()), true
		}

		return types.TypeDefinition{}, false
	})

	resolvedDescriptor, err := newReferenceResolver(Options{Resolver: typeResolver}).resolve(
		conditionPath(0),
		types.Ref("example.Name").Type(),
		0,
	)

	requireNoError(t, err)
	requireEqual(t, resolvedDescriptor.Code(), types.TypeString)
}

func TestReferenceResolverRejectsMissingResolver(t *testing.T) {
	_, err := newReferenceResolver(Options{}).resolve(
		conditionPath(0),
		types.Ref("example.Name").Type(),
		0,
	)

	requireErrorKind(t, err, FailureUnresolvedRef)
}

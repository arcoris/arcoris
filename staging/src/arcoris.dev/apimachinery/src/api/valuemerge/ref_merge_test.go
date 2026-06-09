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

package valuemerge

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestMergeRefScalar(t *testing.T) {
	resolver := resolverFunc(func(name types.TypeName) (types.Definition, bool) {
		if name == "example.Name" {
			return types.Define("example.Name", types.String()), true
		}

		return types.Definition{}, false
	})

	got, err := Merge(
		str("old"),
		str("new"),
		types.Ref("example.Name").Descriptor(),
		pathSet(root()),
		Options{Resolver: resolver},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, str("new"))
}

func TestMergeRefMissingResolver(t *testing.T) {
	_, err := Merge(
		str("old"),
		str("new"),
		types.Ref("example.Name").Descriptor(),
		pathSet(root()),
		Options{},
	)

	requireErrorIs(t, err, ErrUnresolvedRef)
}

func TestMergeRefCycle(t *testing.T) {
	resolver := resolverFunc(func(name types.TypeName) (types.Definition, bool) {
		if name == "example.A" {
			return types.Define("example.A", types.Ref("example.A")), true
		}

		return types.Definition{}, false
	})

	_, err := Merge(
		str("old"),
		str("new"),
		types.Ref("example.A").Descriptor(),
		pathSet(root()),
		Options{Resolver: resolver},
	)

	requireErrorIs(t, err, ErrReferenceCycle)
}
